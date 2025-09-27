package gokord

import (
	"fmt"
	"log/slog"

	"github.com/nyttikord/gokord/logger"
)

type stdLogger struct {
	Level slog.Level
}

func (s *stdLogger) Log(level slog.Level, caller int, format string, args ...any) {
	logger.Log(level, caller+1, format, args...)
}

func (s *stdLogger) LogError(err error, format string, args ...any) {
	format += fmt.Sprintf(" %s%s%s", logger.AnsiRed, err.Error(), logger.AnsiReset)
	s.Log(slog.LevelError, 1, format, args...)
}

func (s *stdLogger) LogWarn(format string, args ...any) {
	format = fmt.Sprintf("%s%s%s ", logger.AnsiYellow, format, logger.AnsiReset)
	s.Log(slog.LevelWarn, 1, format, args...)
}

func (s *stdLogger) LogInfo(format string, args ...any) {
	s.Log(slog.LevelInfo, 1, format, args...)
}

func (s *stdLogger) LogDebug(format string, args ...any) {
	s.Log(slog.LevelDebug, 1, format, args...)
}

func (s *stdLogger) GetLevel() slog.Level {
	return s.Level
}

func (s *stdLogger) ChangeLevel(level slog.Level) {
	s.Level = level
}
