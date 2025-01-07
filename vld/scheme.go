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

type _VldItem struct {
	field  *lion.Field[VldFieldMeta]
	ptrfnc _PtrVldFunc
	valfnc _ValVldFunc
}

type _Scheme[T any] struct {
	typeinfo *lion.TypeInfo[VldFieldMeta]
	vlds     []_VldItem
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
		ptrfn, valfn := makeVldFunction(fptr, fptr.Metainfo(), fptr.StructField().Type)
		if ptrfn == nil {
			continue
		}
		scheme.vlds = append(scheme.vlds, _VldItem{
			ptrfnc: ptrfn,
			valfnc: valfn,
			field:  fptr,
		})
	}
}

type _PtrVldFunc = func(ctx context.Context, ptr unsafe.Pointer) error
type _ValVldFunc = func(ctx context.Context, val any) error

func makeVldFunction(field *lion.Field[VldFieldMeta], meta *VldFieldMeta, gotype reflect.Type) (_PtrVldFunc, _ValVldFunc) {
	EnsureSimpleContainer := func() {
		switch gotype.Elem().Kind() {
		case reflect.Array, reflect.Slice, reflect.Map:
			{
				panic(fmt.Errorf("fv.vld: nested container is not supported, %s", gotype))
			}
		case reflect.Pointer:
			{
				// see fcuntion `makeSliceVld`
				panic(fmt.Errorf("fv.vld: change ele kind to value, %s", gotype))
			}
		}

		if gotype.Kind() == reflect.Map {
			switch gotype.Key().Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
				reflect.String, reflect.Struct, reflect.Bool:
				{
					break
				}
			default:
				{
					panic(fmt.Errorf("fv.vld: unsupported key kind, %s", gotype))
				}
			}
		}
	}
	if meta == nil {
		return nil, nil
	}

	switch gotype.Kind() {
	case reflect.String:
		{
			return makeStringVld(meta)
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
	case reflect.Slice:
		{
			EnsureSimpleContainer()
			return makeSliceVld(field, meta, gotype)
		}
	case reflect.Map:
		{
			EnsureSimpleContainer()
			return makeMapVld(field, meta, gotype)
		}
	case reflect.Pointer:
		{
			return makePointerVld(field, meta, gotype)
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
							return meta.scheme.dovldptr(ctx, uptr)
						}, func(ctx context.Context, val any) error {
							return meta.scheme.dovldval(ctx, val)
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
	dovldptr(ctx context.Context, uptr unsafe.Pointer) error
	dovldval(ctx context.Context, val any) error
}

func (scheme *_Scheme[T]) gettypeinfo() *lion.TypeInfo[VldFieldMeta] {
	return scheme.typeinfo
}

func (scheme *_Scheme[T]) dovldptr(ctx context.Context, uptr unsafe.Pointer) error {
	for _, vld := range scheme.vlds {
		if err := vld.ptrfnc(ctx, unsafe.Add(uptr, vld.field.Offset())); err != nil {
			return err
		}
	}
	return nil
}

type anystruct struct {
	_typeptr unsafe.Pointer
	valptr   unsafe.Pointer
}

func (scheme *_Scheme[T]) dovldval(ctx context.Context, val any) error {
	vv := reflect.ValueOf(val)
	if vv.Kind() == reflect.Pointer {
		return scheme.dovldptr(ctx, vv.UnsafePointer())
	}
	var ptr = (*anystruct)(unsafe.Pointer((&val)))
	return scheme.dovldptr(ctx, ptr.valptr)
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
	return scheme.dovldptr(ctx, unsafe.Pointer(ptr))
}
