package services

import (
	"context"
	"fmt"
	"kvparser/internal/logger"
	"os"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

type ChromeParserService interface {
	DoctorsSchedulePage() (string, error)
	Close()
}

type chromeParserService struct {
	ctx    context.Context
	cancel context.CancelFunc
	logger logger.Logger
}

// Конструктор сервиса
func NewChromeParser(logger logger.Logger) (ChromeParserService, error) {
	dir, err := os.MkdirTemp("", "chromedp-userdata")
	if err != nil {
		logger.Error(err)
	}

	logFn := func(format string, args ...interface{}) {
		logger.Info(fmt.Sprintf(format, args))
	}

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.UserDataDir(dir),
		chromedp.Headless,
		chromedp.NoSandbox,
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-software-rasterizer", true),
	)

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)

	// Родительский контекст (без таймаута, сервис живёт долго)
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(logFn))

	chromedp.ListenTarget(ctx, func(ev interface{}) {
		go func(ev interface{}) {
			switch e := ev.(type) {
			case *runtime.EventConsoleAPICalled:
				for _, arg := range e.Args {
					logger.Info("JS console: %s", arg.Value)
				}
			case *runtime.EventExceptionThrown:
				logger.Info("JS error: %v", e.ExceptionDetails)
			}
		}(ev)
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
		logger: logger,
	}, nil
}

// Метод парсинга
func (svc *chromeParserService) DoctorsSchedulePage() (string, error) {
	var res string

	loadCtx, cancel := context.WithTimeout(svc.ctx, 10*time.Second)
	defer cancel()

	targetUrl := "https://k-vrachu.cifromed35.ru/service/hospitals/doctors/12600087?per_page=999999&type=by_unit"

	tasks := chromedp.Tasks{
		network.Enable(),
		network.SetCacheDisabled(true),
		chromedp.Navigate(targetUrl),
		chromedp.InnerHTML(".docsInLpuTable", &res, chromedp.NodeVisible),
	}

	if err := chromedp.Run(loadCtx, tasks); err != nil {
		return "", err
	}

	return res, nil
}

// Корректное завершение
func (svc *chromeParserService) Close() {
	svc.cancel()
}
