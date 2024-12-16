package fv

import (
	"context"
	"net/http"

	"github.com/zzztttkkk/faceless.void/internal"
)

func HttpRequest(ctx context.Context) *http.Request {
	av := ctx.Value(internal.CtxKeyForHttpRequest)
	if av != nil {
		return av.(*http.Request)
	}
	panic("empty http.Request")
}

func WithLanguage(ctx context.Context, lang string) context.Context {
	return context.WithValue(ctx, internal.CtxKeyForLanguageKind, lang)
}
