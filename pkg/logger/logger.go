package logger

import (
	// "fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func SetupLogger(env string) (*zap.Logger, error) {
	if env == "production" {
		return zap.NewProduction()
	}

	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	log, err := config.Build()
	if err != nil {
		return nil, err
	}

	return log, nil
}
