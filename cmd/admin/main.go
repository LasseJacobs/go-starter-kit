package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/LasseJacobs/go-starter-kit/internal/config"
	"github.com/LasseJacobs/go-starter-kit/pkg/log"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	logger := log.New()
	if len(os.Args) <= 3 {
		logger.Error("Usage:", os.Args[1], "command", "argument")
		return errors.New("invalid command")
	}

	cfg, err := config.Load(os.Args[2])
	if err != nil {
		logger.Errorf("failed to load application configuration: %s", os.Args[2])
		return err
	}

	logger, err = log.ConfigureLogging(cfg.Logging)
	if err != nil {
		logger.Errorf("failed to load application logger: %s", os.Args[2])
		return err
	}
	switch os.Args[1] {
	case "seed":
		err = Seed(cfg, logger, os.Args[3])
	default:
		err = errors.New("must specify a command")
	}

	if err != nil {
		return err
	}

	return nil
}

// Seed to populate database with seed data
func Seed(cfg *config.Config, logger *log.Logger, sqlFilename string) error {
	db, err := sqlx.Open(cfg.DB.Driver, cfg.DB.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	// load from SQL file
	bytes, err := ioutil.ReadFile(sqlFilename)
	if err != nil {
		return err
	}
	sql := string(bytes)
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	if _, err := tx.Exec(sql); err != nil {
		logger.Error(err)
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		logger.Error(err)
		return err
	}
	logger.Info("Seed data complete")
	return nil
}
