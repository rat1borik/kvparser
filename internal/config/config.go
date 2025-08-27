// Package config описывает конфигурацию приложения.
package config

import (
	"kvparser/internal/utils"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Cookie  string `yaml:"cookie"`
	LogPath string `yaml:"log_path"`
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
