package main

import (
	"kvparser/internal/config"
	"kvparser/internal/logger"
	"log"

	"github.com/gin-gonic/gin"

	svc "github.com/kardianos/service"
)

type program struct {
	logger logger.Logger
	cfg    *config.ServerConfig
}

func (p *program) Start(s svc.Service) error {
	if p.cfg.IsProduction {
		gin.SetMode(gin.ReleaseMode)
	} else {
		log.Println("Starting in development mode")

	}

	return nil
}

func (p *program) Stop(s svc.Service) error {
	return nil
}
