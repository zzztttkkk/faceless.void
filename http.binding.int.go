package fv

import (
	"context"
	"fmt"
	"math"
	"reflect"
	"unsafe"

	"github.com/zzztttkkk/faceless.void/internal"
)

type _IntFieldBindding[T internal.IntType] struct {
	ins   *_BindingInstance
	name  string
	alias []string
	where BindingSrcKind
	ptr   *T
	vld   func(context.Context, T) error

	base           int
	optional       bool
	defaultvalue   T
	defaultvalueok bool
}

func (ifb *_IntFieldBindding[T]) Default(v T) *_IntFieldBindding[T] {
	ifb.defaultvalue = v
	ifb.defaultvalueok = true
	return ifb
}

func (ifb *_IntFieldBindding[T]) Optional() *_IntFieldBindding[T] {
	ifb.optional = true
	return ifb
}

func (ifb *_IntFieldBindding[T]) Base(base int) *_IntFieldBindding[T] {
	ifb.base = base
	return ifb
}

func (ifb *_IntFieldBindding[T]) Name(name string, alias ...string) *_IntFieldBindding[T] {
	ifb.name = name
	ifb.alias = alias
	return ifb
}

func (ifb *_IntFieldBindding[T]) From(src BindingSrcKind) *_IntFieldBindding[T] {
	ifb.where = src
	return ifb
}

func (ifb *_IntFieldBindding[T]) Validate(vld func(context.Context, T) error) *_IntFieldBindding[T] {
	ifb.vld = vld
	return ifb
}

func (ifb *_IntFieldBindding[T]) do(ctx context.Context) error {
	if ifb.name == "" {
		ifb.name = ifb.ins.nameof(unsafe.Pointer(ifb.ptr))
	}

	itype := reflect.TypeOf(T(0))
	if itype.Name()[0] == 'i' {
		var min, max int64
		switch itype.Bits() {
		case 8:
			{
				min = math.MinInt8
				max = math.MaxInt8
				break
			}
		case 16:
			{
				min = math.MinInt16
				max = math.MaxInt16
				break
			}
		case 32:
			{
				min = math.MinInt32
				max = math.MaxInt32
				break
			}
		case 64:
			{
				min = math.MinInt64
				max = math.MaxInt64
				break
			}
		}

		iv, ok := _BindingGetter(ctx).Int(ifb.where, ifb.base, ifb.name, ifb.alias...)
		if !ok {
			if !ifb.defaultvalueok {
				if !ifb.optional {
					return fmt.Errorf("missing required field, `%s`", ifb.name)
				}
				return nil
			}
			iv = int64(ifb.defaultvalue)
		}
		if iv < min || iv > max {
			return fmt.Errorf("int overflow")
		}
		*ifb.ptr = T(iv)
	} else {
		var max uint64
		switch itype.Bits() {
		case 8:
			{
				max = math.MaxUint8
				break
			}
		case 16:
			{
				max = math.MaxUint16
				break
			}
		case 32:
			{
				max = math.MaxUint32
				break
			}
		case 64:
			{
				max = math.MaxUint64
				break
			}
		}

		uv, ok := _BindingGetter(ctx).Uint(ifb.where, ifb.name, ifb.alias...)
		if !ok {
			if !ifb.defaultvalueok {
				if !ifb.optional {
					return fmt.Errorf("miss")
				}
				return nil
			}
			uv = uint64(ifb.defaultvalue)
		}
		if uv > max {
			return fmt.Errorf("uint overflow")
		}
		*ifb.ptr = T(uv)
	}
	return nil
}

func (ins *_BindingInstance) Int(ptr *int) *_IntFieldBindding[int] {
	ifb := &_IntFieldBindding[int]{
		ins: ins,
		ptr: ptr,
	}
	ins.fields = append(ins.fields, ifb.do)
	return ifb
}

func (ins *_BindingInstance) Int8(ptr *int8) *_IntFieldBindding[int8] {
	ifb := &_IntFieldBindding[int8]{
		ins: ins,
		ptr: ptr,
	}
	ins.fields = append(ins.fields, ifb.do)
	return ifb
}

func (ins *_BindingInstance) Int16(ptr *int16) *_IntFieldBindding[int16] {
	ifb := &_IntFieldBindding[int16]{
		ins: ins,
		ptr: ptr,
	}
	ins.fields = append(ins.fields, ifb.do)
	return ifb
}

func (ins *_BindingInstance) Int32(ptr *int32) *_IntFieldBindding[int32] {
	ifb := &_IntFieldBindding[int32]{
		ins: ins,
		ptr: ptr,
	}
	ins.fields = append(ins.fields, ifb.do)
	return ifb
}

func (ins *_BindingInstance) Int64(ptr *int64) *_IntFieldBindding[int64] {
	ifb := &_IntFieldBindding[int64]{
		ins: ins,
		ptr: ptr,
	}
	ins.fields = append(ins.fields, ifb.do)
	return ifb
}

func (ins *_BindingInstance) Uint(ptr *uint) *_IntFieldBindding[uint] {
	ifb := &_IntFieldBindding[uint]{
		ins: ins,
		ptr: ptr,
	}
	ins.fields = append(ins.fields, ifb.do)
	return ifb
}

func (ins *_BindingInstance) Uint8(ptr *uint8) *_IntFieldBindding[uint8] {
	ifb := &_IntFieldBindding[uint8]{
		ins: ins,
		ptr: ptr,
	}
	ins.fields = append(ins.fields, ifb.do)
	return ifb
}

func (ins *_BindingInstance) Uint16(ptr *uint16) *_IntFieldBindding[uint16] {
	ifb := &_IntFieldBindding[uint16]{
		ins: ins,
		ptr: ptr,
	}
	ins.fields = append(ins.fields, ifb.do)
	return ifb
}

func (ins *_BindingInstance) Uint32(ptr *uint32) *_IntFieldBindding[uint32] {
	ifb := &_IntFieldBindding[uint32]{
		ins: ins,
		ptr: ptr,
	}
	ins.fields = append(ins.fields, ifb.do)
	return ifb
}

func (ins *_BindingInstance) Uint64(ptr *uint64) *_IntFieldBindding[uint64] {
	ifb := &_IntFieldBindding[uint64]{
		ins: ins,
		ptr: ptr,
	}
	ins.fields = append(ins.fields, ifb.do)
	return ifb
}
