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
		do := func(ctx context.Context, uptr unsafe.Pointer) error {
			ptrptrv := reflect.NewAt(gotype, uptr)
			ptrv := ptrptrv.Elem()
			if !ptrv.IsValid() || ptrv.IsNil() {
				if meta.optional {
					return nil
				}
				return &Error{Field: field, Kind: ErrorKindNilPointer}
			}
			return ptrfnc(ctx, ptrv.UnsafePointer())
		}
		return func(ctx context.Context, uptr unsafe.Pointer) error {
				return do(ctx, uptr)
			},
			func(ctx context.Context, val any) error {
				antptr := (*anystruct)(unsafe.Pointer(&val))
				return do(ctx, antptr.valptr)
			}
	}

	do := func(ctx context.Context, pv reflect.Value) error {
		if !pv.IsValid() || pv.IsNil() {
			if meta.optional {
				return nil
			}
			return fmt.Errorf("missing required")
		}
		return valfnc(ctx, pv.Elem().Interface())
	}
	return func(ctx context.Context, uptr unsafe.Pointer) error {
			ptrval := reflect.NewAt(field.StructField().Type, uptr)
			return do(ctx, ptrval.Elem())
		},
		func(ctx context.Context, val any) error { return do(ctx, reflect.ValueOf(val)) }
}
