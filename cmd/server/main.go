// Package main запускает HTTP сервер приложения.
package main

import (
	"kvparser/internal/config"
	"kvparser/internal/logger"
	"log"
	"os"

	"github.com/kardianos/service"
)

var AppEnv string

func main() {

	// Загрузка конфигурации
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)

	}
	cfg.IsProduction = AppEnv == "production"
	logger := logger.NewLogrusLogger(cfg.LogPath)

	// Установка вывода стандартного log в logrus
	log.SetOutput(logger.Writer())
	log.SetFlags(0) // убираем timestamp, так как logrus добавит свой

	svcConfig := &service.Config{
		Name:        "KVParser",
		DisplayName: "KVParser Service",
		Description: "A tool for k-vrachu parsing",
	}

	prg := &program{logger: logger, cfg: cfg}
	s, err := service.New(prg, svcConfig)

	if err != nil {
		log.Fatal(err)
	}

	// Если запустили с параметром install/uninstall/start/stop
	if len(os.Args) > 1 {
		err = service.Control(s, os.Args[1])
		if err != nil {
			logger.Error("wrong service command: ", err)
		}
		return
	}

	// Запуск сервиса
	err = s.Run()
	if err != nil {
		logger.Error("failed run: ", err)
	}

}
