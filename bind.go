package fv

import (
	"context"
	"net/http"
	"reflect"
)

type Empty struct{}

var (
	bindfncs  = map[reflect.Type]func(ctx context.Context, dest any) error{}
	bindtypes = map[reflect.Type]Empty{}
)

func bindHttp(ctx context.Context, req *http.Request, vtype reflect.Type, dest any) error {
	return nil
}

type bindingOptions struct {
	rawtype reflect.Type
	valtype reflect.Type
	src     string
}
