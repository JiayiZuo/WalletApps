package common

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLogger() {
	logLevelStr := os.Getenv("LOG_LEVEL")
	logFile := os.Getenv("LOG_FILE")
	if logFile == "" {
		logFile = "wallet.log"
	}

	level := zapcore.InfoLevel
	if strings.ToLower(logLevelStr) == "debug" {
		level = zapcore.DebugLevel
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(level)
	cfg.OutputPaths = []string{"stdout", logFile}
	cfg.Encoding = "json"

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	Logger = logger
}
