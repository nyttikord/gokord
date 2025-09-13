package logger

import (
	"fmt"
	"log"
	"runtime"
	"strings"
)

type Level int

const (
	LevelDebug Level = 0
	LevelInfo  Level = 1
	LevelWarn  Level = 2
	LevelError Level = 3
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	}
	return ""
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

// Log logs and formats a message at the given level.
//
// Caller is the number of calls before this one (e.g., 0 if you want to log this call, 1 to log the call before...)
func Log(level Level, caller int, format string, args ...any) {
	pc, file, line, _ := runtime.Caller(caller + 1)

	files := strings.Split(file, "/")
	file = files[len(files)-1]

	name := runtime.FuncForPC(pc).Name()
	fns := strings.Split(name, ".")
	name = fns[len(fns)-1]

	log.Printf("[%s] %s:%d:%s() %s\n", level.String(), file, line, name, fmt.Sprintf(format, args...))
}
