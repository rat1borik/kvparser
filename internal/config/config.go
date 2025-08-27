// Package config описывает конфигурацию приложения.
package config

import (
	"kvparser/internal/utils"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	TgToken       string `yaml:"tg_token"`
	TgChatId      int64  `yaml:"tg_chat_id"`
	LogPath       string `yaml:"log_path"`
	Cron          string `yaml:"cron"`
	FilterOptions struct {
		Subdivision  string   `yaml:"subdivision"`
		Specialists  []string `yaml:"specialists"`
		Specialities []string `yaml:"specialities"`
	} `yaml:"filter_options"`
	// Берется из env
	IsProduction bool
}

func LoadConfig(filename string) (*ServerConfig, error) {
	path, err := utils.ExecPath(filename)
	if err != nil {
		log.Fatalf("can't load config: %v", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg ServerConfig
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
