// Package request contains utils to create and process requests.
//
// Request[T] is a general interface representing any requests returning something.
// Methods should always use this one instead of a more precise returns type.
//
// Empty is a simple request returning nothing (only an error).
// You wan wrap any Simple request as an EmptyRequest with WrapAsEmpty.
package request

import (
	"context"
	"maps"
	"net/http"

	"github.com/nyttikord/gokord/discord"
)

// Request represents an HTTP request.
// It is immutable: when you call a method, it returns a new Request with this option.
type Request[T any] interface {
	// Do executes the request.
	Do(context.Context) (T, error)
	// WithRetryOnRatelimit controls whether the session should retry the request on rate limit.
	WithRetryOnRateLimit(bool) Request[T]
	// WithRestRetries changes maximum amount of retries if request fails.
	WithRestRetries(uint) Request[T]
	// WithHeader sets a header in the request.
	WithHeader(string, string) Request[T]
	// WithAuditLogReason changes audit log reason associated with the request.
	WithAuditLogReason(string) Request[T]
	// WithLocale changes accepted locale of the request.
	WithLocale(discord.Locale) Request[T]
	// Config returns the Config used
	Config() Config
}

// Empty represents an HTTP request that only returns an error when it is executed.
// It is immutable: when you call a method, it returns a new Request with this option.
type Empty interface {
	// Do executes the request.
	Do(context.Context) error
	// WithRetryOnRatelimit controls whether the session should retry the request on rate limit.
	WithRetryOnRateLimit(bool) Empty
	// WithRestRetries changes maximum amount of retries if request fails.
	WithRestRetries(uint) Empty
	// WithHeader sets a header in the request.
	WithHeader(string, string) Empty
	// WithAuditLogReason changes audit log reason associated with the request.
	WithAuditLogReason(string) Empty
	// WithLocale changes accepted locale of the request.
	WithLocale(discord.Locale) Empty
	// Config returns the Config used
	Config() Config
}

type simpleEmpty struct {
	simple Simple
	err    error
	cfg    Config
}

func (r simpleEmpty) Do(ctx context.Context) error {
	if r.err != nil {
		return r.err
	}
	_, err := r.simple.Do(ctx)
	return err
}

func (r simpleEmpty) WithRetryOnRateLimit(b bool) Empty {
	r.cfg.Header = r.cfg.Header.Clone()

	r.cfg.ShouldRetryOnRateLimit = &b
	return r
}

func (r simpleEmpty) WithRestRetries(m uint) Empty {
	r.cfg.Header = r.cfg.Header.Clone()

	r.cfg.MaxRestRetries = &m
	return r
}

func (r simpleEmpty) WithHeader(key, value string) Empty {
	r.cfg.Header = r.cfg.Header.Clone()

	r.cfg.Header.Set(key, value)
	return r
}

func (r simpleEmpty) WithAuditLogReason(reason string) Empty {
	return r.WithHeader("X-Audit-Log-Reason", reason)
}

func (r simpleEmpty) WithLocale(locale discord.Locale) Empty {
	return r.WithHeader("X-Discord-Locale", string(locale))
}

func (r simpleEmpty) Config() Config {
	return r.cfg
}

func WrapAsEmpty(req Simple) Empty {
	return simpleEmpty{simple: req}
}

func WrapErrorAsEmpty(err error) Empty {
	return simpleEmpty{err: err}
}

type Config struct {
	Header                 http.Header
	ShouldRetryOnRateLimit *bool
	MaxRestRetries         *uint
}

func NewConfig() Config {
	return Config{Header: make(http.Header)}
}

type baseRequest[T any] Config

func (r baseRequest[T]) Do(ctx context.Context) (T, error) {
	panic("cannot execute a baseRequest")
}

func (r baseRequest[T]) WithRetryOnRateLimit(b bool) Request[T] {
	dst := make(http.Header, len(r.Header))
	maps.Copy(dst, r.Header)
	r.Header = dst

	r.ShouldRetryOnRateLimit = &b
	return r
}

func (r baseRequest[T]) WithRestRetries(m uint) Request[T] {
	dst := make(http.Header, len(r.Header))
	maps.Copy(dst, r.Header)
	r.Header = dst

	r.MaxRestRetries = &m
	return r
}

func (r baseRequest[T]) WithHeader(key, value string) Request[T] {
	dst := make(http.Header, len(r.Header))
	maps.Copy(dst, r.Header)
	r.Header = dst

	r.Header.Set(key, value)
	return r
}

func (r baseRequest[T]) WithAuditLogReason(reason string) Request[T] {
	return r.WithHeader("X-Audit-Log-Reason", reason)
}

func (r baseRequest[T]) WithLocale(locale discord.Locale) Request[T] {
	return r.WithHeader("X-Discord-Locale", string(locale))
}

func (r baseRequest[T]) Config() Config {
	return Config(r)
}

// Error is a request that returns the specified error when Do is called.
type Error[T any] struct {
	baseRequest[T]
	err error
}

func NewError[T any](err error) Error[T] {
	if err == nil {
		panic("cannot use nil error in request.Error")
	}
	return Error[T]{err: err}
}

func (r Error[T]) Do(ctx context.Context) (T, error) {
	var v T
	return v, r.err
}

// Pre is a function called before the request.
type Pre func(context.Context, *Do) error

// Post is function called after the request.
type Post[T any] func(context.Context, []byte) (T, error)
