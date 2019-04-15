package logger

import (
	"context"

	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func Init(isProd bool) {
	var baseLogger *zap.Logger
	var err error
	if isProd {
		baseLogger, err = zap.NewProduction()
	} else {
		baseLogger, err = zap.NewDevelopment()
	}
	if err != nil {
		panic("Can't start logger: " + err.Error())
	}

	log := baseLogger.Sugar()
	log.Info("Logger initiated")

	logger = baseLogger.Sugar()
}

func New(ctx context.Context) *zap.SugaredLogger {
	newLogger := logger
	if ctx != nil {
		if ctxUserId, ok := ctx.Value("UserID").(string); ok {
			newLogger = newLogger.With(zap.String("userID", ctxUserId))
		}
	}

	return newLogger
}
