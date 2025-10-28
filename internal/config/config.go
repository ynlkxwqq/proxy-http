package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

const (
	defaultAppPort = "8080"
	defaultAppPath = "/"
)

type (
	Configs struct {
		APP AppConfig
	}

	AppConfig struct {
		Port string
		Path string
	}
)

func New() (cfg Configs, err error) {
	root, err := os.Getwd()
	if err != nil {
		return
	}

	err = godotenv.Load(filepath.Join(root, ".env"))
	if err != nil {
		return
	}

	cfg.APP = AppConfig{
		Port: defaultAppPort,
		Path: defaultAppPath,
	}

	if err = envconfig.Process("APP", &cfg.APP); err != nil {
		return
	}

	return
}
