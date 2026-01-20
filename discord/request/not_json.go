package request

import (
	"bytes"
	"context"
	"image"
)

type Image struct {
	baseRequest[image.Image]
	do
}

func NewImage(req REST, method, endpoint string) Image {
	return Image{
		do: newDo(req, method, endpoint),
	}
}

func (r Image) WithBucketID(bucket string) Image {
	r.bucket = bucket
	return r
}

func (r Image) WithData(data any) Image {
	r.data = data
	return r
}

func (r Image) Do(ctx context.Context) (image.Image, error) {
	b, err := r.do.Do(ctx, r.RequestConfig())
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader(b))
	return img, err
}
