package fv

import (
	"context"
	"database/sql"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/zzztttkkk/faceless.void/internal"
	"github.com/zzztttkkk/reflectx"
)

type BindingSrcKind int

const (
	BindingSrcForm = BindingSrcKind(iota)
	BindingSrcHeader
	BindingSrcCookie
	BindingSrcQuery
	BindingSrcPostForm
	BindingSrcPathValue
)

type BindingOptions struct {
	Where      BindingSrcKind
	Alias      []string
	Optional   bool
	Default    sql.Null[any]
	Unmarshal  func(fptr any, data string) error
	IsFile     bool
	TimeLayout string
}

func init() {
	reflectx.RegisterOf[BindingOptions]().TagNames("bnd", "binding")
}

type _Binding[T any] struct {
	typeinfo *reflectx.TypeInfo[BindingOptions]
	binds    []func(ctx context.Context, ptr unsafe.Pointer) error
}

func Binding[T any]() *_Binding[T] {
	obj := &_Binding[T]{
		typeinfo: reflectx.TypeInfoOf[T, BindingOptions](),
	}
	for idx := range obj.typeinfo.Fields {
		obj.binds = append(obj.binds, makeBindFunc(&obj.typeinfo.Fields[idx]))
	}
	return obj
}

const (
	defaultMaxMemory int64 = 32 << 20 // 32 MB
)

type _ReqGetCache struct {
	query   url.Values
	cookies []*http.Cookie
}

func (rgc *_ReqGetCache) Query(req *http.Request) url.Values {
	if rgc.query == nil {
		rgc.query = req.URL.Query()
	}
	return rgc.query
}

func (rgc *_ReqGetCache) Cookies(req *http.Request) []*http.Cookie {
	if rgc.cookies == nil {
		rgc.cookies = req.Cookies()
		if rgc.cookies == nil {
			rgc.cookies = make([]*http.Cookie, 0, 1)
		}
	}
	return rgc.cookies
}

func reqGetCache(ctx context.Context) *_ReqGetCache {
	v := ctx.Value(internal.CtxKeyForBindingGetter)
	if v == nil {
		panic("empty request get cache")
	}
	return v.(*_ReqGetCache)
}

func makeGetter(field *reflectx.Field[BindingOptions]) func(ctx context.Context, req *http.Request) ([]string, bool) {
	if field.Meta == nil {
		field.Meta = &BindingOptions{}
	}

	parseForm := func(req *http.Request) {
		v := defaultMaxMemory
		_ = req.ParseMultipartForm(v)
	}

	var getonce func(ctx context.Context, req *http.Request, name string) ([]string, bool)

	switch field.Meta.Where {
	case BindingSrcForm:
		{
			getonce = func(_ context.Context, req *http.Request, name string) ([]string, bool) {
				parseForm(req)
				vs, ok := req.Form[name]
				return vs, ok
			}
		}
	case BindingSrcQuery:
		{
			getonce = func(ctx context.Context, req *http.Request, name string) ([]string, bool) {
				vs, ok := reqGetCache(ctx).Query(req)[name]
				return vs, ok
			}
		}
	case BindingSrcPostForm:
		{
			getonce = func(_ context.Context, req *http.Request, name string) ([]string, bool) {
				parseForm(req)
				vs, ok := req.PostForm[name]
				return vs, ok
			}
		}
	case BindingSrcPathValue:
		{
			getonce = func(ctx context.Context, req *http.Request, name string) ([]string, bool) {
				v := req.PathValue(name)
				if v != "" {
					return []string{v}, true
				}
				return nil, false
			}
		}
	case BindingSrcCookie:
		{
			getonce = func(ctx context.Context, req *http.Request, name string) ([]string, bool) {
				var vs []string
				for _, c := range reqGetCache(ctx).Cookies(req) {
					if c.Name == name {
						vs = append(vs, c.Value)
					}
				}
				if len(vs) > 0 {
					return vs, true
				}
				return nil, false
			}
		}
	case BindingSrcHeader:
		{
			getonce = func(_ context.Context, req *http.Request, name string) ([]string, bool) {
				vs, ok := req.Header[name]
				return vs, ok
			}
		}
	default:
		{
			panic("unreachable code")
		}
	}

	return func(ctx context.Context, req *http.Request) ([]string, bool) {
		vs, ok := getonce(ctx, req, field.Name)
		if ok {
			return vs, ok
		}
		for _, alias := range field.Meta.Alias {
			vs, ok = getonce(ctx, req, alias)
			if ok {
				return vs, ok
			}
		}
		return nil, false
	}
}

var (
	msgForBindingMissingRequired = NewI18nString("fv.binding: missing required. %s")
	msgForBadDefaultValueType    = "fv.binding: bad default value type, required type is `%s`. %s"
)

