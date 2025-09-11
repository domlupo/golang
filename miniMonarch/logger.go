package main

import (
	"log/slog"
	"os"
)

type level int
type logType int

const (
	debugLevel = iota
	infoLevel
	warnLevel
	errorLevel
)

const (
	noLogType = iota
	mainLogType
)

const initialLogLevel = slog.LevelInfo
const filterTimeAndLevelFromLogs = true

type l struct {
	*slog.Logger
	*slog.LevelVar
	enabledLogTypes map[logType]bool
}

func newLogger() l {
	var levelVar = new(slog.LevelVar)
	levelVar.Set(initialLogLevel)
	logOptions := &slog.HandlerOptions{
		Level: levelVar,
		ReplaceAttr: func(groups []string, attribute slog.Attr) slog.Attr {
			if filterTimeAndLevelFromLogs && (attribute.Key == "time" || attribute.Key == "level") {
				return slog.Attr{}
			} else {
				return attribute
			}
		},
	}

	var logger = slog.New(slog.NewJSONHandler(os.Stderr, logOptions))
	l := l{
		Logger:          logger,
		LevelVar:        levelVar,
		enabledLogTypes: map[logType]bool{mainLogType: true},
	}

	return l
}

func (logger l) toggleLevel() {
	if logger.LevelVar.Level() == slog.LevelError {
		logger.Set(slog.LevelWarn)
	} else if logger.LevelVar.Level() == slog.LevelWarn {
		logger.Set(slog.LevelInfo)
	} else if logger.LevelVar.Level() == slog.LevelInfo {
		logger.Set(slog.LevelDebug)
	} else if logger.LevelVar.Level() == slog.LevelDebug {
		logger.Set(slog.LevelError)
	}
}

func (logger l) toggleType(logType logType) {
	logger.enabledLogTypes[logType] = !logger.enabledLogTypes[logType]
}

func (logger l) log(logType logType, level level, msg string, args ...any) {
	if !logger.enabledLogTypes[logType] {
		return
	}

	switch level {
	case debugLevel:
		logger.Debug(msg, args...)
	case infoLevel:
		logger.Info(msg, args...)
	case warnLevel:
		logger.Warn(msg, args...)
	case errorLevel:
		logger.Error(msg, args...)
	}
}
