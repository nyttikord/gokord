package gokord

import (
	"fmt"

	"github.com/nyttikord/gokord/logger"
)

type stdLogger struct {
	Level logger.Level
}

func (s *stdLogger) Log(level logger.Level, caller int, format string, args ...any) {
	if level < s.Level {
		return
	}

	logger.Log(level, caller+1, format, args...)
}

func (s *stdLogger) LogError(err error, format string, args ...any) {
	format += fmt.Sprintf(" %s%s%s", logger.AnsiRed, err.Error(), logger.AnsiReset)
	s.Log(logger.LevelError, 1, format, args...)
}

func (s *stdLogger) LogWarn(format string, args ...any) {
	format = fmt.Sprintf("%s%s%s ", logger.AnsiYellow, format, logger.AnsiReset)
	s.Log(logger.LevelWarn, 1, format, args...)
}

func (s *stdLogger) LogInfo(format string, args ...any) {
	s.Log(logger.LevelInfo, 1, format, args...)
}

func (s *stdLogger) LogDebug(format string, args ...any) {
	s.Log(logger.LevelDebug, 1, format, args...)
}

func (s *stdLogger) GetLevel() logger.Level {
	return s.Level
}

func (s *stdLogger) ChangeLevel(level logger.Level) {
	s.Level = level
}
