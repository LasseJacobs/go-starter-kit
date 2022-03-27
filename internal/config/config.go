package config

import (
	"github.com/LasseJacobs/go-starter-kit/pkg/cleanenv"
	"github.com/LasseJacobs/go-starter-kit/pkg/log"
)

// Config holds data for application configuration
type Config struct {
	Name    string            `yaml:"name" env:"NAME" env-default:"app"`
	Desc    string            `yaml:"description" env:"DESCRIPTION" env-default:"app template"`
	Server  Server            `yaml:"server" env-prefix:"SERVER_"`
	DB      Database          `yaml:"database" env-prefix:"DB_"`
	Logging log.LoggingConfig `yaml:"logging" env-prefix:"LOGGING_"`
}

// Database holds data for database configuration
type Database struct {
	User     string `yaml:"user" env:"USER" env-required:"true"`
	Password string `yaml:"password" env:"PASSWORD" env-required:"true"`
	Host     string `yaml:"host" env:"HOST" env-required:"true"`
	Database string `yaml:"database" env:"DATABASE" env-required:"true"`

	Port    string `yaml:"port" env:"PORT" env-default:"5432"`
	SSLMode string `yaml:"SSLMode" env:"SSL_MODE" env-default:"disable"`
}

// Server holds data for server configuration
type Server struct {
	Port string `yaml:"port" env-default:"8080"`
}

// Load returns config from yaml and environment variables.
func Load(prefix string, file string) (*Config, error) {
	// default config
	c := Config{}
	var err = cleanenv.ReadConfig(file, prefix, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
