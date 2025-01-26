package fv

import (
	"context"

	"github.com/zzztttkkk/faceless.void/internal"
)

type TAppScope struct {
	Marshaler IMarshaler
}

func AppScope(ctx context.Context) *TAppScope {
	val := ctx.Value(internal.CtxKeyForAppScope)
	if val != nil {
		return val.(*TAppScope)
	}
	panic("nil app scope")
}
