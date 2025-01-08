package vld

import (
	"context"
	"fmt"
	"reflect"
	"unsafe"

	"github.com/zzztttkkk/lion"
)

type _SliceHeader struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}

func makeSliceVld(field *lion.Field[VldFieldMeta], meta *VldFieldMeta, gotype reflect.Type) (_PtrVldFunc, _ValVldFunc) {
	var slicefncs = []func(ctx context.Context, suptr unsafe.Pointer, head _SliceHeader) error{}
	if meta.maxSize.Valid {
		maxs := meta.maxSize.V
		slicefncs = append(slicefncs, func(ctx context.Context, suptr unsafe.Pointer, head _SliceHeader) error {
			if head.Len > maxs {
				return fmt.Errorf("slice too long")
			}
			return nil
		})
	}
	if meta.minSize.Valid {
		mins := meta.minSize.V
		slicefncs = append(slicefncs, func(ctx context.Context, suptr unsafe.Pointer, head _SliceHeader) error {
			if head.Len < mins {
				return fmt.Errorf("slice too short")
			}
			return nil
		})
	}

	eletype := gotype.Elem()
	eleptrfnc, eleanyfnc := makeVldFunction(field, meta.ele, gotype.Elem())
	if eleptrfnc != nil || eleanyfnc != nil {
		elesize := eletype.Size()
		if perferptr(eleptrfnc, eleanyfnc, gotype.Elem()) {
			slicefncs = append(slicefncs, func(ctx context.Context, suptr unsafe.Pointer, head _SliceHeader) error {
				for i := 0; i < head.Len; i++ {
					eleuptr := unsafe.Add(head.Data, i*int(elesize))
					if err := eleptrfnc(ctx, eleuptr); err != nil {
						return err
					}
				}
				return nil
			})
		} else {
			slicefncs = append(slicefncs, func(ctx context.Context, suptr unsafe.Pointer, head _SliceHeader) error {
				for i := 0; i < head.Len; i++ {
					eleuptr := unsafe.Add(head.Data, i*int(elesize))
					eleptrv := reflect.NewAt(eletype, eleuptr)
					if err := eleanyfnc(ctx, eleptrv.Elem().Interface()); err != nil {
						return err
					}
				}
				return nil
			})
		}
	}
	do := func(ctx context.Context, suptr unsafe.Pointer) error {
		head := *((*_SliceHeader)(suptr))
		if uintptr(head.Data) == 0 {
			if meta.optional {
				return nil
			}
			return fmt.Errorf("nil slice")
		}

		for _, fnc := range slicefncs {
			err := fnc(ctx, suptr, head)
			if err != nil {
				return err
			}
		}

		if meta._Func != nil {
			sptrv := reflect.NewAt(gotype, suptr)
			if err := meta._Func(ctx, sptrv.Elem().Interface()); err != nil {
				return err
			}
		}
		return nil
	}
	return func(ctx context.Context, uptr unsafe.Pointer) error {
			return do(ctx, uptr)
		},
		func(ctx context.Context, val any) error {
			anyptr := (*anystruct)(unsafe.Pointer(&val))
			return do(ctx, anyptr.valptr)
		}
}
