package applog

import "go.uber.org/zap"

type Logger struct {
	Log *zap.Logger
}

func NewLogger() *Logger {
	lg, _ := zap.NewProduction()
	return &Logger{
		Log: lg,
	}
}
