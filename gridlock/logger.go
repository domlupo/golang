// logger wrapper which adds runtime dynamic log level and log types (specific game logic e.g. state)
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

type level int
type logType int

const (
	debugLevel level = iota
	infoLevel
	warnLevel
	errorLevel
)

const (
	noLogType logType = iota
	alwaysLogType
	undecidedBehavior
	mainLogType
	stateLogType
	selectedUnitLogType
)

type loggerConfig struct {
	initialLogLevel            slog.Level
	filterTimeAndLevelFromLogs bool
	logLoggerLogs              bool
	logUndecidedBehavior       bool
	logMainLogs                bool
	logStateLogs               bool
	logSelectedUnitLogs        bool
}

type l struct {
	*slog.Logger
	*slog.LevelVar
	enabledLogTypes map[logType]bool
}

func newLogger(lc loggerConfig) l {
	var levelVar = new(slog.LevelVar)
	levelVar.Set(lc.initialLogLevel)
	logOptions := &slog.HandlerOptions{
		Level: levelVar,
		ReplaceAttr: func(groups []string, attribute slog.Attr) slog.Attr {
			if lc.filterTimeAndLevelFromLogs && (attribute.Key == "time" || attribute.Key == "level") {
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
			alwaysLogType:       lc.logLoggerLogs,
			undecidedBehavior:   lc.logUndecidedBehavior,
			mainLogType:         lc.logMainLogs,
			stateLogType:        lc.logStateLogs,
			selectedUnitLogType: lc.logSelectedUnitLogs},
	}

	return l
}

func (logger l) toggleLevel() {
	if logger.LevelVar.Level() == slog.LevelError {
		logger.Set(slog.LevelWarn)
		logger.log(errorLevel, alwaysLogType, "logger set to log WARN logs")
	} else if logger.LevelVar.Level() == slog.LevelWarn {
		logger.Set(slog.LevelInfo)
		logger.log(errorLevel, alwaysLogType, "logger set to log INFO logs")
	} else if logger.LevelVar.Level() == slog.LevelInfo {
		logger.Set(slog.LevelDebug)
		logger.log(errorLevel, alwaysLogType, "logger set to log DEBUG logs")
	} else if logger.LevelVar.Level() == slog.LevelDebug {
		logger.Set(slog.LevelError)
		logger.log(errorLevel, alwaysLogType, "logger set to log ERROR logs")
	}
}

func (logger l) toggleType(lt logType) {
	enabled := logger.enabledLogTypes[lt]
	logger.enabledLogTypes[lt] = !enabled

	var ltName string
	switch lt {
	case noLogType:
		ltName = "noLogType"
	case alwaysLogType:
		ltName = "alwaysLogType"
	case undecidedBehavior:
		ltName = "undecidedBehavior"
	case mainLogType:
		ltName = "mainLog"
	case stateLogType:
		ltName = "stateLog"
	case selectedUnitLogType:
		ltName = "selectedUnitLog"
	default:
		logger.log(errorLevel, alwaysLogType, "Did not find a matching log type name")
	}

	if ltName != "" && enabled {
		logger.log(errorLevel, alwaysLogType, fmt.Sprintf("toggled log type off: %s", ltName))
	} else if ltName != "" && !enabled {
		logger.log(errorLevel, alwaysLogType, fmt.Sprintf("toggled log type on: %s", ltName))
	}
}

func (logger l) logSelected(selected bool, level level, logType logType, msg string, attrs ...slog.Attr) {
	if selected {
		logger.log(level, logType, msg, attrs...)
	}
}

func (logger l) log(level level, logType logType, msg string, attrs ...slog.Attr) {
	if !logger.enabledLogTypes[logType] {
		return
	}

	switch level {
	case debugLevel:
		logger.LogAttrs(context.TODO(), slog.LevelDebug, msg, attrs...)
	case infoLevel:
		logger.LogAttrs(context.TODO(), slog.LevelInfo, msg, attrs...)
	case warnLevel:
		logger.LogAttrs(context.TODO(), slog.LevelWarn, msg, attrs...)
	case errorLevel:
		logger.LogAttrs(context.TODO(), slog.LevelError, msg, attrs...)
	}
}
