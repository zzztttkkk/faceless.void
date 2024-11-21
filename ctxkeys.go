package fv

import (
	"context"
	"net/http"
)

type ctxKeyType int

const (
	ctxKeyForHttpMarshaler = ctxKeyType(iota)
	ctxKeyForHttpRequest
)

func HttpRequest(ctx context.Context) *http.Request {
	av := ctx.Value(ctxKeyForHttpRequest)
	if av != nil {
		return av.(*http.Request)
	}
	panic("empty http.Request")
}
