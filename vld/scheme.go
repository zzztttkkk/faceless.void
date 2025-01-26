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

type _CustomAnyVldFunc func(ctx context.Context, v any) error

type VldFieldMeta struct {
	gotype reflect.Type

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
	_Func _CustomAnyVldFunc
}

func (meta *VldFieldMeta) Clone() *VldFieldMeta {
	nptr := &VldFieldMeta{}
	*nptr = *meta
	return nptr
}

type _VldItem struct {
	field  *lion.Field
	ptrfnc _PtrVldFunc
	valfnc _ValVldFunc
}

type _Scheme[T any] struct {
	typeinfo *lion.TypeInfo
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
		typeinfo: lion.TypeInfoOf[T](),
	}
	schemes[lion.Typeof[T]()] = obj
	return obj
}

func (scheme *_Scheme[T]) Scope(fnc func(ctx context.Context, mptr *T)) {
	defer scheme.Finish()
	fnc(context.WithValue(context.Background(), internal.CtxKeyForVldScheme, scheme), lion.Ptr[T]())
}

func (scheme *_Scheme[T]) Field(fptr any, meta *VldFieldMeta) *_Scheme[T] {
	field := scheme.typeinfo.FieldByPtr(fptr)
	lion.UpdateMetaFor(field, meta)
	return scheme
}

func (scheme *_Scheme[T]) Finish() {
	for _, fptr := range scheme.typeinfo.AllTagedFields("vld") {
		meta := lion.MetaOf[VldFieldMeta](fptr)
		if meta == nil {
			continue
		}
		ptrfn, valfn := makeVldFunction(fptr, meta, fptr.StructField().Type)
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

type _PtrVldFunc = func(ctx context.Context, ptr unsafe.Pointer) *Error
type _ValVldFunc = func(ctx context.Context, val any) *Error

func makeVldFunction(field *lion.Field, meta *VldFieldMeta, gotype reflect.Type) (_PtrVldFunc, _ValVldFunc) {
	EnsureSimpleContainer := func() {
		switch gotype.Elem().Kind() {
		case reflect.Slice, reflect.Map:
			{
				panic(fmt.Errorf("fv.vld: nested container is not supported, %s", gotype))
			}
		}
		if gotype.Kind() == reflect.Map {
			if !lion.Kinds.IsValue(gotype.Key().Kind()) || gotype.Key().Kind() == reflect.Struct {
				panic(fmt.Errorf("fv.vld: unsupported key kind, %s", gotype))
			}
		}
	}
	if meta == nil {
		return nil, nil
	}
	if meta.gotype != gotype {
		panic(fmt.Errorf(
			"fv.vld: bad meta gotype, field path: `%s.%s`. expected `%s`, but got `%s`.",
			field.TypeInfo().GoType, field.StructField().Name,
			gotype, meta.gotype,
		))
	}

	switch gotype.Kind() {
	case reflect.String:
		{
			return makeStringVld(field, meta)
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
					var fncs = []func(tv time.Time) *Error{}
					if meta.maxTime.Valid {
						fncs = append(fncs, func(tv time.Time) *Error {
							if tv.After(meta.maxTime.V) {
								return newerr(field, meta, ErrorKindTimeTooLate).withbv(tv)
							}
							return nil
						})
					}
					if meta.minTime.Valid {
						fncs = append(fncs, func(tv time.Time) *Error {
							if tv.Before(meta.minTime.V) {
								return newerr(field, meta, ErrorKindTimeTooEarly).withbv(tv)
							}
							return nil
						})
					}
					do := func(ctx context.Context, tv time.Time) *Error {
						for _, fn := range fncs {
							if err := fn(tv); err != nil {
								return err
							}
						}
						if meta._Func != nil {
							if ce := meta._Func(ctx, tv); ce != nil {
								return newerr(field, meta, ErrorKindCustom).with(tv, ce)
							}
						}
						return nil
					}
					return func(ctx context.Context, uptr unsafe.Pointer) *Error {
							return do(ctx, *((*time.Time)(unsafe.Add(uptr, field.Offset()))))
						},
						func(ctx context.Context, val any) *Error { return do(ctx, val.(time.Time)) }
				}
			default:
				{
					if meta.scheme.gettypeinfo().GoType != gotype {
						panic(fmt.Errorf("fv.vld: bad meta scheme, %s.%s", field.TypeInfo().GoType, field.StructField().Name))
					}
					return func(ctx context.Context, uptr unsafe.Pointer) *Error {
							return meta.scheme.dovldptr(ctx, uptr)
						}, func(ctx context.Context, val any) *Error {
							return meta.scheme.dovldval(ctx, val)
						}
				}
			}
		}
	}
	panic(fmt.Errorf("unsupported type: %s", gotype))
}

type _IScheme interface {
	gettypeinfo() *lion.TypeInfo
	updatemeta(ptr any, meta *VldFieldMeta)
	dovldptr(ctx context.Context, uptr unsafe.Pointer) *Error
	dovldval(ctx context.Context, val any) *Error
}

func (scheme *_Scheme[T]) gettypeinfo() *lion.TypeInfo {
	return scheme.typeinfo
}

func (scheme *_Scheme[T]) dovldptr(ctx context.Context, uptr unsafe.Pointer) *Error {
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

func (scheme *_Scheme[T]) dovldval(ctx context.Context, val any) *Error {
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
	lion.UpdateMetaFor(fv, meta)
}

var (
	typeofTime = lion.Typeof[time.Time]()
)

func Validate[T any](ctx context.Context, ptr *T) error {
	gotype := lion.Typeof[T]()
	scheme, ok := schemes[gotype]
	if !ok {
		return fmt.Errorf("fv.vld: can not found scheme for type: %s", gotype)
	}
	return scheme.dovldptr(ctx, unsafe.Pointer(ptr))
}
