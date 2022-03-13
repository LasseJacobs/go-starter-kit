package main

import (
	"github.com/sirupsen/logrus"
	"os"
)

// Version set current code version
var Version = "1.0.0"

func main() {
	app := App()
	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
