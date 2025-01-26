package fv

import (
	"context"
	"errors"
	"net/http"

	"github.com/zzztttkkk/faceless.void/bnd"
	"github.com/zzztttkkk/faceless.void/vld"
)

var (
	ErrNilMarshaler = errors.New("nil marshaler")
)

type HandleFunc func(ctx context.Context, req *http.Request, respw http.ResponseWriter) error

func MakeHandleFunc[Input any, Output any](fnc func(ctx context.Context, params *Input) (Output, error)) HandleFunc {
	return func(ctx context.Context, req *http.Request, respw http.ResponseWriter) error {
		var params Input
		var err error
		if err = bnd.Bind(ctx, &params, req); err != nil {
			return err
		}
		if err = vld.Validate(ctx, &params); err != nil {
			return err
		}
		out, err := fnc(ctx, &params)
		if err != nil {
			return err
		}

		app := AppScope(ctx)
		return app.Marshaler.Marshal(respw, out)
	}
}
