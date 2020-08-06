package context

type MinimalLogger interface {
	Errorf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

type Logger interface {
	MinimalLogger

	SetLogLevel(LogLevel)
	WithField(key string, value interface{}) Logger
	WithFields(Fields) Logger
}

type loggerWrapper struct {
	MinimalLogger
}

func WrapLogger(logger MinimalLogger) Logger {
	if logger == nil {
		logger = &dummyLogger{}
	}
	return &loggerWrapper{MinimalLogger: logger}
}

type LogAdapter interface {
	SetLogLevel(MinimalLogger, LogLevel) MinimalLogger
	WithField(MinimalLogger, string, interface{}) MinimalLogger
	WithFields(MinimalLogger, Fields) MinimalLogger
}

var (
	LogAdapters []LogAdapter
)

func (logger *loggerWrapper) SetLogLevel(logLevel LogLevel) {
	for _, logAdapter := range LogAdapters {
		if modifiedLogger := logAdapter.SetLogLevel(logger.MinimalLogger, logLevel); modifiedLogger != nil {
			logger.MinimalLogger = modifiedLogger
			return
		}
	}
}

func (logger *loggerWrapper) WithField(field string, value interface{}) Logger {
	for _, logAdapter := range LogAdapters {
		if modifiedLogger := logAdapter.WithField(logger.MinimalLogger, field, value); modifiedLogger != nil {
			return &loggerWrapper{MinimalLogger: modifiedLogger}
		}
	}

	return logger
}

func (logger *loggerWrapper) WithFields(fields Fields) Logger {
	for _, logAdapter := range LogAdapters {
		if modifiedLogger := logAdapter.WithFields(logger.MinimalLogger, fields); modifiedLogger != nil {
			return &loggerWrapper{MinimalLogger: modifiedLogger}
		}
	}

	return logger
}

type dummyLogger struct{}

func (l *dummyLogger) Fatalf(format string, args ...interface{}) {}
func (l *dummyLogger) Panicf(format string, args ...interface{}) {}
func (l *dummyLogger) Errorf(format string, args ...interface{}) {}
func (l *dummyLogger) Warnf(format string, args ...interface{})  {}
func (l *dummyLogger) Infof(format string, args ...interface{})  {}
func (l *dummyLogger) Debugf(format string, args ...interface{}) {}
