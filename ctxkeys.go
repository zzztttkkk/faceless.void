package fv

import (
	"context"
	"net/http"
)

type ctxKeyType int

const (
	ctxKeyForHttpMarshaler = ctxKeyType(iota)
	ctxKeyForHttpRequest
	ctxKeyForBindingGetter
)

func HttpRequest(ctx context.Context) *http.Request {
	av := ctx.Value(ctxKeyForHttpRequest)
	if av != nil {
		return av.(*http.Request)
	}
	panic("empty http.Request")
}

func BindingGetter(ctx context.Context) *_Getter {
	av := ctx.Value(ctxKeyForBindingGetter)
	if av != nil {
		return av.(*_Getter)
	}
	panic("empty BindingHelper")
}
