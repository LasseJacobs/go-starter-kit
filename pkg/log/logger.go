package log

import (
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

type Logger struct {
	*logrus.Entry
}

type LoggingConfig struct {
	Level            string                 `yaml:"level" env:"LEVEL" env-default:"DEBUG"`
	File             string                 `yaml:"file" env:"FILE"`
	DisableColors    bool                   `split_words:"true" yaml:"disable_colors"`
	QuoteEmptyFields bool                   `split_words:"true" yaml:"quote_empty_fields"`
	TSFormat         string                 `yaml:"ts_format" env:"TS_FORMAT"`
	Fields           map[string]interface{} `yaml:"fields" env:"FIELDS"`
}

func New() *Logger {
	return &Logger{logrus.NewEntry(logrus.StandardLogger())}
}

func ConfigureLogging(config LoggingConfig) (*Logger, error) {
	logger := logrus.New()

	tsFormat := time.RFC3339Nano
	if config.TSFormat != "" {
		tsFormat = config.TSFormat
	}
	// always use the full timestamp
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:    true,
		DisableTimestamp: false,
		TimestampFormat:  tsFormat,
		DisableColors:    config.DisableColors,
		QuoteEmptyFields: config.QuoteEmptyFields,
	})

	// use a file if you want
	if config.File != "" {
		f, errOpen := os.OpenFile(config.File, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0664)
		if errOpen != nil {
			return nil, errOpen
		}
		logger.SetOutput(f)
		logger.Infof("Set output file to %s", config.File)
	}

	if config.Level != "" {
		level, err := logrus.ParseLevel(config.Level)
		if err != nil {
			return nil, err
		}
		logger.SetLevel(level)
		logger.Debug("Set log level to: " + logger.GetLevel().String())
	}

	f := logrus.Fields{}
	for k, v := range config.Fields {
		f[k] = v
	}

	return &Logger{logger.WithFields(f)}, nil
}
