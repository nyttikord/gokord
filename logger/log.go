// Package logger provides a general Logger interface to log information and useful constant to color messages.
package logger

import (
	"context"
	"log/slog"
	"os"
)

// Logger is an interface describing a custom logger implementation.
type Logger interface {
	// Log the message if level is equal to or greater than previously defined
	//
	// See logger.Log
	Log(level slog.Level, caller int, format string, args ...any)
	LogError(err error, format string, args ...any)
	LogWarn(format string, args ...any)
	LogInfo(format string, args ...any)
	LogDebug(format string, args ...any)
	// GetLevel returns the minimum Level logged
	GetLevel() slog.Level
	// ChangeLevel changes the minimum Level logged
	ChangeLevel(level slog.Level)
}

const (
	AnsiReset       = "\033[0m"
	AnsiRed         = "\033[91m"
	AnsiGreen       = "\033[32m"
	AnsiYellow      = "\033[33m"
	AnsiBlue        = "\033[34m"
	AnsiMagenta     = "\033[35m"
	AnsiCyan        = "\033[36m"
	AnsiWhite       = "\033[37m"
	AnsiBlueBold    = "\033[34;1m"
	AnsiMagentaBold = "\033[35;1m"
	AnsiRedBold     = "\033[31;1m"
	AnsiYellowBold  = "\033[33;1m"
)

var test = New(os.Stdout, &Options{Level: slog.LevelInfo})

// Log logs and formats a message at the given level.
//
// Caller is the number of calls before this one (e.g., 0 if you want to log this call, 1 to log the call before...)
func Log(level slog.Level, caller int, format string, args ...any) {
	//pc, file, line, _ := runtime.Caller(caller + 1)
	//
	//files := strings.Split(file, "/")
	//file = files[len(files)-1]
	//
	//name := runtime.FuncForPC(pc).Name()
	//fns := strings.Split(name, ".")
	//name = fns[len(fns)-1]

	t := slog.New(test)
	t.Log(context.Background(), level, format, args)

	//log.Printf("[%s] %s:%d:%s() %s\n", level.String(), file, line, name, fmt.Sprintf(format, args...))
}
