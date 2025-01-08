package vld

import (
	"context"
	"fmt"
	"reflect"
	"unsafe"

	"github.com/zzztttkkk/lion"
)

func makeMapVld(field *lion.Field[VldFieldMeta], meta *VldFieldMeta, gotype reflect.Type) (_PtrVldFunc, _ValVldFunc) {
	var mapfncs = []func(ctx context.Context, sv reflect.Value) error{}
	if meta.maxSize.Valid {
		maxs := meta.maxSize.V
		mapfncs = append(mapfncs, func(ctx context.Context, sv reflect.Value) error {
			lv := sv.Len()
			if lv > maxs {
				return newerr(field, meta, ErrorKindContainerSizeTooLarge).withbv(lv)
			}
			return nil
		})
	}
	if meta.minSize.Valid {
		mins := meta.minSize.V
		mapfncs = append(mapfncs, func(ctx context.Context, sv reflect.Value) error {
			lv := sv.Len()
			if lv < mins {
				return newerr(field, meta, ErrorKindContainerSizeTooSmall).withbv(lv)
			}
			return nil
		})
	}

	eletype := gotype.Elem()
	eleisptr := eletype.Kind() == reflect.Pointer

	var eleptrvld _PtrVldFunc
	var elevalvld _ValVldFunc
	var keyvld _ValVldFunc
	if meta.ele != nil {
		eleptrvld, elevalvld = makeVldFunction(field, meta.ele, eletype)
	}
	if meta.key != nil {
		_, keyvld = makeVldFunction(field, meta.key, gotype.Key())
	}
	if elevalvld != nil || eleptrvld != nil {
		isperferptr := perferptr(eleptrvld, elevalvld, eletype)

		var elevvfnc func(ctx context.Context, v reflect.Value) error
		if isperferptr {
			elevvfnc = func(ctx context.Context, vv reflect.Value) error {
				var eleuptr unsafe.Pointer
				if eleisptr {
					valany := vv.Interface()
					valanyptr := (*anystruct)(unsafe.Pointer(&valany))
					elevptr := valanyptr.valptr
					eleuptr = unsafe.Pointer(&elevptr)
				} else {
					valany := vv.Interface()
					valanyptr := (*anystruct)(unsafe.Pointer(&valany))
					eleuptr = valanyptr.valptr
				}
				if ve := eleptrvld(ctx, eleuptr); ve != nil {
					return ve
				}
				return nil
			}
		} else {
			elevvfnc = func(ctx context.Context, val reflect.Value) error {
				return elevalvld(ctx, val.Interface())
			}
		}

		if keyvld != nil {
			if isperferptr {
				mapfncs = append(mapfncs, func(ctx context.Context, sv reflect.Value) error {
					iter := sv.MapRange()
					for iter.Next() {
						key := iter.Key()
						if ke := keyvld(ctx, key.Interface()); ke != nil {
							return ke
						}
						if ve := elevvfnc(ctx, iter.Value()); ve != nil {
							return ve
						}
					}
					return nil
				})
			} else {
				mapfncs = append(mapfncs, func(ctx context.Context, sv reflect.Value) error {
					iter := sv.MapRange()
					for iter.Next() {
						key := iter.Key()
						if ke := keyvld(ctx, key.Interface()); ke != nil {
							return ke
						}
						if ve := elevvfnc(ctx, iter.Value()); ve != nil {
							return ve
						}
					}
					return nil
				})
			}
		} else {
			mapfncs = append(mapfncs, func(ctx context.Context, sv reflect.Value) error {
				iter := sv.MapRange()
				for iter.Next() {
					if ve := elevvfnc(ctx, iter.Value()); ve != nil {
						return ve
					}
				}
				return nil
			})
		}
	} else if keyvld != nil {
		mapfncs = append(mapfncs, func(ctx context.Context, sv reflect.Value) error {
			iter := sv.MapRange()
			for iter.Next() {
				key := iter.Key()
				if ke := keyvld(ctx, key.Interface()); ke != nil {
					return ke
				}
			}
			return nil
		})
	}

	do := func(ctx context.Context, mapv reflect.Value) error {
		if !mapv.IsValid() || mapv.IsNil() {
			if meta.optional {
				return nil
			}
			return fmt.Errorf("missing required")
		}
		for _, fnc := range mapfncs {
			if err := fnc(ctx, mapv); err != nil {
				return err
			}
		}
		if meta._Func != nil {
			mapany := mapv.Interface()
			if ce := meta._Func(ctx, mapany); ce != nil {
				return newerr(field, meta, ErrorKindCustom).with(mapany, ce)
			}
		}
		return nil
	}
	return func(ctx context.Context, uptr unsafe.Pointer) error {
			mapptrv := reflect.NewAt(gotype, uptr)
			return do(ctx, mapptrv.Elem())
		},
		func(ctx context.Context, val any) error { return do(ctx, reflect.ValueOf(val)) }
}
