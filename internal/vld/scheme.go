package internalvld

import (
	"database/sql"
	"fmt"
	"reflect"
	"regexp"
	"time"
	"unsafe"

	"github.com/zzztttkkk/lion"
)

type VldFieldMeta struct {
	Optional bool

	// string
	MaxLength    sql.Null[int]
	MinLength    sql.Null[int]
	Regexp       *regexp.Regexp
	StringRanges []string

	// int
	MaxInt    sql.Null[int64]
	MinInt    sql.Null[int64]
	IntRanges []int64

	// uint
	MaxUint    sql.Null[uint64]
	UintRanges []uint64

	// time
	MaxTime sql.Null[time.Time]
	MinTime sql.Null[time.Time]

	// slice/map
	MaxSize sql.Null[int]
	MinSize sql.Null[int]
	MapKey  *VldFieldMeta

	// custom
	Func func(v any) error
}

func init() {
	lion.RegisterOf[VldFieldMeta]().TagNames("vld", "bnd", "db", "json").Unexposed()
}

type _Scheme struct {
	typeinfo *lion.TypeInfo[VldFieldMeta]
	vlds     []func(ptr unsafe.Pointer) error
}

var (
	schemes = map[reflect.Type]*_Scheme{}
)

func SchemeOf[T any]() *_Scheme {
	v, ok := schemes[lion.Typeof[T]()]
	if ok {
		return v
	}
	obj := &_Scheme{
		typeinfo: lion.TypeInfoOf[T, VldFieldMeta](),
	}
	schemes[lion.Typeof[T]()] = obj
	return obj
}

func (scheme *_Scheme) Field(fptr any, meta *VldFieldMeta) *_Scheme {
	field := scheme.typeinfo.FieldByPtr(fptr)
	field.UpdateMetainfo(meta)
	return scheme
}

func (scheme *_Scheme) Finish() {
	for idx := range scheme.typeinfo.Fields {
		fptr := &scheme.typeinfo.Fields[idx]
		if fptr.Metainfo() == nil {
			continue
		}
		fn, _ := makeVldFunction(fptr, nil, fptr.StructField().Type)
		if fn == nil {
			continue
		}
		scheme.vlds = append(scheme.vlds, fn)
	}
}

type inttypes interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func makeIntVld[T inttypes](field *lion.Field[VldFieldMeta], meta *VldFieldMeta) (func(uptr unsafe.Pointer) error, func(val any) error) {
	if meta == nil {
		meta = field.Metainfo()
	}
	var fncs = []func(iv T) error{}
	if meta.MaxInt.Valid {
		maxv := T(meta.MaxInt.V)
		fncs = append(fncs, func(iv T) error {
			if iv > maxv {
				return fmt.Errorf("int gt max")
			}
			return nil
		})
	}
	if meta.MinInt.Valid {
		minv := T(meta.MinInt.V)
		fncs = append(fncs, func(iv T) error {
			if iv < minv {
				return fmt.Errorf("int lt min")
			}
			return nil
		})
	}
	if meta.Func != nil {
		fncs = append(fncs, func(iv T) error { return meta.Func(iv) })
	}
	if len(fncs) < 1 {
		return nil, nil
	}

	do := func(iv T) error {
		for _, fn := range fncs {
			if err := fn(iv); err != nil {
				return err
			}
		}
		return nil
	}
	return func(uptr unsafe.Pointer) error { return do(*((*T)(unsafe.Add(uptr, field.Offset())))) },
		func(val any) error { return do(val.(T)) }
}

func makeUintVld[T inttypes](field *lion.Field[VldFieldMeta], meta *VldFieldMeta) (func(uptr unsafe.Pointer) error, func(val any) error) {
	if meta == nil {
		meta = field.Metainfo()
	}
	var fncs = []func(iv T) error{}
	if meta.MaxUint.Valid {
		maxv := T(meta.MaxUint.V)
		fncs = append(fncs, func(iv T) error {
			if iv > maxv {
				return fmt.Errorf("int gt max")
			}
			return nil
		})
	}
	if meta.Func != nil {
		fncs = append(fncs, func(iv T) error { return meta.Func(iv) })
	}
	if len(fncs) < 1 {
		return nil, nil
	}

	do := func(iv T) error {
		for _, fn := range fncs {
			if err := fn(iv); err != nil {
				return err
			}
		}
		return nil
	}
	return func(uptr unsafe.Pointer) error {
			iv := *((*T)(unsafe.Add(uptr, field.Offset())))
			return do(iv)
		}, func(val any) error {
			return do(val.(T))
		}
}

