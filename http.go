package fv

import (
	"context"
	"net/http"
	"reflect"
)

var (
	ictxType = reflect.TypeOf((*context.Context)(nil)).Elem()
)

type HttpHandlerFunc func(ctx context.Context, req *http.Request, respw http.ResponseWriter) error
