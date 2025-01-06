package vld

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"regexp"
	"time"
	"unsafe"

	"github.com/zzztttkkk/faceless.void/internal"
	"github.com/zzztttkkk/lion"
)

type VldFieldMeta struct {
	optional bool

	// string
	maxLength    sql.Null[int]
	minLength    sql.Null[int]
	regexp       *regexp.Regexp
	stringRanges []string

	// int
	maxInt    sql.Null[int64]
	minInt    sql.Null[int64]
	intRanges []int64

	// uint
	maxUint    sql.Null[uint64]
	minUint    sql.Null[uint64]
	uintRanges []uint64

	// time
	maxTime sql.Null[time.Time]
	minTime sql.Null[time.Time]

	// slice/map
	maxSize sql.Null[int]
	minSize sql.Null[int]
	key     *VldFieldMeta
	ele     *VldFieldMeta

	// scheme
	scheme _IScheme

	// custom
	_Func func(ctx context.Context, v any) error
}

func init() {
	lion.RegisterOf[VldFieldMeta]().TagNames("vld").Unexposed()
}

type _Scheme[T any] struct {
	typeinfo *lion.TypeInfo[VldFieldMeta]
	vlds     []_PtrVldFunc
}

var (
	schemes         = map[reflect.Type]_IScheme{}
	copyToHeapFuncs = map[reflect.Type]func(v any) reflect.Value{}
)

func SchemeOf[T any]() *_Scheme[T] {
	gotype := lion.Typeof[T]()
	v, ok := schemes[gotype]
	if ok {
		return (any(v)).(*_Scheme[T])
	}
	copyToHeapFuncs[gotype] = func(v any) reflect.Value {
		ptr := new(T)
		*ptr = (v.(T))
		return reflect.ValueOf(ptr)
	}

	obj := &_Scheme[T]{
		typeinfo: lion.TypeInfoOf[T, VldFieldMeta](),
	}
	schemes[lion.Typeof[T]()] = obj
	return obj
}

func (scheme *_Scheme[T]) Scope(fnc func(ctx context.Context, mptr *T)) {
	defer scheme.Finish()
	fnc(context.WithValue(context.Background(), internal.CtxKeyForVldScheme, scheme), Ptr[T]())
}

func (scheme *_Scheme[T]) Field(fptr any, meta *VldFieldMeta) *_Scheme[T] {
	field := scheme.typeinfo.FieldByPtr(fptr)
	field.UpdateMetainfo(meta)
	return scheme
}

func (scheme *_Scheme[T]) Finish() {
	for idx := range scheme.typeinfo.Fields {
		fptr := &scheme.typeinfo.Fields[idx]
		if fptr.Metainfo() == nil {
			continue
		}
		fn, _ := makeVldFunction(fptr, fptr.Metainfo(), fptr.StructField().Type)
		if fn == nil {
			continue
		}
		scheme.vlds = append(scheme.vlds, fn)
	}
}

type _PtrVldFunc = func(ctx context.Context, ptr unsafe.Pointer) error
type _ValVldFunc = func(ctx context.Context, val any) error

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

