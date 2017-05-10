package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Opts struct {
	Environment string
	Level       zapcore.Level
}

// Init allocates new logger based on given options.
func Init(opts Opts) (*zap.Logger, error) {
	var cfg zap.Config
	switch opts.Environment {
	case "production":
		cfg = zap.NewProductionConfig()

	case "development":
		cfg = zap.NewDevelopmentConfig()
	}
	if opts.Level >= zapcore.DebugLevel && opts.Level <= zapcore.FatalLevel {
		cfg.Level.SetLevel(opts.Level)
	}

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	logger.Info("logger has been initialized", zap.String("environment", opts.Environment), zap.String("level", opts.Level.String()))

	return logger, nil
}
