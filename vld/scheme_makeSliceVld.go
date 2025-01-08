package vld

import (
	"context"
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
	var slicefncs = []func(ctx context.Context, suptr unsafe.Pointer, head _SliceHeader) *Error{}
	if meta.maxSize.Valid {
		maxs := meta.maxSize.V
		slicefncs = append(slicefncs, func(ctx context.Context, suptr unsafe.Pointer, head _SliceHeader) *Error {
			if head.Len > maxs {
				return newerr(field, meta, ErrorKindContainerSizeTooLarge)
			}
			return nil
		})
	}
	if meta.minSize.Valid {
		mins := meta.minSize.V
		slicefncs = append(slicefncs, func(ctx context.Context, suptr unsafe.Pointer, head _SliceHeader) *Error {
			if head.Len < mins {
				return newerr(field, meta, ErrorKindContainerSizeTooSmall)
			}
			return nil
		})
	}

	eletype := gotype.Elem()
	eleptrfnc, eleanyfnc := makeVldFunction(field, meta.ele, gotype.Elem())
	if eleptrfnc != nil || eleanyfnc != nil {
		elesize := eletype.Size()
		if perferptr(eleptrfnc, eleanyfnc, gotype.Elem()) {
			slicefncs = append(slicefncs, func(ctx context.Context, suptr unsafe.Pointer, head _SliceHeader) *Error {
				for i := 0; i < head.Len; i++ {
					eleuptr := unsafe.Add(head.Data, i*int(elesize))
					if err := eleptrfnc(ctx, eleuptr); err != nil {
						return err.appendfield(field)
					}
				}
				return nil
			})
		} else {
			slicefncs = append(slicefncs, func(ctx context.Context, suptr unsafe.Pointer, head _SliceHeader) *Error {
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
	do := func(ctx context.Context, suptr unsafe.Pointer) *Error {
		head := *((*_SliceHeader)(suptr))
		if uintptr(head.Data) == 0 {
			if meta.optional {
				return nil
			}
			return newerr(field, meta, ErrorKindNilSlice)
		}

		for _, fnc := range slicefncs {
			err := fnc(ctx, suptr, head)
			if err != nil {
				return err
			}
		}

		if meta._Func != nil {
			sptrv := reflect.NewAt(gotype, suptr)
			sliceany := sptrv.Elem().Interface()
			if err := meta._Func(ctx, sliceany); err != nil {
				return newerr(field, meta, ErrorKindCustom).with(sliceany, err)
			}
		}
		return nil
	}
	return func(ctx context.Context, uptr unsafe.Pointer) *Error {
			return do(ctx, uptr)
		},
		func(ctx context.Context, val any) *Error {
			anyptr := (*anystruct)(unsafe.Pointer(&val))
			return do(ctx, anyptr.valptr)
		}
}
