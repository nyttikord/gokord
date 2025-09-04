package gokord

import (
	"github.com/nyttikord/gokord/logger"
)

// Log the message if level is equal to or greater than Session.LogLevel
//
// See logger.Log
func (s *Session) Log(level logger.Level, caller int, format string, args ...any) {
	if level < s.LogLevel {
		return
	}

	logger.Log(level, caller+1, format, args...)
}

func (s *Session) LogError(format string, args ...any) {
	s.Log(logger.LevelError, 1, format, args...)
}

func (s *Session) LogWarn(format string, args ...any) {
	s.Log(logger.LevelWarn, 1, format, args...)
}

func (s *Session) LogInfo(format string, args ...any) {
	s.Log(logger.LevelInfo, 1, format, args...)
}

func (s *Session) LogDebug(format string, args ...any) {
	s.Log(logger.LevelDebug, 1, format, args...)
}

// Log the message if level is equal to or greater than VoiceConnection.LogLevel
//
// See logger.Log
func (v *VoiceConnection) Log(level logger.Level, caller int, format string, args ...any) {
	if level < v.LogLevel {
		return
	}

	logger.Log(level, caller+1, format, args...)
}

func (v *VoiceConnection) LogError(format string, args ...any) {
	v.Log(logger.LevelError, 1, format, args...)
}

func (v *VoiceConnection) LogWarn(format string, args ...any) {
	v.Log(logger.LevelWarn, 1, format, args...)
}

func (v *VoiceConnection) LogInfo(format string, args ...any) {
	v.Log(logger.LevelInfo, 1, format, args...)
}

func (v *VoiceConnection) LogDebug(format string, args ...any) {
	v.Log(logger.LevelDebug, 1, format, args...)
}
