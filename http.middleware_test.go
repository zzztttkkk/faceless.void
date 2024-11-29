package fv

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

func TestMiddleware(t *testing.T) {
	middleware := []HttpMiddlewareFunc{
		func(ctx context.Context, next func(context.Context) error, req *http.Request, respw http.ResponseWriter) error {
			fmt.Println("A", ctx.Value("M"))
			defer fmt.Println("D-A")
			return next(context.WithValue(ctx, "M", "A"))
		},
		func(ctx context.Context, next func(context.Context) error, req *http.Request, respw http.ResponseWriter) error {
			fmt.Println("B", ctx.Value("M"))
			defer fmt.Println("D-B")
			return next(context.WithValue(ctx, "M", "B"))
		},
	}

	fnc := wrapMiddleware(
		func(ctx context.Context, req *http.Request, respw http.ResponseWriter) error {
			fmt.Println(">>>", ctx.Value("M"))
			return nil
		}, middleware)

	fnc(context.Background(), nil, nil)
}
