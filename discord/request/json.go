package request

import (
	"context"
)

type Do struct {
	Method   string
	Bucket   string
	Endpoint string
	Data     any
}

func newDo(method, endpoint string) Do {
	return Do{
		Method:   method,
		Endpoint: endpoint,
		Data:     nil,
	}
}

func (r Do) do(ctx context.Context, cfg Config) ([]byte, error) {
	req := getREST(ctx)
	if len(r.Bucket) == 0 {
		return req.Request(ctx, r.Method, r.Endpoint, r.Data, cfg)
	}
	return req.RequestWithBucketID(ctx, r.Method, r.Endpoint, r.Data, r.Bucket, cfg)
}

// Simple is a basic [Request] that returns raw bytes.
type Simple struct {
	baseRequest[[]byte]
	do Do
}

func NewSimple(method, endpoint string) Simple {
	return Simple{
		do: newDo(method, endpoint),
	}
}

func (r Simple) WithBucketID(bucket string) Simple {
	r.do.Bucket = bucket
	return r
}

func (r Simple) WithData(data any) Simple {
	r.do.Data = data
	return r
}

func (r Simple) Do(ctx context.Context) ([]byte, error) {
	return r.do.do(ctx, r.Config())
}

type Data[T any] struct {
	baseRequest[T]
	do  Do
	pre Pre
}

func NewData[T any](method, endpoint string) Data[T] {
	return Data[T]{
		do: newDo(method, endpoint),
	}
}

func (r Data[T]) WithBucketID(bucket string) Data[T] {
	r.do.Bucket = bucket
	return r
}

func (r Data[T]) WithData(data any) Data[T] {
	r.do.Data = data
	return r
}

func (r Data[T]) WithPre(pre Pre) Data[T] {
	r.pre = pre
	return r
}

func (r Data[T]) Do(ctx context.Context) (T, error) {
	var v T
	if r.pre != nil {
		err := r.pre(ctx, &r.do)
		if err != nil {
			return v, err
		}
	}
	b, err := r.do.do(ctx, r.Config())
	if err != nil {
		return v, err
	}
	err = Unmarshal(ctx, b, &v)
	return v, err
}
