package request

import (
	"context"
)

// Simple is a basic request that returns nothing if there is no error
type Simple struct {
	baseRequest[[]byte]
	req      REST
	method   string
	bucket   string
	endpoint string
	data     any
}

func NewSimple(req REST, method, endpoint string) Simple {
	return Simple{
		req:      req,
		method:   method,
		endpoint: endpoint,
		data:     nil,
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
	if len(r.bucket) == 0 {
		return r.req.Request(ctx, r.method, r.endpoint, r.data, r.RequestConfig())
	}
	return r.req.RequestWithBucketID(ctx, r.method, r.endpoint, r.data, r.bucket, r.RequestConfig())
}

type SimpleData[T any] struct {
	baseRequest[T]
	req      REST
	method   string
	bucket   string
	endpoint string
	data     any
}

func NewSimpleData[T any](req REST, method, endpoint string) SimpleData[T] {
	return SimpleData[T]{
		req:      req,
		method:   method,
		endpoint: endpoint,
		data:     nil,
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
	var err error
	var b []byte
	if len(r.bucket) == 0 {
		b, err = r.req.Request(ctx, r.method, r.endpoint, r.data, r.RequestConfig())
	} else {
		b, err = r.req.RequestWithBucketID(ctx, r.method, r.endpoint, r.data, r.bucket, r.RequestConfig())
	}
	var v T
	r.req.Unmarshal(b, &v)
	return v, err
}
