package vld

import (
	"context"
	"fmt"
	"reflect"
	"unsafe"

	"github.com/zzztttkkk/lion"
)

const (
	PerferPtrVldSize = uintptr(1)
)

func perferptr(ptrfnc _PtrVldFunc, valfnc _ValVldFunc, eletype reflect.Type) bool {
	if ptrfnc == nil {
		return false
	}
	if valfnc == nil {
		return true
	}
	return unsafe.Sizeof(eletype) > PerferPtrVldSize
}

func makePointerVld(field *lion.Field[VldFieldMeta], meta *VldFieldMeta, gotype reflect.Type) (_PtrVldFunc, _ValVldFunc) {
	eletype := gotype.Elem()
	if eletype.Kind() == reflect.Pointer {
		panic(fmt.Errorf("fv.vld: nested pointer is not supported"))
	}
	ptrfnc, valfnc := makeVldFunction(field, meta, eletype)
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
				return fmt.Errorf("missing required")
			}
			return ptrfnc(ctx, ptrv.UnsafePointer())
		}
		return func(ctx context.Context, uptr unsafe.Pointer) error {
				return do(ctx, uptr)
			},
			func(ctx context.Context, val any) error {
				antptr := (*anystruct)(unsafe.Pointer(&val))
				return do(ctx, unsafe.Pointer(antptr.valptr))
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