func makeVldFunction(field *lion.Field[VldFieldMeta], meta *VldFieldMeta, gotype reflect.Type) (func(uptr unsafe.Pointer) error, func(val any) error) {
	_rawmeta := meta
	if meta == nil {
		meta = field.Metainfo()
	}
	getter := field.Getter()
	switch gotype.Kind() {
	case reflect.String:
		{
			var fncs = []func(v string) error{}

			if meta.MaxLength.Valid {
				maxl := meta.MaxLength.V
				fncs = append(fncs, func(v string) error {
					if len(v) > maxl {
						return fmt.Errorf("string too long")
					}
					return nil
				})
			}
			if meta.MinLength.Valid {
				minl := meta.MinLength.V
				fncs = append(fncs, func(v string) error {
					if len(v) < minl {
						return fmt.Errorf("string too short")
					}
					return nil
				})
			}
			if meta.Regexp != nil {
				fncs = append(fncs, func(v string) error {
					if !meta.Regexp.MatchString(v) {
						return fmt.Errorf("string not match")
					}
					return nil
				})
			}
			if len(meta.StringRanges) > 0 {
				if len(meta.StringRanges) > 15 {
					var rangemap = map[string]struct{}{}
					for _, rv := range meta.StringRanges {
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
						for _, rv := range meta.StringRanges {
							if v == rv {
								return nil
							}
						}
						return fmt.Errorf("string not in range")
					})
				}
			}
			if meta.Func != nil {
				fncs = append(fncs, func(v string) error { return meta.Func(v) })
			}
			if len(fncs) < 1 {
				return nil, nil
			}

			do := func(fv string) error {
				for _, fn := range fncs {
					if err := fn(fv); err != nil {
						return err
					}
				}
				return nil
			}
			return func(uptr unsafe.Pointer) error {
					fv := *((*string)(unsafe.Add(uptr, field.Offset())))
					return do(fv)
				}, func(val any) error {
					return do(val.(string))
				}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		{
			switch gotype.Bits() {
			case 8:
				return makeIntVld[int8](field, _rawmeta)
			case 16:
				return makeIntVld[int16](field, _rawmeta)
			case 32:
				return makeIntVld[int32](field, _rawmeta)
			case 64:
				return makeIntVld[int64](field, _rawmeta)
			}
			break
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		{
			switch gotype.Bits() {
			case 8:
				return makeUintVld[uint8](field, _rawmeta)
			case 16:
				return makeUintVld[uint16](field, _rawmeta)
			case 32:
				return makeUintVld[uint32](field, _rawmeta)
			case 64:
				return makeUintVld[uint64](field, _rawmeta)
			}
			break
		}
	case reflect.Slice, reflect.Array:
		{
			var slicefncs = []func(sv reflect.Value) error{}
			if gotype.Kind() == reflect.Slice {
				if meta.MaxSize.Valid {
					maxs := meta.MaxSize.V
					slicefncs = append(slicefncs, func(sv reflect.Value) error {
						if sv.Len() > maxs {
							return fmt.Errorf("slice too long")
						}
						return nil
					})
				}
				if meta.MinSize.Valid {
					mins := meta.MinSize.V
					slicefncs = append(slicefncs, func(sv reflect.Value) error {
						if sv.Len() < mins {
							return fmt.Errorf("slice too short")
						}
						return nil
					})
				}
			}
			_, eleanyfnc := makeVldFunction(field, _rawmeta, gotype.Elem())
			if eleanyfnc != nil {
				slicefncs = append(slicefncs, func(sv reflect.Value) error {
					slen := sv.Len()
					for i := 0; i < slen; i++ {
						elev := sv.Index(i)
						if ee := eleanyfnc(elev.Interface()); ee != nil {
							return ee
						}
					}
					return nil
				})
			}
			if len(slicefncs) < 1 {
				return nil, nil
			}
			do := func(sv reflect.Value) error {
				if !sv.IsValid() || sv.IsNil() {
					if meta.Optional {
						return nil
					}
					return fmt.Errorf("missing required")
				}
				for _, fn := range slicefncs {
					if err := fn(sv); err != nil {
						return err
					}
				}
				return nil
			}
			return func(uptr unsafe.Pointer) error { return do(reflect.ValueOf(getter(uptr))) },
				func(val any) error { return do(reflect.ValueOf(val)) }
		}
	case reflect.Map:
		{
			var mapfncs = []func(sv reflect.Value) error{}
			if meta.MaxSize.Valid {
				maxs := meta.MaxSize.V
				mapfncs = append(mapfncs, func(sv reflect.Value) error {
					if sv.Len() > maxs {
						return fmt.Errorf("map too long")
					}
					return nil
				})
			}
			if meta.MinSize.Valid {
				mins := meta.MinSize.V
				mapfncs = append(mapfncs, func(sv reflect.Value) error {
					if sv.Len() < mins {
						return fmt.Errorf("map too short")
					}
					return nil
				})
			}
			if gotype.Elem().Kind() == reflect.Map {
				panic(fmt.Errorf("unsupported nested map"))
			}
			_, eleanyfnc := makeVldFunction(field, _rawmeta, gotype.Elem())
			var keyanyfnc func(val any) error
			if meta.MapKey != nil {
				_, keyanyfnc = makeVldFunction(field, meta.MapKey, gotype.Key())
			}
			if eleanyfnc != nil {
				if keyanyfnc != nil {
					mapfncs = append(mapfncs, func(sv reflect.Value) error {
						iter := sv.MapRange()
						for iter.Next() {
							key := iter.Key()
							if ke := keyanyfnc(key.Interface()); ke != nil {
								return ke
							}
							val := iter.Value()
							if ve := eleanyfnc(val.Interface()); ve != nil {
								return ve
							}
						}
						return nil
					})
				} else {
					mapfncs = append(mapfncs, func(sv reflect.Value) error {
						iter := sv.MapRange()
						for iter.Next() {
							val := iter.Value()
							if ve := eleanyfnc(val.Interface()); ve != nil {
								return ve
							}
						}
						return nil
					})
				}
			} else {
				if keyanyfnc != nil {
					mapfncs = append(mapfncs, func(sv reflect.Value) error {
						iter := sv.MapRange()
						for iter.Next() {
							key := iter.Key()
							if ke := keyanyfnc(key.Interface()); ke != nil {
								return ke
							}
						}
						return nil
					})
				}
			}

			if len(mapfncs) < 1 {
				return nil, nil
			}
			do := func(mapv reflect.Value) error {
				if !mapv.IsValid() || mapv.IsNil() {
					if meta.Optional {
						return nil
					}
					return fmt.Errorf("missing required")
				}
				for _, fnc := range mapfncs {
					if err := fnc(mapv); err != nil {
						return err
					}
				}
				return nil
			}
			return func(uptr unsafe.Pointer) error { return do(reflect.ValueOf(getter(uptr))) },
				func(val any) error { return do(reflect.ValueOf(val)) }
		}
	case reflect.Pointer:
		{
			_, anyfnc := makeVldFunction(field, _rawmeta, gotype.Elem())
			if anyfnc == nil {
				return nil, nil
			}
			do := func(pv reflect.Value) error {
				if !pv.IsValid() || pv.IsNil() {
					if meta.Optional {
						return nil
					}
					return fmt.Errorf("missing required")
				}
				return anyfnc(pv.Elem().Interface())
			}
			return func(uptr unsafe.Pointer) error { return do(reflect.ValueOf(getter(uptr))) },
				func(val any) error { return do(reflect.ValueOf(val)) }
		}
	default:
		{
			switch gotype {
			case typeofTime:
				{
					var fncs = []func(tv time.Time) error{}
					if meta.MaxTime.Valid {
						fncs = append(fncs, func(tv time.Time) error {
							if tv.After(meta.MaxTime.V) {
								return fmt.Errorf("time too late")
							}
							return nil
						})
					}
					if meta.MinTime.Valid {
						fncs = append(fncs, func(tv time.Time) error {
							if tv.Before(meta.MinTime.V) {
								return fmt.Errorf("time too early")
							}
							return nil
						})
					}
					if len(fncs) < 1 {
						return nil, nil
					}

					do := func(tv time.Time) error {
						for _, fn := range fncs {
							if err := fn(tv); err != nil {
								return err
							}
						}
						return nil
					}

					return func(uptr unsafe.Pointer) error { return do(*((*time.Time)(unsafe.Add(uptr, field.Offset())))) },
						func(val any) error { return do(val.(time.Time)) }
				}
			}
		}
	}
	panic(fmt.Errorf("unsupported type: %s", gotype))
}

var (
	typeofTime = lion.Typeof[time.Time]()
)

func Vld[T any](ptr *T) error {
	gotype := lion.Typeof[T]()
	scheme, ok := schemes[gotype]
	if !ok {
		return fmt.Errorf("faceless.void.vld: can not found scheme for type: %s", gotype)
	}
	uptr := unsafe.Pointer(ptr)
	for _, vld := range scheme.vlds {
		if err := vld(uptr); err != nil {
			return err
		}
	}
	return nil
}
