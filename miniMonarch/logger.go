package main

import (
	"log/slog"
	"os"
)

const initialLogLevel = slog.LevelInfo
const filterTimeAndLevelFromLogs = true

type l struct {
	*slog.Logger
	*slog.LevelVar
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