func makeBindFunc(field *reflectx.Field[BindingOptions]) func(ctx context.Context, ptr unsafe.Pointer) error {
	if field.Meta == nil {
		field.Meta = &BindingOptions{}
	}

	if field.Meta.IsFile {
		switch field.Field.Type {
		case reflectx.Typeof[[]byte]():
			{
				break
			}
		case reflectx.Typeof[[][]byte]():
			{
				break
			}
		case reflectx.Typeof[*multipart.FileHeader]():
			{
				break
			}
		case reflectx.Typeof[[]*multipart.FileHeader]():
			{
				break
			}
		default:
			{
				panic(fmt.Errorf("fv.binding: bad field type for file. %s.%s", field.Field.PkgPath, field.Field.Type))
			}
		}
	}

	ptrgetter := field.PtrGetter()
	vsgetter := makeGetter(field)

	ftype := field.Field.Type
	isslice := false
	if field.Field.Type.Kind() == reflect.Slice {
		isslice = true
		ftype = ftype.Elem()
	}

	switch ftype {
	case reflectx.Typeof[string]():
		{
			if isslice {
				var dv []string
				if field.Meta.Default.Valid {
					switch tv := field.Meta.Default.V.(type) {
					case string:
						{
							dv = []string{tv}
							break
						}
					case []string:
						{
							dv = tv
							break
						}
					default:
						{
							panic(fmt.Errorf(msgForBadDefaultValueType, "[]string", field.Name))
						}
					}
				}

				return func(ctx context.Context, ptr unsafe.Pointer) error {
					fptr := ptrgetter(ptr).(*[]string)
					vs, ok := vsgetter(ctx, HttpRequest(ctx))
					if !ok {
						if field.Meta.Default.Valid {
							*fptr = dv
							return nil
						}
						if field.Meta.Optional {
							return nil
						}
						return NewError(ctx, ErrorKindBindingMissingRequired, msgForBindingMissingRequired, field.Name)
					}
					*fptr = vs
					return nil
				}
			}

			var dv string
			if field.Meta.Default.Valid {
				sv, ok := field.Meta.Default.V.(string)
				if !ok {
					panic(fmt.Errorf(msgForBadDefaultValueType, "string", field.Name))
				}
				dv = sv
			}
			return func(ctx context.Context, ptr unsafe.Pointer) error {
				fptr := ptrgetter(ptr).(*string)
				vs, ok := vsgetter(ctx, HttpRequest(ctx))
				if !ok || len(vs) < 1 {
					if field.Meta.Default.Valid {
						*fptr = dv
						return nil
					}
					if field.Meta.Optional {
						return nil
					}
					return NewError(ctx, ErrorKindBindingMissingRequired, msgForBindingMissingRequired, field.Name)
				}
				*fptr = vs[0]
				return nil
			}
		}
	case reflectx.Typeof[bool]():
		{
			if isslice {
				var dv []bool
				if field.Meta.Default.Valid {
					switch tv := field.Meta.Default.V.(type) {
					case bool:
						{
							dv = []bool{tv}
							break
						}
					case []bool:
						{
							dv = tv
							break
						}
					}
				}
				return func(ctx context.Context, ptr unsafe.Pointer) error {
					fptr := ptrgetter(ptr).(*[]bool)

					vs, ok := vsgetter(ctx, HttpRequest(ctx))
					if !ok {
						if field.Meta.Default.Valid {
							*fptr = dv
							return nil
						}
						if field.Meta.Optional {
							return nil
						}
						return NewError(ctx, ErrorKindBindingMissingRequired, msgForBindingMissingRequired, field.Name)
					}
					*fptr = make([]bool, len(vs))
					for idx, v := range vs {
						bv, err := strconv.ParseBool(v)
						if err != nil {
							return err
						}
						(*fptr)[idx] = bv
					}
					return nil
				}
			}

			var dv bool
			if field.Meta.Default.Valid {
				bv, ok := field.Meta.Default.V.(bool)
				if !ok {
					panic(fmt.Errorf(msgForBadDefaultValueType, "bool", field.Name))
				}
				dv = bv
			}
			return func(ctx context.Context, ptr unsafe.Pointer) error {
				fptr := ptrgetter(ptr).(*bool)
				*fptr = dv
				return nil
			}
		}
	default:
		{
			if field.Meta.Unmarshal == nil {
				panic(fmt.Errorf(""))
			}

			if isslice {
				return func(ctx context.Context, ptr unsafe.Pointer) error {
					return nil
				}
			}
			return func(ctx context.Context, ptr unsafe.Pointer) error {
				return nil
			}
		}
	}
}

func (builder *_Binding[T]) Bind(ctx context.Context, ptr *T) error {
	for _, bind := range builder.binds {
		if err := bind(ctx, unsafe.Pointer(ptr)); err != nil {
			return err
		}
	}
	return nil
}
