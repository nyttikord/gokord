package request

import (
	"bytes"
	"context"
	"image"
	_ "image/jpeg" // For JPEG decoding
	_ "image/png"  // For PNG decoding
)

// Custom is a request with custom treatment after bytes received.
type Custom[T any] struct {
	baseRequest[T]
	do   Do
	pre  Pre
	post Post[T]
}

func NewCustom[T any](req REST, method, endpoint string) Custom[T] {
	return Custom[T]{
		do: newDo(req, method, endpoint),
	}
}

func (r Custom[T]) WithBucketID(bucket string) Custom[T] {
	r.do.Bucket = bucket
	return r
}

func (r Custom[T]) WithData(data any) Custom[T] {
	r.do.Data = data
	return r
}

func (r Custom[T]) WithPre(pre Pre) Custom[T] {
	r.pre = pre
	return r
}

func (r Custom[T]) WithPost(post Post[T]) Custom[T] {
	r.post = post
	return r
}

func (r Custom[T]) Do(ctx context.Context) (T, error) {
	var v T
	if r.pre != nil {
		err := r.pre(ctx, &r.do)
		if err != nil {
			return v, err
		}
	}
	b, err := r.do.Do(ctx, r.Config())
	if err != nil {
		return v, err
	}
	if r.post == nil {
		panic("invalid Custom request: post is nil")
	}
	return r.post(ctx, b)
}

// Image is a request that returns an image.Image.
type Image struct {
	baseRequest[image.Image]
	do  Do
	pre Pre
}

func NewImage(req REST, method, endpoint string) Image {
	return Image{
		do: newDo(req, method, endpoint),
	}
}

func (r Image) WithBucketID(bucket string) Image {
	r.do.Bucket = bucket
	return r
}

func (r Image) WithData(data any) Image {
	r.do.Data = data
	return r
}

func (r Image) WithPre(pre Pre) Image {
	r.pre = pre
	return r
}

func (r Image) Do(ctx context.Context) (image.Image, error) {
	if r.pre != nil {
		err := r.pre(ctx, &r.do)
		if err != nil {
			return nil, err
		}
	}
	b, err := r.do.Do(ctx, r.Config())
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader(b))
	return img, err
}
