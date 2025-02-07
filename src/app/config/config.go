package config

import (
	"errors"
	"go-mma/util/env"
)

var (
	ErrInvalidHTTPPort = errors.New("HTTP_PORT must be a positive integer")
	ErrNoDSN           = errors.New("DB_DSN is required")
)

type Config struct {
	HTTPPort int
	DSN      string
}

func Load() (*Config, error) {
	config := &Config{
		HTTPPort: env.GetIntDefault("HTTP_PORT", 8080),
		DSN:      env.Get("DB_DSN"),
	}
	err := config.Validate()
	if err != nil {
		return nil, err
	}
	return config, err
}

func (c *Config) Validate() error {
	if c.HTTPPort <= 0 {
		return ErrInvalidHTTPPort
	}

	if len(c.DSN) <= 0 {
		return ErrNoDSN
	}
	return nil
}
