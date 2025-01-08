package vld

import (
	"context"
	"fmt"
	"reflect"
	"unsafe"

	"github.com/zzztttkkk/lion"
)

func perferptr(ptrfnc _PtrVldFunc, valfnc _ValVldFunc, eletype reflect.Type) bool {
	if ptrfnc == nil {
		return false
	}
	if valfnc == nil {
		return true
	}
	return eletype.Size() > PerferPtrVldSizeThreshold
}

func makePointerVld(field *lion.Field[VldFieldMeta], meta *VldFieldMeta, gotype reflect.Type) (_PtrVldFunc, _ValVldFunc) {
	eletype := gotype.Elem()
	if !lion.Kinds.IsValue(eletype.Kind()) {
		panic(fmt.Errorf("fv.vld: pointer ele is not value, %s", gotype))
	}
	ptrfnc, valfnc := makeVldFunction(field, meta.ele, eletype)
	if ptrfnc == nil && valfnc == nil {
		return nil, nil
	}
	if perferptr(ptrfnc, valfnc, eletype) {
		do := func(ctx context.Context, uptr unsafe.Pointer) *Error {
			ptrptrv := reflect.NewAt(gotype, uptr)
			ptrv := ptrptrv.Elem()
			if !ptrv.IsValid() || ptrv.IsNil() {
				if meta.optional {
					return nil
				}
				return newerr(field, meta, ErrorKindNilPointer)
			}
			if eev := ptrfnc(ctx, ptrv.UnsafePointer()); eev != nil {
				return eev.appendfield(field)
			}
			return nil
		}
		return func(ctx context.Context, uptr unsafe.Pointer) *Error {
				return do(ctx, uptr)
			},
			func(ctx context.Context, val any) *Error {
				antptr := (*anystruct)(unsafe.Pointer(&val))
				return do(ctx, antptr.valptr)
			}
	}

	do := func(ctx context.Context, pv reflect.Value) *Error {
		if !pv.IsValid() || pv.IsNil() {
			if meta.optional {
				return nil
			}
			return newerr(field, meta, ErrorKindNilPointer)
		}
		if eev := valfnc(ctx, pv.Elem().Interface()); eev != nil {
			return eev.appendfield(field)
		}
		return nil

	}
	return func(ctx context.Context, uptr unsafe.Pointer) *Error {
			ptrval := reflect.NewAt(field.StructField().Type, uptr)
			return do(ctx, ptrval.Elem())
		},
		func(ctx context.Context, val any) *Error { return do(ctx, reflect.ValueOf(val)) }
}
