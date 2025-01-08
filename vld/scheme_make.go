package vld

import (
	"context"
	"database/sql"
	"fmt"
	"unsafe"

	"github.com/zzztttkkk/faceless.void/internal"
	"github.com/zzztttkkk/lion"
)

type _MVType interface {
	~int64 | ~uint64
}

func makeIntOrUintVld[T lion.IntType, MV _MVType](_field *lion.Field[VldFieldMeta], meta *VldFieldMeta, minv sql.Null[MV], maxv sql.Null[MV], ranges []MV) (_PtrVldFunc, _ValVldFunc) {
	var fncs = []func(iv T) error{}
	if maxv.Valid {
		maxv := T(maxv.V)
		fncs = append(fncs, func(iv T) error {
			if iv > maxv {
				return fmt.Errorf("int gt max")
			}
			return nil
		})
	}
	if minv.Valid {
		minv := T(minv.V)
		fncs = append(fncs, func(iv T) error {
			if iv < minv {
				return fmt.Errorf("int lt min")
			}
			return nil
		})
	}
	if len(ranges) > 0 {
		var rangs []T
		for _, v := range ranges {
			rangs = append(rangs, T(v))
		}
		if len(rangs) > 16 {
			rangeset := internal.MakeSet(rangs)
			fncs = append(fncs, func(iv T) error {
				_, ok := rangeset[iv]
				if ok {
					return nil
				}
				return fmt.Errorf("int not in range")
			})
		} else {
			fncs = append(fncs, func(iv T) error {
				for _, rv := range rangs {
					if iv == rv {
						return nil
					}
				}
				return fmt.Errorf("int not in range")
			})
		}
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

func makeIntVld[T lion.SingedInt](field *lion.Field[VldFieldMeta], meta *VldFieldMeta) (_PtrVldFunc, _ValVldFunc) {
	if meta == nil {
		meta = field.Metainfo()
	}
	return makeIntOrUintVld[T](field, meta, meta.maxInt, meta.minInt, meta.intRanges)
}

func makeUintVld[T lion.UnsignedInt](field *lion.Field[VldFieldMeta], meta *VldFieldMeta) (_PtrVldFunc, _ValVldFunc) {
	if meta == nil {
		meta = field.Metainfo()
	}
	return makeIntOrUintVld[T](field, meta, meta.maxUint, meta.minUint, meta.uintRanges)
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
