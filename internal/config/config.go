package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	BaseUrl       string `envconfig:"BASE_URL"`
	Token         string `envconfig:"TOKEN"`
	JWTSigningKey string `envconfig:"JWT_SIGNING_KEY"`
	Address       string `envconfig:"ADDRESS"`
	DatabaseUrl   string `envconfig:"DATABASE_URL"`
}

func NewConfig() (*Config, error) {
	cfg := new(Config)
	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
