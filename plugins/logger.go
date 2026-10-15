// Logger which adds dynamic runtime log levels and log types for various game layers and concepts.

package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

type level int

const (
	debugLevel level = iota
	infoLevel
	warnLevel
	errorLevel
)

type logType int

// Separate various game layers and concepts for easier log reading and debugging.
const (
	noLogType logType = iota
	alwaysLogType
	undecidedBehavior
	mainLogic
	stateLogic
	setupLogic
	selectedUnit
	devTool
)

type loggerConfig struct {
	initialLogLevel            slog.Level
	filterTimeAndLevelFromLogs bool
	logUndecidedBehavior       bool
	logMainLogic               bool
	logStateLogic              bool
	logSetupLogic              bool
	logSelectedUnit            bool
	logDevTool                 bool
}

func newLoggerConfig() loggerConfig {
	c := loggerConfig{
		initialLogLevel:            slog.LevelInfo,
		filterTimeAndLevelFromLogs: true,
		logUndecidedBehavior:       true,
		logSetupLogic:              true,
		logDevTool:                 true,
	}
	return c
}

type l struct {
	*slog.Logger
	*slog.LevelVar
	enabledLogTypes map[logType]bool
}

func setupLogger() l {
	lc := newLoggerConfig()

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

	l := l{
		Logger:   slog.New(slog.NewJSONHandler(os.Stderr, logOptions)),
		LevelVar: levelVar,
		enabledLogTypes: map[logType]bool{
			alwaysLogType:     true,
			undecidedBehavior: lc.logUndecidedBehavior,
			mainLogic:         lc.logMainLogic,
			stateLogic:        lc.logStateLogic,
			setupLogic:        lc.logSetupLogic,
			selectedUnit:      lc.logSelectedUnit,
			devTool:           lc.logDevTool,
		},
	}
	return l
}

func (l *l) registerToggles(tr *toggleRegister) {
	tr.register(
		keyL,
		func() { l.toggleLevel() },
	)
	tr.register(
		keyM,
		func() { l.toggleType(mainLogic) },
	)
	tr.register(
		keyS,
		func() { l.toggleType(stateLogic) },
	)
	tr.register(
		keyU,
		func() { l.toggleType(selectedUnit) },
	)
}

// Override the logger state and force log the data.
func (logger l) force(msg string, attrs ...slog.Attr) {
	logger.log(errorLevel, alwaysLogType, msg, attrs...)
}

// Override the logger state to disable all log types except the passed log type.
func (logger l) onlyEnable(logType logType) {
	for k := range logger.enabledLogTypes {
		logger.enabledLogTypes[k] = false
	}
	logger.enabledLogTypes[alwaysLogType] = true
	logger.enabledLogTypes[logType] = true
}

func (logger l) toggleLevel() {
	if logger.LevelVar.Level() == slog.LevelError {
		logger.Set(slog.LevelWarn)
		logger.force("logger.toggleLevel: set to log WARN level")
	} else if logger.LevelVar.Level() == slog.LevelWarn {
		logger.Set(slog.LevelInfo)
		logger.force("logger.toggleLevel: set to log INFO level")
	} else if logger.LevelVar.Level() == slog.LevelInfo {
		logger.Set(slog.LevelDebug)
		logger.force("logger.toggleLevel: set to log DEBUG level")
	} else if logger.LevelVar.Level() == slog.LevelDebug {
		logger.Set(slog.LevelError)
		logger.force("logger.toggleLevel: set to log ERROR level")
	}
}

func (logger l) toggleType(lt logType) {
	var ltName string
	switch lt {
	case noLogType:
		ltName = "noLogType"
	case alwaysLogType:
		logger.force("logger.toggleType: cannot toggle always log type")
		return
	case undecidedBehavior:
		ltName = "undecidedBehavior"
	case mainLogic:
		ltName = "mainLogic"
	case stateLogic:
		ltName = "stateLogic"
	case setupLogic:
		ltName = "setupLogic"
	case selectedUnit:
		ltName = "selectedUnit"
	default:
		logger.force("logger.toggleType: unknown log type")
	}

	enabled := logger.enabledLogTypes[lt]
	logger.enabledLogTypes[lt] = !enabled

	if ltName != "" && enabled {
		logger.force(fmt.Sprintf("logger.toggleType: toggled %s log type off", ltName))
	} else if ltName != "" && !enabled {
		logger.force(fmt.Sprintf("logger.toggleType: toggled %s log type on", ltName))
	}
}

func (logger l) logSelected(
	selected bool, level level, logType logType, msg string, attrs ...slog.Attr) {

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
		logger.LogAttrs(context.Background(), slog.LevelDebug, msg, attrs...)
	case infoLevel:
		logger.LogAttrs(context.Background(), slog.LevelInfo, msg, attrs...)
	case warnLevel:
		logger.LogAttrs(context.Background(), slog.LevelWarn, msg, attrs...)
	case errorLevel:
		logger.LogAttrs(context.Background(), slog.LevelError, msg, attrs...)
	}
}

func nest(key string, attrs []slog.Attr) slog.Attr {
	interfaceSlice := make([]any, len(attrs))
	for k, v := range attrs {
		interfaceSlice[k] = v
	}
	return slog.Group(key, interfaceSlice...)
}
