package request

import (
	"context"
)

type do struct {
	req      REST
	method   string
	bucket   string
	endpoint string
	data     any
}

func newDo(req REST, method, endpoint string) do {
	return do{
		req:      req,
		method:   method,
		endpoint: endpoint,
		data:     nil,
	}
}

func (r do) Do(ctx context.Context, cfg Config) ([]byte, error) {
	if len(r.bucket) == 0 {
		return r.req.Request(ctx, r.method, r.endpoint, r.data, cfg)
	}
	return r.req.RequestWithBucketID(ctx, r.method, r.endpoint, r.data, r.bucket, cfg)
}

// Simple is a basic request that returns raw bytes.
type Simple struct {
	baseRequest[[]byte]
	do
}

func NewSimple(req REST, method, endpoint string) Simple {
	return Simple{
		do: newDo(req, method, endpoint),
	}
}

func (r Simple) WithBucketID(bucket string) Simple {
	r.bucket = bucket
	return r
}

func (r Simple) WithData(data any) Simple {
	r.data = data
	return r
}

func (r Simple) Do(ctx context.Context) ([]byte, error) {
	return r.do.Do(ctx, r.RequestConfig())
}

type SimpleData[T any] struct {
	baseRequest[T]
	do
}

func NewSimpleData[T any](req REST, method, endpoint string) SimpleData[T] {
	return SimpleData[T]{
		do: newDo(req, method, endpoint),
	}
}

func (r SimpleData[T]) WithBucketID(bucket string) SimpleData[T] {
	r.bucket = bucket
	return r
}

func (r SimpleData[T]) WithData(data any) SimpleData[T] {
	r.data = data
	return r
}

func (r SimpleData[T]) Do(ctx context.Context) (T, error) {
	b, err := r.do.Do(ctx, r.RequestConfig())
	var v T
	if err != nil {
		return v, err
	}
	r.req.Unmarshal(b, &v)
	return v, err
}
