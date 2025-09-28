// Package logger contains the structures used by gokord's custom slog.Logger.
//
// If you want to modify the line logged in the output, see NewContext.
package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"runtime"
	"strings"
	"sync"
	"time"
)

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

// ConsoleHandler represents the default slog.Handler used by gokord.
//
// See New to create a new ConsoleHandler with the given Options.
type ConsoleHandler struct {
	opts Options
	goas []groupOrAttrs
	mu   *sync.Mutex
	out  io.Writer
}

// Options of the ConsoleHandler.
type Options struct {
	// Level reports the minimum level to log.
	// Levels with lower levels are discarded.
	// If nil, the Handler uses [slog.LevelInfo].
	Level slog.Leveler
}

// New creates a new ConsoleHandler.
func New(out io.Writer, opts *Options) *ConsoleHandler {
	h := &ConsoleHandler{out: out, mu: &sync.Mutex{}}
	if opts != nil {
		h.opts = *opts
	}
	if h.opts.Level == nil {
		h.opts.Level = slog.LevelInfo
	}
	return h
}

type key int

const (
	callerSkipKey key = 0
)

// NewContext returns a new context.Context with the callerSkip given.
//
// callerSkip is the number of runtime calls to log before this one.
// 0 is for the current.
// 1 is for the precedent call.
// n is for the n times precedent call.
// The calls to the log is already skipped.
//
// See FromContext to extract the caller from a context.Context.
func NewContext(ctx context.Context, callerSkip int) context.Context {
	return context.WithValue(ctx, callerSkipKey, callerSkip)
}

// FromContext returns the caller in the given context.Context.
//
// See NewContext to create a context.Context.
func FromContext(ctx context.Context) (int, bool) {
	caller, ok := ctx.Value(callerSkipKey).(int)
	return caller, ok
}

// Enabled indicates if the given slog.Level is enabled.
func (h *ConsoleHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

// Handle a slog.Record.
func (h *ConsoleHandler) Handle(ctx context.Context, r slog.Record) error {
	buf := make([]byte, 0, 1024)
	if !r.Time.IsZero() {
		buf = fmt.Appendf(buf, "%s ", r.Time.Format(time.DateTime))
	}
	buf = fmt.Appendf(buf, "[%s] ", r.Level)
	if r.PC != 0 {
		caller, ok := FromContext(ctx)
		var file string
		var line int
		if ok {
			_, file, line, ok = runtime.Caller(caller + 3)
		} else {
			_, file, line, ok = runtime.Caller(3)
		}
		files := strings.Split(file, "/")
		if len(files) == 1 {
			file = files[len(files)-1]
		} else {
			file = files[len(files)-2] + "/" + files[len(files)-1]
		}
		buf = fmt.Appendf(buf, "%s:%d ", file, line)
	}
	if r.Level >= slog.LevelError {
		buf = fmt.Appendf(buf, AnsiRed)
	} else if r.Level >= slog.LevelWarn {
		buf = fmt.Appendf(buf, AnsiYellow)
	}
	buf = fmt.Appendf(buf, "%s%s", r.Message, AnsiReset)
	// Handle state from WithGroup and WithAttrs.
	goas := h.goas
	if r.NumAttrs() == 0 {
		// If the record has no Attrs, remove groups at the end of the list;
		// they are empty.
		for len(goas) > 0 && goas[len(goas)-1].group != "" {
			goas = goas[:len(goas)-1]
		}
	}
	for _, goa := range goas {
		if goa.group != "" {
			buf = fmt.Appendf(buf, " %s={", goa.group)
		} else {
			for _, a := range goa.attrs {
				buf = h.appendAttr(buf, a)
			}
			buf = fmt.Appendf(buf, "!}!") // I don't know where I should put it
		}
	}
	r.Attrs(func(a slog.Attr) bool {
		buf = h.appendAttr(buf, a)
		return true
	})
	buf = fmt.Appendf(buf, "\n")
	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.out.Write(buf)
	return err
}

func (h *ConsoleHandler) appendAttr(buf []byte, a slog.Attr) []byte {
	// Resolve the Attr's value before doing anything else.
	a.Value = a.Value.Resolve()
	// Ignore empty Attrs.
	if a.Equal(slog.Attr{}) {
		return buf
	}
	buf = fmt.Appendf(buf, " ")
	if strings.Contains(a.Key, " ") {
		a.Key = fmt.Sprintf("%q", a.Key)
	}
	switch a.Value.Kind() {
	case slog.KindString:
		buf = fmt.Appendf(buf, "%s=%q", a.Key, a.Value.String())
	case slog.KindTime:
		buf = fmt.Appendf(buf, "%s=%s", a.Key, a.Value.Time().Format(time.RFC3339))
	case slog.KindGroup:
		attrs := a.Value.Group()
		// Ignore empty groups.
		if len(attrs) == 0 {
			return buf
		}
		if a.Key != "" {
			buf = fmt.Appendf(buf, "%s={", a.Key)
		}
		for _, ga := range attrs {
			buf = h.appendAttr(buf, ga)
		}
		if a.Key != "" {
			buf[len(buf)-1] = '}' // replace last space by }
		}
	default:
		var val any
		val = a.Value
		if s, ok := val.(fmt.Stringer); ok {
			val = s.String()
		} else if b, ok := val.([]byte); ok {
			val = string(b)
		}
		buf = fmt.Appendf(buf, "%s=%s", a.Key, a.Value)
	}
	return buf
}

// groupOrAttrs holds either a group name or a list of slog.Attrs.
type groupOrAttrs struct {
	group string      // group name if non-empty
	attrs []slog.Attr // attrs if non-empty
}

func (h *ConsoleHandler) withGroupOrAttrs(goa groupOrAttrs) *ConsoleHandler {
	h2 := *h
	h2.goas = make([]groupOrAttrs, len(h.goas)+1)
	copy(h2.goas, h.goas)
	h2.goas[len(h2.goas)-1] = goa
	return &h2
}

func (h *ConsoleHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	return h.withGroupOrAttrs(groupOrAttrs{group: name})
}

func (h *ConsoleHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	return h.withGroupOrAttrs(groupOrAttrs{attrs: attrs})
}
