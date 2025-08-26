package main

import (
	"fmt"
	"kvparser/internal/config"
	"kvparser/internal/logger"
	"kvparser/internal/services"
	"log"

	"github.com/gin-gonic/gin"

	svc "github.com/kardianos/service"
)

type program struct {
	logger logger.Logger
	cfg    *config.ServerConfig
	parser services.ChromeParserService
}

func (p *program) Start(s svc.Service) error {
	if p.cfg.IsProduction {
		gin.SetMode(gin.ReleaseMode)
	} else {
		log.Println("Starting in development mode")

	}

	svc, err := services.NewChromeParser()
	if err != nil {
		return err
	}

	p.parser = svc

	content, err := svc.ParsePage("https://12341ya.r1234u")
	if err != nil {
		return err
	}

	fmt.Println(content)

	return nil
}

func (p *program) Stop(s svc.Service) error {
	log.Println("Stopping service...")
	if p.parser != nil {
		p.parser.Close()
	}

	return nil
}
