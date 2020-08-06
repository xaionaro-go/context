package context

import (
	"flag"
	"fmt"
	"strings"

	"github.com/spf13/pflag"
)

type LogLevel int

const (
	LogLevelUndefined = LogLevel(iota)
	LogLevelDebug
	LogLevelInfo
	LogLevelWarning
	LogLevelError
	LogLevelPanic
	LogLevelFatal
)

func ParseLogLevel(in string) (LogLevel, error) {
	switch strings.ToLower(in) {
	case "d", "debug":
		return LogLevelDebug, nil
	case "i", "info":
		return LogLevelInfo, nil
	case "w", "warning", "warn":
		return LogLevelWarning, nil
	case "e", "err", "error":
		return LogLevelError, nil
	case "p", "panic":
		return LogLevelPanic, nil
	case "f", "fatal":
		return LogLevelFatal, nil
	}
	return LogLevelUndefined, fmt.Errorf("unknown logging level '%s'", in)
}

func (logLevel LogLevel) String() string {
	switch logLevel {
	case LogLevelUndefined:
		return "undefined"
	case LogLevelDebug:
		return "debug"
	case LogLevelInfo:
		return "info"
	case LogLevelWarning:
		return "warning"
	case LogLevelError:
		return "error"
	case LogLevelPanic:
		return "panic"
	case LogLevelFatal:
		return "fatal"
	}
	return "unknown"
}

func (logLevel *LogLevel) Set(value string) error {
	newLogLevel, err := ParseLogLevel(value)
	if err != nil {
		return err
	}
	*logLevel = newLogLevel
	return nil
}

func (logLevel *LogLevel) Type() string {
	return "LogLevel"
}

type FlagSet interface {
	Var(value flag.Value, name string, usage string)
}

func LogLevelFlag(varFunc func(pflag.Value, string, string), name string, value string, usage string) *LogLevel {
	if varFunc == nil {
		varFunc = func(value pflag.Value, n, u string) {
			flag.CommandLine.Var(value, n, u)
		}
	}
	logLevel, err := ParseLogLevel(value)
	if err != nil {
		panic(err)
	}
	varFunc(&logLevel, name, usage)
	return &logLevel
}
