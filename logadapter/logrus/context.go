package logrus

import (
	"github.com/sirupsen/logrus"
	"github.com/xaionaro-go/context"
)

func NewCLIContext(logLevel context.LogLevel) context.Context {
	logger := logrus.New()
	logger.SetLevel(logrusLogLevel(logLevel))
	entry := logrus.NewEntry(logger)
	entry.Level = logrusLogLevel(logLevel)
	entry.Logger.SetFormatter(&logrus.JSONFormatter{})
	return context.NewContext(context.Background(), context.NewTraceID(), entry, nil, nil)
}
