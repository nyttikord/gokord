package request

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"

	"github.com/nyttikord/gokord/discord"
)

// File stores info about files you send in messages.
type File struct {
	Name        string
	ContentType string
	Reader      io.Reader
}

// Multipart is multipart body with json request.
// The files field is not immutable.
type Multipart[T any] struct {
	baseRequest[T]
	do    Do
	pre   Pre
	post  Post[T]
	files []*File // files is not immutable.
}

func NewMultipart[T any](req REST, method, endpoint string, data any, files []*File) Multipart[T] {
	base := Multipart[T]{
		do:    newDo(req, method, endpoint),
		files: files,
	}
	base.do.Data = data
	return base
}

func WrapMultipartAsEmpty(r Multipart[[]byte]) Empty {
	return multipartEmpty{multipart: r}
}

func (r Multipart[T]) WithBucketID(bucket string) Multipart[T] {
	r.do.Bucket = bucket
	return r
}

func (r Multipart[T]) WithPre(pre Pre) Multipart[T] {
	r.pre = pre
	return r
}

func (r Multipart[T]) WithPost(post Post[T]) Multipart[T] {
	r.post = post
	return r
}

func (r Multipart[T]) Do(ctx context.Context) (T, error) {
	var v T
	if r.pre != nil {
		err := r.pre(ctx, &r.do)
		if err != nil {
			return v, err
		}
	}
	if r.do.Data == nil || len(r.files) == 0 {
		panic("invalid multipart body: data or files nil")
	}
	contentType, body, err := multipartBodyWithJSON(r.do.Data, r.files)
	if err != nil {
		return v, err
	}
	bucket := r.do.Endpoint
	if r.do.Bucket != "" {
		bucket = r.do.Bucket
	}
	b, err := r.do.req.RequestRaw(
		ctx, http.MethodPatch, r.do.Endpoint, contentType, body, bucket, 0, r.Config(),
	)
	if err != nil {
		return v, err
	}
	if r.post != nil {
		return r.post(ctx, b)
	}
	err = r.do.req.Unmarshal(b, v)
	if err != nil {
		return v, err
	}
	return v, nil
}

// MultipartBodyWithJSON returns the contentType and body for a discord request.
//
// data is the object to encode for payload_json in the multipart request.
// files is the files to include in the request.
func multipartBodyWithJSON(data any, files []*File) (requestContentType string, requestBody []byte, err error) {
	body := &bytes.Buffer{}
	bodywriter := multipart.NewWriter(body)

	payload, err := json.Marshal(data)
	if err != nil {
		return
	}

	var p io.Writer

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="payload_json"`)
	h.Set("Content-Type", "application/json")

	p, err = bodywriter.CreatePart(h)
	if err != nil {
		return
	}

	if _, err = p.Write(payload); err != nil {
		return
	}

	for i, file := range files {
		h := make(textproto.MIMEHeader)
		h.Set(
			"Content-Disposition",
			fmt.Sprintf(`form-data; name="files[%d]"; filename="%s"`, i, discord.QuoteEscaper.Replace(file.Name)),
		)
		contentType := file.ContentType
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		h.Set("Content-Type", contentType)

		p, err = bodywriter.CreatePart(h)
		if err != nil {
			return
		}

		if _, err = io.Copy(p, file.Reader); err != nil {
			return
		}
	}

	err = bodywriter.Close()
	if err != nil {
		return
	}

	return bodywriter.FormDataContentType(), body.Bytes(), nil
}

type multipartEmpty struct {
	multipart Multipart[[]byte]
	cfg       Config
}

func (r multipartEmpty) Do(ctx context.Context) error {
	_, err := r.multipart.WithPost(func(ctx context.Context, b []byte) ([]byte, error) { return b, nil }).Do(ctx)
	return err
}

func (r multipartEmpty) WithRetryOnRateLimit(b bool) Empty {
	r.cfg.Header = r.cfg.Header.Clone()

	r.cfg.ShouldRetryOnRateLimit = &b
	return r
}

func (r multipartEmpty) WithRestRetries(m uint) Empty {
	r.cfg.Header = r.cfg.Header.Clone()

	r.cfg.MaxRestRetries = &m
	return r
}

func (r multipartEmpty) WithHeader(key, value string) Empty {
	r.cfg.Header = r.cfg.Header.Clone()

	r.cfg.Header.Set(key, value)
	return r
}

func (r multipartEmpty) WithAuditLogReason(reason string) Empty {
	return r.WithHeader("X-Audit-Log-Reason", reason)
}

func (r multipartEmpty) WithLocale(locale discord.Locale) Empty {
	return r.WithHeader("X-Discord-Locale", string(locale))
}

func (r multipartEmpty) Config() Config {
	return r.cfg
}
