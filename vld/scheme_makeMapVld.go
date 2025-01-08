package vld

import (
	"context"
	"reflect"
	"unsafe"

	"github.com/zzztttkkk/lion"
)

func makeMapVld(field *lion.Field[VldFieldMeta], meta *VldFieldMeta, gotype reflect.Type) (_PtrVldFunc, _ValVldFunc) {
	var mapfncs = []func(ctx context.Context, sv reflect.Value) *Error{}
	if meta.maxSize.Valid {
		maxs := meta.maxSize.V
		mapfncs = append(mapfncs, func(ctx context.Context, sv reflect.Value) *Error {
			lv := sv.Len()
			if lv > maxs {
				return newerr(field, meta, ErrorKindContainerSizeTooLarge).withbv(lv)
			}
			return nil
		})
	}
	if meta.minSize.Valid {
		mins := meta.minSize.V
		mapfncs = append(mapfncs, func(ctx context.Context, sv reflect.Value) *Error {
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

		var elevvfnc func(ctx context.Context, v reflect.Value) *Error
		if isperferptr {
			elevvfnc = func(ctx context.Context, vv reflect.Value) *Error {
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
					return ve.appendfield(field)
				}
				return nil
			}
		} else {
			elevvfnc = func(ctx context.Context, val reflect.Value) *Error {
				if ve := elevalvld(ctx, val.Interface()); ve != nil {
					return ve.appendfield(field)
				}
				return nil
			}
		}

		if keyvld != nil {
			if isperferptr {
				mapfncs = append(mapfncs, func(ctx context.Context, sv reflect.Value) *Error {
					iter := sv.MapRange()
					for iter.Next() {
						key := iter.Key()
						if ke := keyvld(ctx, key.Interface()); ke != nil {
							return ke.appendfield(field)
						}
						if ve := elevvfnc(ctx, iter.Value()); ve != nil {
							return ve
						}
					}
					return nil
				})
			} else {
				mapfncs = append(mapfncs, func(ctx context.Context, sv reflect.Value) *Error {
					iter := sv.MapRange()
					for iter.Next() {
						key := iter.Key()
						if ke := keyvld(ctx, key.Interface()); ke != nil {
							return ke.appendfield(field)
						}
						if ve := elevvfnc(ctx, iter.Value()); ve != nil {
							return ve
						}
					}
					return nil
				})
			}
		} else {
			mapfncs = append(mapfncs, func(ctx context.Context, sv reflect.Value) *Error {
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
		mapfncs = append(mapfncs, func(ctx context.Context, sv reflect.Value) *Error {
			iter := sv.MapRange()
			for iter.Next() {
				key := iter.Key()
				if ke := keyvld(ctx, key.Interface()); ke != nil {
					return ke.appendfield(field)
				}
			}
			return nil
		})
	}

	do := func(ctx context.Context, mapv reflect.Value) *Error {
		if !mapv.IsValid() || mapv.IsNil() {
			if meta.optional {
				return nil
			}
			return newerr(field, meta, ErrorKindNilMap)
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
	return func(ctx context.Context, uptr unsafe.Pointer) *Error {
			mapptrv := reflect.NewAt(gotype, uptr)
			return do(ctx, mapptrv.Elem())
		},
		func(ctx context.Context, val any) *Error { return do(ctx, reflect.ValueOf(val)) }
}
