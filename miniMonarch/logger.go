// logger wrapper which adds runtime dynamic log level and log types (specific game logic e.g. state)
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
	loggerLogType
	mainLogType
	stateLogType
	selectedUnitLogType
)

// config
const initialLogLevel = slog.LevelDebug
const filterTimeAndLevelFromLogs = true
const logLoggerLogs = true
const logMainLogs = true
const logStateLogs = false
const logSelectedUnitLogs = true

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
		Logger:   logger,
		LevelVar: levelVar,
		enabledLogTypes: map[logType]bool{
			loggerLogType:       logLoggerLogs,
			mainLogType:         logMainLogs,
			stateLogType:        logStateLogs,
			selectedUnitLogType: logSelectedUnitLogs},
	}

	return l
}

func (logger l) toggleLevel() {
	if logger.LevelVar.Level() == slog.LevelError {
		logger.Set(slog.LevelWarn)
		logger.log(errorLevel, loggerLogType, "logger set to log WARN logs")
	} else if logger.LevelVar.Level() == slog.LevelWarn {
		logger.Set(slog.LevelInfo)
		logger.log(errorLevel, loggerLogType, "logger set to log INFO logs")
	} else if logger.LevelVar.Level() == slog.LevelInfo {
		logger.Set(slog.LevelDebug)
		logger.log(errorLevel, loggerLogType, "logger set to log DEBUG logs")
	} else if logger.LevelVar.Level() == slog.LevelDebug {
		logger.Set(slog.LevelError)
		logger.log(errorLevel, loggerLogType, "logger set to log ERROR logs")
	}
}

func (logger l) toggleType(logType logType) {
	enabled := logger.enabledLogTypes[logType]
	// TODO: Should I WARN log what log type is enabled and disabled?
	logger.enabledLogTypes[logType] = !enabled
}

func (logger l) logSelected(selected bool, level level, logType logType, msg string, args ...any) {
	if selected {
		logger.log(level, logType, msg, args...)
	}
}

func (logger l) log(level level, logType logType, msg string, args ...any) {
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
