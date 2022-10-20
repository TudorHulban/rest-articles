package configs

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

const _HTTPServerPort = 80

type Database struct {
	Host     string `envconfig:"DATABASE_HOST" required:"true"`
	Port     int    `envconfig:"DATABASE_PORT" required:"true"`
	User     string `envconfig:"DATABASE_USER" required:"true"`
	Password string `envconfig:"DATABASE_PASSWORD" required:"true"`
	Name     string `envconfig:"DATABASE_NAME" required:"true"`
}

type Config struct {
	Database       Database
	HTTPServerPort int `envconfig:"SERVER_PORT" default:"80"`
}

func NewParsedConfig(path string) (*Config, error) {
	if errLoad := godotenv.Load(path); errLoad != nil {
		return nil, errLoad
	}

	var config Config

	err := envconfig.Process("", &config)

	return &config, err
}
