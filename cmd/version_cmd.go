package main

import (
	"fmt"
	"github.com/urfave/cli"
)

var versionCmd = cli.Command{
	Name:   "version",
	Usage:  "show api version",
	Action: showVersion,
}

func showVersion(c *cli.Context) error {
	fmt.Println(Version)
	return nil
}
