package zap

import (
	"fmt"

	"github.com/xaionaro-go/context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	context.LogAdapters = append(context.LogAdapters, Adapter{})
}

type Adapter struct{}

func (_ Adapter) SetLogLevel(logger context.MinimalLogger, logLevel context.LogLevel) context.MinimalLogger {
	if logger, ok := logger.(*zap.SugaredLogger); ok {
		return logger.Desugar().WithOptions(zap.IncreaseLevel(zapLogLevel(logLevel))).Sugar()
	}
	return nil
}

func (_ Adapter) WithField(logger context.MinimalLogger, field string, value interface{}) context.MinimalLogger {
	if logger, ok := logger.(*zap.SugaredLogger); ok {
		return logger.With(field, value)
	}
	return nil
}

func (_ Adapter) WithFields(logger context.MinimalLogger, fields context.Fields) context.MinimalLogger {
	if logger, ok := logger.(*zap.SugaredLogger); ok {
		args := make([]interface{}, 0, len(fields)*2)
		for k, v := range fields {
			args = append(args, k, v)
		}
		return context.WrapLogger(logger.With(args...))
	}
	return nil
}

func zapLogLevel(logLevel context.LogLevel) zapcore.Level {
	switch logLevel {
	case context.LogLevelDebug:
		return zapcore.DebugLevel
	case context.LogLevelInfo:
		return zapcore.InfoLevel
	case context.LogLevelWarning:
		return zapcore.WarnLevel
	case context.LogLevelError:
		return zapcore.ErrorLevel
	case context.LogLevelPanic:
		return zapcore.PanicLevel
	case context.LogLevelFatal:
		return zapcore.FatalLevel
	default:
		panic(fmt.Sprintf("should never happened: %v", logLevel))
	}
}
