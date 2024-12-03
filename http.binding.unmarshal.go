package fv

import (
	"context"
	"encoding/json"
	"io"
	"reflect"

	"github.com/zzztttkkk/faceless.void/internal"
)

type _BindingAnyField struct {
	ins         *_BindingInstance
	ptr         any
	optional    bool
	unmarshaler func([]byte, any) error
	validator   func(context.Context) error
}

func (ins *_BindingInstance) Any(ptr any) *_BindingAnyField {
	af := &_BindingAnyField{
		ins: ins,
		ptr: ptr,
	}
	ins.fields = append(ins.fields, af.do)
	return af
}

func (af *_BindingAnyField) Optional() *_BindingAnyField {
	af.optional = true
	return af
}
func (af *_BindingAnyField) Unmarshal(fnc func([]byte, any) error) *_BindingAnyField {
	af.unmarshaler = fnc
	return af
}

func (af *_BindingAnyField) Validate(fnc func(ctx context.Context) error) *_BindingAnyField {
	af.validator = fnc
	return af
}

var (
	msgForBindingUnmarshalFailed = internal.NewI18nString("fv.binding: unmarshal failed, %s, %s")
)

func (af *_BindingAnyField) do(ctx context.Context) error {
	req := HttpRequest(ctx)
	defer req.Body.Close()
	payload, err := io.ReadAll(req.Body)
	if err != nil || len(payload) < 1 {
		if af.optional {
			return nil
		}
		name := af.ins.nameof(reflect.ValueOf(af.ptr).UnsafePointer())
		return internal.NewError(ctx, internal.ErrorKind(ErrorKindBindingMissingRequired), msgForBindingMissingRequired, name)
	}
	if af.unmarshaler == nil {
		err = json.Unmarshal(payload, af.ptr)
	} else {
		err = af.unmarshaler(payload, af.ptr)
	}
	if err != nil {
		name := af.ins.nameof(reflect.ValueOf(af.ptr).UnsafePointer())
		return internal.NewError(ctx, internal.ErrorKind(ErrorKindBindingUnmarshalFailed), msgForBindingUnmarshalFailed, name, err)
	}
	if af.validator == nil {
		return nil
	}
	return af.validator(ctx)
}
