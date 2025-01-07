package vld

import (
	"context"
	"fmt"
	"reflect"
	"unsafe"

	"github.com/zzztttkkk/lion"
)

func makeSliceVld(field *lion.Field[VldFieldMeta], meta *VldFieldMeta, gotype reflect.Type) (_PtrVldFunc, _ValVldFunc) {
	var slicefncs = []func(ctx context.Context, sv reflect.Value) error{}
	if gotype.Kind() == reflect.Slice {
		if meta.maxSize.Valid {
			maxs := meta.maxSize.V
			slicefncs = append(slicefncs, func(ctx context.Context, sv reflect.Value) error {
				if sv.Len() > maxs {
					return fmt.Errorf("slice too long")
				}
				return nil
			})
		}
		if meta.minSize.Valid {
			mins := meta.minSize.V
			slicefncs = append(slicefncs, func(ctx context.Context, sv reflect.Value) error {
				if sv.Len() < mins {
					return fmt.Errorf("slice too short")
				}
				return nil
			})
		}
	}
	eleptrfnc, eleanyfnc := makeVldFunction(field, meta.ele, gotype.Elem())
	if eleptrfnc != nil || eleanyfnc != nil {
		if perferptr(eleptrfnc, eleanyfnc, gotype.Elem()) {
			slicefncs = append(slicefncs, func(ctx context.Context, sv reflect.Value) error {
				slen := sv.Len()
				for i := 0; i < slen; i++ {
					eleav := sv.Index(i).Interface()
					eleavptr := (*anystruct)(unsafe.Pointer(&eleav))
					fmt.Println(eleav, eleavptr.valptr)
					if ee := eleptrfnc(ctx, unsafe.Pointer(eleavptr.valptr)); ee != nil {
						return ee
					}
				}
				return nil
			})
		} else {
			slicefncs = append(slicefncs, func(ctx context.Context, sv reflect.Value) error {
				slen := sv.Len()
				for i := 0; i < slen; i++ {
					elev := sv.Index(i)
					if ee := eleanyfnc(ctx, elev.Interface()); ee != nil {
						return ee
					}
				}
				return nil
			})
		}
	}
	do := func(ctx context.Context, sv reflect.Value) error {
		if !sv.IsValid() || sv.IsNil() {
			if meta.optional {
				return nil
			}
			return fmt.Errorf("missing required")
		}
		for _, fn := range slicefncs {
			if err := fn(ctx, sv); err != nil {
				return err
			}
		}
		if meta._Func != nil {
			return meta._Func(ctx, sv.Interface())
		}
		return nil
	}
	return func(ctx context.Context, uptr unsafe.Pointer) error {
			sliceptrv := reflect.NewAt(field.StructField().Type, uptr)
			return do(ctx, sliceptrv.Elem())
		},
		func(ctx context.Context, val any) error { return do(ctx, reflect.ValueOf(val)) }
}
