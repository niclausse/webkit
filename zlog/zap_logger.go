package zlog

import "go.uber.org/zap"

func getZapLogger() *zap.Logger {
	if zapLogger == nil {
		zapLogger = newLogger().WithOptions(zap.AddCallerSkip(1))
	}

	return zapLogger
}
