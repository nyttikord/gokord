package request

import (
	"context"

	"github.com/nyttikord/gokord/logger"
)

type warnRequest[T any] struct {
	Request[T]
	warn  string
	args  []any
	depth int
}

func (r warnRequest[T]) Do(ctx context.Context) (T, error) {
	getLogger(ctx).WarnContext(logger.NewContext(context.Background(), 1+r.depth), r.warn, r.args...)
	return r.Request.Do(ctx)
}

// WrapWarnDepth sends a warn when the request is executed.
//
// depth specifies the call to log.
// depth = 1 logs the previous call.
// depth = 2 logs the previous of the previous call.
//
// See [WrapWarn] if you are using depth = 1.
func WrapWarnDepth[T any](req Request[T], depth uint, warn string, args ...any) Request[T] {
	return warnRequest[T]{
		Request: req,
		warn:    warn,
		args:    args,
		depth:   int(depth) + 1,
	}
}

// WrapWarn sends a warn when the request is executed.
//
// See [WrapWarnDepth] to specify the depth of the call.
func WrapWarn[T any](req Request[T], warn string, args ...any) Request[T] {
	return WrapWarnDepth(req, 1, warn, args...)
}

type emptyWarnRequest struct {
	Empty
	warn  string
	args  []any
	depth int
}

func (r emptyWarnRequest) Do(ctx context.Context) error {
	getLogger(ctx).WarnContext(logger.NewContext(context.Background(), 1+r.depth), r.warn, r.args...)
	return r.Empty.Do(ctx)
}

// WrapEmptyWarnDepth sends a warn when the request is executed.
//
// depth specifies the call to log.
// depth = 1 logs the previous call.
// depth = 2 logs the previous of the previous call.
//
// See [WrapWarn] if you are using depth = 1.
func WrapEmptyWarnDepth(req Empty, depth uint, warn string, args ...any) Empty {
	return emptyWarnRequest{
		Empty: req,
		warn:  warn,
		args:  args,
		depth: int(depth) + 1,
	}
}

// WrapEmptyWarn sends a warn when the request is executed.
//
// See [WrapWarnDepth] to specify the depth of the call.
func WrapEmptyWarn(req Empty, warn string, args ...any) Empty {
	return WrapEmptyWarnDepth(req, 1, warn, args...)
}
