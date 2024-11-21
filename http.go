package fv

import (
	"context"
	"net/http"
	"reflect"
)

type IHttpHandler interface {
	ServeHTTP(ctx context.Context, req *http.Request, respw http.ResponseWriter) error
}

var (
	iHttpHandlerType = reflect.TypeOf((*IHttpHandler)(nil)).Elem()
	ictxType         = reflect.TypeOf((*context.Context)(nil)).Elem()
)

type HttpHandlerFunc func(ctx context.Context, req *http.Request, respw http.ResponseWriter) error

func (fnc HttpHandlerFunc) ServeHTTP(ctx context.Context, req *http.Request, respw http.ResponseWriter) error {
	return fnc(ctx, req, respw)
}

var _ http.Handler = (*httpEndpoint)(nil)
