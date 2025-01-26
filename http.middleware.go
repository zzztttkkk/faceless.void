package fv

import (
	"context"
	"net/http"

	"github.com/zzztttkkk/faceless.void/internal"
)

type HttpMiddlewareFunc func(ctx context.Context, next func(context.Context) error, req *http.Request, respw http.ResponseWriter) error

func wrapMiddleware(fnc HandleFunc, middleware []HttpMiddlewareFunc) HandleFunc {
	mc := len(middleware)
	return func(ctx context.Context, req *http.Request, respw http.ResponseWriter) error {
		var idx = -1
		var next func(context.Context) error
		next = func(ctx context.Context) error {
			idx++
			if idx >= mc {
				return fnc(ctx, req, respw)
			}
			return middleware[idx](ctx, next, req, respw)
		}
		return next(ctx)
	}
}

type middlewareBuilder struct {
	pairs []internal.Pair[string]
}

func (builder *middlewareBuilder) Filename(filename string) *middlewareBuilder {
	builder.pairs = append(builder.pairs, internal.PairOf("filename", filename))
	return builder
}

func (builder *middlewareBuilder) Func(fnc HttpMiddlewareFunc) *middlewareBuilder {
	builder.pairs = append(builder.pairs, internal.PairOf("fnc", fnc))
	return builder
}

var (
	allMiddlewareBuilders []*middlewareBuilder
)

func (builder *middlewareBuilder) Register() {
	allMiddlewareBuilders = append(allMiddlewareBuilders, builder)
}

func Middleware() *middlewareBuilder {
	return &middlewareBuilder{}
}
