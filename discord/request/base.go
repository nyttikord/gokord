// Package request contains utils to create and process requests.
//
// Request[T] is a general interface representing any requests returning something.
// Methods should always use this one instead of a more precise returns type.
//
// Empty is a simple request returning nothing (only an error).
// You wan wrap any Simple request as an EmptyRequest with WrapAsEmpty.
// You can unwrap it with UnwrapEmpty.
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
	// RequestConfig returns the Config used
	RequestConfig() Config
}

// Empty is a Request that only returns an error when it is executed.
type Empty struct {
	Simple
}

func (r Empty) Do(ctx context.Context) error {
	_, err := r.Simple.Do(ctx)
	return err
}

func WrapAsEmpty(req Simple) Empty {
	return Empty{req}
}

func UnwrapEmpty(req Empty) Simple {
	return req.Simple
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

func (r baseRequest[T]) RequestConfig() Config {
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
