package fv

import (
	"context"
	"encoding/json"
	"io"
	"unsafe"

	"github.com/zzztttkkk/faceless.void/internal"
)

type _BindingAnyField[T any] struct {
	ins         *_BindingInstance
	ptr         *T
	optional    bool
	unmarshaler func([]byte, *T) error
	validator   func(context.Context, *T) error
}

func AnyField[T any](ins *_BindingInstance, ptr *T) *_BindingAnyField[T] {
	af := &_BindingAnyField[T]{
		ins: ins,
		ptr: ptr,
	}
	ins.fields = append(ins.fields, af.do)
	return af
}

func (af *_BindingAnyField[T]) Optional() *_BindingAnyField[T] {
	af.optional = true
	return af
}
func (af *_BindingAnyField[T]) Unmarshal(fnc func([]byte, *T) error) *_BindingAnyField[T] {
	af.unmarshaler = fnc
	return af
}

func (af *_BindingAnyField[T]) Validate(fnc func(context.Context, *T) error) *_BindingAnyField[T] {
	af.validator = fnc
	return af
}

var (
	msgForBindingUnmarshalFailed = internal.NewI18nString("fv.binding: unmarshal failed, %s, %s")
)

func (af *_BindingAnyField[T]) do(ctx context.Context) error {
	req := HttpRequest(ctx)
	defer req.Body.Close()
	payload, err := io.ReadAll(req.Body)
	if err != nil || len(payload) < 1 {
		if af.optional {
			return nil
		}
		name := af.ins.nameof(unsafe.Pointer(af.ptr))
		return internal.NewError(ctx, internal.ErrorKind(ErrorKindBindingMissingRequired), msgForBindingMissingRequired, name)
	}
	if af.unmarshaler == nil {
		err = json.Unmarshal(payload, af.ptr)
	} else {
		err = af.unmarshaler(payload, af.ptr)
	}
	if err != nil {
		name := af.ins.nameof(unsafe.Pointer(af.ptr))
		return internal.NewError(ctx, internal.ErrorKind(ErrorKindBindingUnmarshalFailed), msgForBindingUnmarshalFailed, name, err)
	}
	if af.validator == nil {
		return nil
	}
	return af.validator(ctx, af.ptr)
}
