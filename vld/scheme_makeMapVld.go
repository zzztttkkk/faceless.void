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
			if sv.Len() > maxs {
				return fmt.Errorf("map too long")
			}
			return nil
		})
	}
	if meta.minSize.Valid {
		mins := meta.minSize.V
		mapfncs = append(mapfncs, func(ctx context.Context, sv reflect.Value) error {
			if sv.Len() < mins {
				return fmt.Errorf("map too short")
			}
			return nil
		})
	}
	var eleptrvld _PtrVldFunc
	var elevalvld _ValVldFunc
	var keyvld _ValVldFunc
	if meta.ele != nil {
		eleptrvld, elevalvld = makeVldFunction(field, meta.ele, gotype.Elem())
	}
	if meta.key != nil {
		_, keyvld = makeVldFunction(field, meta.key, gotype.Key())
	}
	if elevalvld != nil || eleptrvld != nil {
		isperferptr := perferptr(eleptrvld, elevalvld, gotype.Elem())
		if keyvld != nil {
			if isperferptr {
				mapfncs = append(mapfncs, func(ctx context.Context, sv reflect.Value) error {
					iter := sv.MapRange()
					for iter.Next() {
						key := iter.Key()
						if ke := keyvld(ctx, key.Interface()); ke != nil {
							return ke
						}
						valany := iter.Value().Interface()
						valanyptr := (*anystruct)(unsafe.Pointer(&valany))
						if ve := eleptrvld(ctx, valanyptr.valptr); ve != nil {
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
						val := iter.Value()
						if ve := elevalvld(ctx, val.Interface()); ve != nil {
							return ve
						}
					}
					return nil
				})
			}
		} else {
			if isperferptr {
				mapfncs = append(mapfncs, func(ctx context.Context, sv reflect.Value) error {
					iter := sv.MapRange()
					for iter.Next() {
						valany := iter.Value().Interface()
						valanyptr := (*anystruct)(unsafe.Pointer(&valany))
						if ve := eleptrvld(ctx, valanyptr.valptr); ve != nil {
							return ve
						}
					}
					return nil
				})
			} else {
				mapfncs = append(mapfncs, func(ctx context.Context, sv reflect.Value) error {
					iter := sv.MapRange()
					for iter.Next() {
						val := iter.Value()
						if ve := elevalvld(ctx, val.Interface()); ve != nil {
							return ve
						}
					}
					return nil
				})
			}
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
			return meta._Func(ctx, mapv.Interface())
		}
		return nil
	}
	return func(ctx context.Context, uptr unsafe.Pointer) error {
			mapptrv := reflect.NewAt(gotype, uptr)
			return do(ctx, mapptrv.Elem())
		},
		func(ctx context.Context, val any) error { return do(ctx, reflect.ValueOf(val)) }
}
