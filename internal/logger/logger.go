// internal/logger/logger.go
package logger

import "go.uber.org/zap"

var globalLogger *zap.Logger

func Init() {
    globalLogger, _ = zap.NewDevelopment()
}

func L() *zap.Logger {
    return globalLogger
}

// Для компонентов
func Component(name string) *zap.Logger {
    return globalLogger.With(zap.String("component", name))
}