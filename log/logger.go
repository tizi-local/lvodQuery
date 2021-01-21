package log

import (
	"fmt"
	"github.com/tizi-local/llib/log"
	"github.com/tizi-local/lvodQuery/config"
	"os"
	"path/filepath"
)

var (
	logger *log.Logger
)

func init() {
}

func NewLogger(config *config.LoggerConfig) *log.Logger {
	logger = log.New()
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	logger.SetReportCaller(true)

	level := log.InfoLevel
	switch config.Level {
	case "debug":
		level = log.DebugLevel
	case "info":
		level = log.InfoLevel
	case "warn":
		level = log.WarnLevel
	case "error":
		level = log.ErrorLevel
	case "fatal":
		level = log.FatalLevel
	}

	if err := os.MkdirAll(filepath.Base(config.Path), 0666); err != nil {
		panic("create directory failed" + config.Path)
	}
	f, err := os.OpenFile(config.Path, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_SYNC, 0666)
	if err != nil {
		fmt.Printf("log file init failed: %s\n", err.Error())
		return nil
	}
	logger.SetOutput(f)
	logger.SetLevel(level)

	return logger
}

func Default() *log.Logger {
	if logger == nil {
		logger = NewLogger(&config.LoggerConfig{})
	}
	return logger
}
