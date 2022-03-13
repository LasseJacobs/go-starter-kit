package main

import (
	"github.com/LasseJacobs/go-starter-kit/internal/api"
	"github.com/LasseJacobs/go-starter-kit/internal/config"
	"github.com/LasseJacobs/go-starter-kit/internal/storage"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var serveCmd = cli.Command{
	Name:  "serve",
	Usage: "Start API server",
	Action: func(c *cli.Context) error {
		return execWithConfig(c, run)
	},
}

func run(c *cli.Context, conf *config.Config) error {
	// connect to the database
	db, err := storage.Dial(conf.DB)
	if err != nil {
		logrus.Errorf("failed to connect to database: %s", err)
		return err
	}

	api := api.NewAPIWithVersion(conf, db, Version)
	return api.Start()
}
