package request

import (
	"bytes"
	"context"
	"image"
)

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
