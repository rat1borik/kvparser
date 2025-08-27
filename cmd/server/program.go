package main

import (
	"fmt"
	"kvparser/internal/config"
	"kvparser/internal/domain"
	"kvparser/internal/logger"
	"kvparser/internal/services"

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
		p.logger.Info("Starting in development mode")
	}

	svc, err := services.NewChromeParser(p.logger)
	if err != nil {
		return err
	}

	p.parser = svc

	res, err := svc.DoctorsSchedulePage()
	if err != nil {
		return err
	}

	parsed, err := domain.FindMatches(res, domain.DoctorOptions{})
	if err != nil {
		return err
	}

	for _, val := range parsed.Matches {
		fmt.Printf("%s %s %s %s\n", val.Name, val.Speciality, val.Status, val.Subdivision)
	}

	return nil
}

func (p *program) Stop(s svc.Service) error {
	p.logger.Info("Stopping service...")
	if p.parser != nil {
		p.parser.Close()
	}

	return nil
}
