package fv

import (
	"bytes"
	"context"
	"net/http"
	"sync"
)

type errorResponse struct {
	err error
}

// Send implements IHttpResponse.
func (e *errorResponse) Send(ctx context.Context, resp http.ResponseWriter) error {
	code := http.StatusInternalServerError
	var msg []byte
	he, ok := e.err.(IHttpError)
	if ok {
		code = he.StatusCode()
		msg = he.BodyMessage(ctx)
	}
	resp.WriteHeader(code)
	_, err := resp.Write(msg)
	return err
}

var _ IHttpResponse = (*errorResponse)(nil)

type codeResponse int64

// Send implements IHttpResponse.
func (c codeResponse) Send(ctx context.Context, resp http.ResponseWriter) error {
	resp.WriteHeader(int(c))
	_, err := resp.Write(nil)
	return err
}

var _ IHttpResponse = (codeResponse)(0)

type anyResponse struct {
	val any
}

var (
	bufpool = sync.Pool{
		New: func() any {
			return bytes.NewBuffer(nil)
		},
	}
)

// Send implements IHttpResponse.
func (a anyResponse) Send(ctx context.Context, resp http.ResponseWriter) error {
	marshaler := getHttpMarshal(ctx)
	resp.Header().Add("Content-Type", marshaler.ContentType())

	buf := bufpool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		bufpool.Put(buf)
	}()

	err := marshaler.Marshal(a.val, buf)
	if err != nil {
		return err
	}
	_, err = resp.Write(buf.Bytes())
	return err
}

var _ IHttpResponse = anyResponse{val: nil}

type FileResponse struct {
	path   string
	nocopy bool
}

// Send implements IHttpResponse.
func (f *FileResponse) Send(ctx context.Context, resp http.ResponseWriter) error {
	panic("unimplemented")
}

var _ IHttpResponse = (*FileResponse)(nil)
