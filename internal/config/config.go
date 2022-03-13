package config

import (
	"github.com/LasseJacobs/go-starter-kit/pkg/cleanenv"
	"github.com/LasseJacobs/go-starter-kit/pkg/log"
)

// Config holds data for application configuration
type Config struct {
	Name    string             `yaml:"name" required:"true"`
	Desc    string             `yaml:"description" required:"true"`
	Server  *Server            `yaml:"server" required:"true"`
	DB      *Database          `yaml:"database" required:"true"`
	Logging *log.LoggingConfig `yaml:"logging" required:"true"`
}

// Database holds data for database configuration
type Database struct {
	Driver   string `yaml:"driver" required:"true"`
	URL      string `yaml:"url" required:"true"`
	Password string `yaml:"password" required:"true"`
}

// Server holds data for server configuration
type Server struct {
	Port string `yaml:"port" required:"true"`
}

// Load returns config from yaml and environment variables.
func Load(file string) (*Config, error) {
	// default config
	c := Config{}
	var err = cleanenv.ReadConfig(file, "APP", &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
