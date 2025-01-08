package vld

import (
	"context"
	"fmt"
	"unsafe"

	"github.com/zzztttkkk/lion"
)

func makeIntVld[T lion.SingedInt](field *lion.Field[VldFieldMeta], meta *VldFieldMeta) (_PtrVldFunc, _ValVldFunc) {
	if meta == nil {
		meta = field.Metainfo()
	}
	var fncs = []func(iv T) error{}
	if meta.maxInt.Valid {
		maxv := T(meta.maxInt.V)
		fncs = append(fncs, func(iv T) error {
			if iv > maxv {
				return fmt.Errorf("int gt max")
			}
			return nil
		})
	}
	if meta.minInt.Valid {
		minv := T(meta.minInt.V)
		fncs = append(fncs, func(iv T) error {
			if iv < minv {
				return fmt.Errorf("int lt min")
			}
			return nil
		})
	}

	do := func(ctx context.Context, iv T) error {
		for _, fn := range fncs {
			if err := fn(iv); err != nil {
				return err
			}
		}
		if meta._Func != nil {
			return meta._Func(ctx, iv)
		}
		return nil
	}
	return func(ctx context.Context, uptr unsafe.Pointer) error { return do(ctx, *((*T)(uptr))) },
		func(ctx context.Context, val any) error { return do(ctx, val.(T)) }
}

func makeUintVld[T lion.UnsignedInt](field *lion.Field[VldFieldMeta], meta *VldFieldMeta) (_PtrVldFunc, _ValVldFunc) {
	if meta == nil {
		meta = field.Metainfo()
	}
	var fncs = []func(iv T) error{}
	if meta.maxUint.Valid {
		maxv := T(meta.maxUint.V)
		fncs = append(fncs, func(iv T) error {
			if iv > maxv {
				return fmt.Errorf("uint gt max")
			}
			return nil
		})
	}
	if meta.minUint.Valid {
		minv := T(meta.minUint.V)
		fncs = append(fncs, func(iv T) error {
			if iv < minv {
				return fmt.Errorf("uint lt min")
			}
			return nil
		})
	}
	do := func(ctx context.Context, iv T) error {
		for _, fn := range fncs {
			if err := fn(iv); err != nil {
				return err
			}
		}
		if meta._Func != nil {
			return meta._Func(ctx, iv)
		}
		return nil
	}
	return func(ctx context.Context, uptr unsafe.Pointer) error {
			return do(ctx, *((*T)(uptr)))
		}, func(ctx context.Context, val any) error {
			return do(ctx, val.(T))
		}
}

func makeStringVld(meta *VldFieldMeta) (_PtrVldFunc, _ValVldFunc) {
	var fncs = []func(v string) error{}

	if meta.maxLength.Valid {
		maxl := meta.maxLength.V
		fncs = append(fncs, func(v string) error {
			if len(v) > maxl {
				return fmt.Errorf("string too long")
			}
			return nil
		})
	}
	if meta.minLength.Valid {
		minl := meta.minLength.V
		fncs = append(fncs, func(v string) error {
			if len(v) < minl {
				return fmt.Errorf("string too short")
			}
			return nil
		})
	}
	if meta.regexp != nil {
		fncs = append(fncs, func(v string) error {
			if !meta.regexp.MatchString(v) {
				return fmt.Errorf("string not match, %s", v)
			}
			return nil
		})
	}
	if len(meta.stringRanges) > 0 {
		if len(meta.stringRanges) > 15 {
			var rangemap = map[string]struct{}{}
			for _, rv := range meta.stringRanges {
				rangemap[rv] = struct{}{}
			}
			fncs = append(fncs, func(v string) error {
				_, ok := rangemap[v]
				if ok {
					return nil
				}
				return fmt.Errorf("string not in range")
			})
		} else {
			fncs = append(fncs, func(v string) error {
				for _, rv := range meta.stringRanges {
					if v == rv {
						return nil
					}
				}
				return fmt.Errorf("string not in range")
			})
		}
	}

	do := func(ctx context.Context, fv string) error {
		for _, fn := range fncs {
			if err := fn(fv); err != nil {
				return err
			}
		}
		if meta._Func != nil {
			return meta._Func(ctx, fv)
		}
		return nil
	}
	return func(ctx context.Context, uptr unsafe.Pointer) error {
			fv := *((*string)(uptr))
			return do(ctx, fv)
		}, func(ctx context.Context, val any) error {
			return do(ctx, val.(string))
		}
}