func makeVldFunction(field *lion.Field[VldFieldMeta], meta *VldFieldMeta, gotype reflect.Type) (_PtrVldFunc, _ValVldFunc) {
	EnsureSimpleContainer := func() {
		switch gotype.Elem().Kind() {
		case reflect.Array, reflect.Slice, reflect.Map:
			{
				panic(fmt.Errorf("fv.vld: nested container is not supported, %s", gotype))
			}
		}
	}
	if meta == nil {
		return nil, nil
	}

	getter := field.Getter()
	switch gotype.Kind() {
	case reflect.String:
		{
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
						return fmt.Errorf("string not match")
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
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		{
			switch gotype.Bits() {
			case 8:
				return makeIntVld[int8](field, meta)
			case 16:
				return makeIntVld[int16](field, meta)
			case 32:
				return makeIntVld[int32](field, meta)
			case 64:
				return makeIntVld[int64](field, meta)
			}
			break
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		{
			switch gotype.Bits() {
			case 8:
				return makeUintVld[uint8](field, meta)
			case 16:
				return makeUintVld[uint16](field, meta)
			case 32:
				return makeUintVld[uint32](field, meta)
			case 64:
				return makeUintVld[uint64](field, meta)
			}
			break
		}
	case reflect.Slice, reflect.Array:
		{
			EnsureSimpleContainer()
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
			_, eleanyfnc := makeVldFunction(field, meta.ele, gotype.Elem())
			if eleanyfnc != nil {
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
			return func(ctx context.Context, uptr unsafe.Pointer) error { return do(ctx, reflect.ValueOf(getter(uptr))) },
				func(ctx context.Context, val any) error { return do(ctx, reflect.ValueOf(val)) }
		}
	case reflect.Map:
		{
			EnsureSimpleContainer()

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
			var elevld _ValVldFunc
			var keyvld _ValVldFunc
			if meta.ele != nil {
				_, elevld = makeVldFunction(field, meta.ele, gotype.Elem())
			}
			if meta.key != nil {
				_, keyvld = makeVldFunction(field, meta.key, gotype.Key())
			}
			if elevld != nil {
				if keyvld != nil {
					mapfncs = append(mapfncs, func(ctx context.Context, sv reflect.Value) error {
						iter := sv.MapRange()
						for iter.Next() {
							key := iter.Key()
							if ke := keyvld(ctx, key.Interface()); ke != nil {
								return ke
							}
							val := iter.Value()
							if ve := elevld(ctx, val.Interface()); ve != nil {
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
							if ve := elevld(ctx, val.Interface()); ve != nil {
								return ve
							}
						}
						return nil
					})
				}
			} else {
				if keyvld != nil {
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
			return func(ctx context.Context, uptr unsafe.Pointer) error { return do(ctx, reflect.ValueOf(getter(uptr))) },
				func(ctx context.Context, val any) error { return do(ctx, reflect.ValueOf(val)) }
		}
	case reflect.Pointer:
		{
			if gotype.Elem().Kind() == reflect.Pointer {
				panic(fmt.Errorf("fv.vld: nested pointer is not supported"))
			}
			_, valfnc := makeVldFunction(field, meta, gotype.Elem())
			if valfnc == nil {
				return nil, nil
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
			return func(ctx context.Context, uptr unsafe.Pointer) error { return do(ctx, reflect.ValueOf(getter(uptr))) },
				func(ctx context.Context, val any) error { return do(ctx, reflect.ValueOf(val)) }
		}
	default:
		{
			switch gotype {
			case typeofTime:
				{
					var fncs = []func(tv time.Time) error{}
					if meta.maxTime.Valid {
						fncs = append(fncs, func(tv time.Time) error {
							if tv.After(meta.maxTime.V) {
								return fmt.Errorf("time too late")
							}
							return nil
						})
					}
					if meta.minTime.Valid {
						fncs = append(fncs, func(tv time.Time) error {
							if tv.Before(meta.minTime.V) {
								return fmt.Errorf("time too early")
							}
							return nil
						})
					}
					do := func(ctx context.Context, tv time.Time) error {
						for _, fn := range fncs {
							if err := fn(tv); err != nil {
								return err
							}
						}
						if meta._Func != nil {
							return meta._Func(ctx, tv)
						}
						return nil
					}
					return func(ctx context.Context, uptr unsafe.Pointer) error {
							return do(ctx, *((*time.Time)(unsafe.Add(uptr, field.Offset()))))
						},
						func(ctx context.Context, val any) error { return do(ctx, val.(time.Time)) }
				}
			default:
				{
					if meta.scheme.gettypeinfo().GoType != gotype {
						panic(fmt.Errorf("fv.vld: bad meta scheme, %s.%s", field.Typeinfo().GoType, field.StructField().Name))
					}
					return func(ctx context.Context, uptr unsafe.Pointer) error {
							return meta.scheme.dovld(ctx, unsafe.Add(uptr, field.Offset()))
						}, func(ctx context.Context, val any) error {
							ptrv := copyToHeapFuncs[gotype](val)
							return meta.scheme.dovld(ctx, ptrv.UnsafePointer())
						}
				}
			}
		}
	}
	panic(fmt.Errorf("unsupported type: %s", gotype))
}

type _IScheme interface {
	gettypeinfo() *lion.TypeInfo[VldFieldMeta]
	updatemeta(ptr any, meta *VldFieldMeta)
	dovld(ctx context.Context, uptr unsafe.Pointer) error
}

func (scheme *_Scheme[T]) gettypeinfo() *lion.TypeInfo[VldFieldMeta] {
	return scheme.typeinfo
}

func (scheme *_Scheme[T]) dovld(ctx context.Context, uptr unsafe.Pointer) error {
	for _, vld := range scheme.vlds {
		if err := vld(ctx, uptr); err != nil {
			return err
		}
	}
	return nil
}

func (scheme *_Scheme[T]) updatemeta(fptr any, meta *VldFieldMeta) {
	fv := scheme.typeinfo.FieldByPtr(fptr)
	if fv == nil {
		panic(fmt.Errorf("fv.vld: bad field point. %s, pointer: %p", scheme.typeinfo.GoType, fptr))
	}
	fv.UpdateMetainfo(meta)
}

var (
	typeofTime = lion.Typeof[time.Time]()
)

func Vld[T any](ctx context.Context, ptr *T) error {
	gotype := lion.Typeof[T]()
	scheme, ok := schemes[gotype]
	if !ok {
		return fmt.Errorf("fv.vld: can not found scheme for type: %s", gotype)
	}
	return scheme.dovld(ctx, unsafe.Pointer(ptr))
}
