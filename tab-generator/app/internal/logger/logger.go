package logger

import (
	"sync"

	"go.uber.org/zap"
)

var (
	Log  *zap.Logger
	once sync.Once
)

func init() {
	once.Do(func() {
		Log, _ = zap.NewDevelopment()
	})
}

func Info(msg string, fields ...zap.Field) {
	Log.Info(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	Log.Debug(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	Log.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Log.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	Log.Fatal(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	Log.Panic(msg, fields...)
}
