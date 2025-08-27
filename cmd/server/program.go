package main

import (
	"fmt"
	"kvparser/internal/config"
	"kvparser/internal/domain"
	"kvparser/internal/infrastructure"
	"kvparser/internal/jobs"
	"kvparser/internal/logger"
	"kvparser/internal/services"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"

	svc "github.com/kardianos/service"
)

type program struct {
	logger    logger.Logger
	cfg       *config.ServerConfig
	parser    services.ChromeParserService
	scheduler *gocron.Scheduler
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

	tg, err := infrastructure.NewTgBot(p.cfg.TgToken, p.cfg.TgChatId)
	if err != nil {
		return fmt.Errorf("can't create bot api: %w", err)
	}

	rem := services.NewRemember[domain.DoctorMatch]()

	sched := gocron.NewScheduler(time.UTC)

	cron := "*/2 * * * *"
	if p.cfg.Cron != "" {
		cron = p.cfg.Cron
	}

	if err := jobs.RegisterFetchSlotsJob(sched, rem, cron, p.logger, svc, domain.DoctorOptions{
		Subdivision:  p.cfg.FilterOptions.Subdivision,
		Specialists:  p.cfg.FilterOptions.Specialists,
		Specialities: p.cfg.FilterOptions.Specialities,
	}, tg); err != nil {
		return err
	}

	p.parser = svc
	p.scheduler = sched

	go sched.StartBlocking()

	return nil
}

func (p *program) Stop(s svc.Service) error {
	p.logger.Info("Stopping service...")

	if p.scheduler != nil {
		p.scheduler.Stop()
	}

	if p.parser != nil {
		p.parser.Close()
	}

	return nil
}
