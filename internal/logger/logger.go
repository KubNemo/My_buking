package logger

import (
	"go.uber.org/zap"
)

// Init инициализирует логгер в зависимости от окружения (prod или dev)
func Init(isProd bool) (*zap.Logger, error) {
	if isProd {
		return zap.NewProduction()
	}
	return zap.NewDevelopment()
}
