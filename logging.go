package gokord

import (
	"github.com/nyttikord/gokord/logger"
)

// Log the message if level is equal to or greater than Session.LogLevel
func (s *Session) Log(level logger.Level, format string, args ...any) {
	if level < s.LogLevel {
		return
	}

	logger.Log(level, format, args...)
}

func (s *Session) LogError(format string, args ...any) {
	s.Log(logger.LevelError, format, args...)
}

func (s *Session) LogWarn(format string, args ...any) {
	s.Log(logger.LevelWarn, format, args...)
}

func (s *Session) LogInfo(format string, args ...any) {
	s.Log(logger.LevelInfo, format, args...)
}

func (s *Session) LogDebug(format string, args ...any) {
	s.Log(logger.LevelDebug, format, args...)
}

// Log the message if level is equal to or greater than VoiceConnection.LogLevel
func (v *VoiceConnection) Log(level logger.Level, format string, args ...any) {
	if level < v.LogLevel {
		return
	}

	logger.Log(level, format, args...)
}

func (v *VoiceConnection) LogError(format string, args ...any) {
	v.Log(logger.LevelError, format, args...)
}

func (v *VoiceConnection) LogWarn(format string, args ...any) {
	v.Log(logger.LevelWarn, format, args...)
}

func (v *VoiceConnection) LogInfo(format string, args ...any) {
	v.Log(logger.LevelInfo, format, args...)
}

func (v *VoiceConnection) LogDebug(format string, args ...any) {
	v.Log(logger.LevelDebug, format, args...)
}
