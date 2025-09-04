package logger

import (
	"go.uber.org/zap"
)

var log *zap.Logger

func Init() {

	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stderr"}

	var err error
	log, err = cfg.Build()
	if err != nil {
		panic(err)
	}
}

func L() *zap.Logger {
	if log == nil {
		Init()
	}
	return log
}

func Sync() {
	_ = log.Sync()
}
