package main

import (
	"fmt"
	"github.com/LasseJacobs/go-starter-kit/internal/config"
	"github.com/LasseJacobs/go-starter-kit/pkg/log"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// App will setup CLI
func App() *cli.App {
	app := &cli.App{
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "config",
				Usage: "Load configuration from `FILE`",
			},
		},
		Action: func(c *cli.Context) error {
			fmt.Println(c.String("config"))
			return nil
		},
		Commands: []cli.Command{
			serveCmd,
			versionCmd,
			adminCmd,
		},
	}

	return app
}

// execWithConfig requires `config` flag is set
func execWithConfig(c *cli.Context, fn func(c *cli.Context, conf *config.Config) error) error {
	var configFile = c.GlobalString("config")
	conf, err := config.Load("APP", configFile)
	if err != nil {
		logrus.Fatalf("Failed to load configuration: %+v", err)
	}

	//todo: check if this is persistent
	_, err = log.ConfigureLogging(conf.Logging)
	if err != nil {
		logrus.Fatalf("Failed to create logger: %+v", err)
	}

	return fn(c, conf)
}
