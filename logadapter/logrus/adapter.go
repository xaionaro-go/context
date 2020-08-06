package logrus

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/xaionaro-go/context"
)

func init() {
	context.LogAdapters = append(context.LogAdapters, Adapter{})
}

type Adapter struct{}

func (_ Adapter) SetLogLevel(logger context.MinimalLogger, logLevel context.LogLevel) context.MinimalLogger {
	switch logger := logger.(type) {
	case *logrus.Entry:
		logger.Logger.SetLevel(logrusLogLevel(logLevel))
		logger.Level = logrusLogLevel(logLevel)
		return logger
	case *logrus.Logger:
		logger.SetLevel(logrusLogLevel(logLevel))
		return logger
	}
	return nil
}

func (_ Adapter) WithField(logger context.MinimalLogger, field string, value interface{}) context.MinimalLogger {
	switch logger := logger.(type) {
	case *logrus.Entry:
		return logger.WithField(field, value)
	case *logrus.Logger:
		return logger.WithField(field, value)
	}
	return nil
}

func (_ Adapter) WithFields(logger context.MinimalLogger, fields context.Fields) context.MinimalLogger {
	switch logger := logger.(type) {
	case *logrus.Entry:
		return context.WrapLogger(logger.WithFields(map[string]interface{}(fields)))
	case *logrus.Logger:
		return context.WrapLogger(logger.WithFields(map[string]interface{}(fields)))
	}
	return nil
}

func logrusLogLevel(logLevel context.LogLevel) logrus.Level {
	switch logLevel {
	case context.LogLevelDebug:
		return logrus.DebugLevel
	case context.LogLevelInfo:
		return logrus.InfoLevel
	case context.LogLevelWarning:
		return logrus.WarnLevel
	case context.LogLevelError:
		return logrus.ErrorLevel
	case context.LogLevelPanic:
		return logrus.PanicLevel
	case context.LogLevelFatal:
		return logrus.FatalLevel
	default:
		panic(fmt.Sprintf("should never happened: %v", logLevel))
	}
}
