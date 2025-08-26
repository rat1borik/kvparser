package services

import (
	"context"
	"log"
	"os"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

type ChromeParserService interface {
	ParsePage(string) (string, error)
	Close()
}

type chromeParserService struct {
	ctx    context.Context
	cancel context.CancelFunc
}

// Конструктор сервиса
func NewChromeParser() (ChromeParserService, error) {
	dir, err := os.MkdirTemp("", "chromedp-userdata")
	if err != nil {
		log.Fatal(err)
	}

	logger := func(format string, args ...interface{}) {
		log.Printf(format, args...)
	}

	opts := append(chromedp.DefaultExecAllocatorOptions[:]) // chromedp.DisableGPU,
	// chromedp.UserDataDir(dir),
	// chromedp.Headless,
	// chromedp.NoSandbox,
	// chromedp.Flag("disable-dev-shm-usage", true),
	// chromedp.Flag("disable-software-rasterizer", true),

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)

	// Родительский контекст (без таймаута, сервис живёт долго)
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(logger))

	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch e := ev.(type) {
		case *runtime.EventConsoleAPICalled:
			for _, arg := range e.Args {
				log.Printf("JS console: %s", arg.Value)
			}
		case *runtime.EventExceptionThrown:
			log.Printf("JS error: %v", e.ExceptionDetails)
		}
	})

	cancelAll := func() {
		cancel()
		allocCancel()
		os.RemoveAll(dir)
	}

	// Можно сразу прогреть браузер
	if err := chromedp.Run(ctx, chromedp.Navigate("about:blank")); err != nil {
		cancelAll()
		return nil, err
	}

	return &chromeParserService{
		ctx:    ctx,
		cancel: cancelAll,
	}, nil
}

// Метод парсинга
func (svc *chromeParserService) ParsePage(url string) (string, error) {
	var res string

	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.Text("body", &res),
	}

	if err := chromedp.Run(svc.ctx, tasks); err != nil {
		return "", err
	}

	return res, nil
}

// Корректное завершение
func (svc *chromeParserService) Close() {
	svc.cancel()
}
