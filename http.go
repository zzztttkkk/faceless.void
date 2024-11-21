package fv

import (
	"context"
	"net/http"
	"reflect"
)

type IHttpResponse interface {
	Send(ctx context.Context, resp http.ResponseWriter) error
}

type IHttpHandler interface {
	ServeHTTP(ctx context.Context) (IHttpResponse, error)
}

var (
	iHttpHandlerType = reflect.TypeOf((*IHttpHandler)(nil)).Elem()
	ictxType         = reflect.TypeOf((*context.Context)(nil)).Elem()
)

type HttpHandlerFunc func(ctx context.Context) (IHttpResponse, error)

func (fnc HttpHandlerFunc) ServeHTTP(ctx context.Context) (IHttpResponse, error) {
	return fnc(ctx)
}

var _ http.Handler = (*httpEndpoint)(nil)
