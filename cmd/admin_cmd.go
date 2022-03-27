package main

import (
	"github.com/LasseJacobs/go-starter-kit/internal/config"
	"github.com/LasseJacobs/go-starter-kit/internal/storage"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"io/ioutil"
)

var adminCmd = cli.Command{
	Name:  "admin",
	Usage: "Run admin processes",
	Subcommands: cli.Commands{
		{
			Name:  "seed",
			Usage: "Seed the database with an input file",
			Action: func(c *cli.Context) error {
				return execWithConfig(c, seed)
			},
		},
	},
}

func seed(c *cli.Context, conf *config.Config) error {
	seedFile := c.Args().Get(0)
	content, err := ioutil.ReadFile(seedFile)
	if err != nil {
		return err
	}
	conn, err := storage.Dial(conf.DB)
	if err != nil {
		return err
	}
	res := conn.MustExec(string(content))
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	logrus.Infof("Seeding completed, rows affected: %d", count)
	return nil
}
