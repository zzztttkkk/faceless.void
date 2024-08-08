package fv

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zzztttkkk/faceless.void/internal"
	"net/http"
	"reflect"
)

type httpFunc struct {
	Methods    string
	Path       string
	Fnc        gin.HandlerFunc
	InputTypes []reflect.Type
	OutputType reflect.Type
}

type HttpGroup struct {
	dir        string
	fncs       []httpFunc
	middleware []gin.HandlerFunc
}

func NewHttpGroup(dir string) *HttpGroup {
	return &HttpGroup{dir: dir}
}

func (g *HttpGroup) Use(middleware ...gin.HandlerFunc) *HttpGroup {
	g.middleware = append(g.middleware, middleware...)
	return g
}

var (
	ctxInterfaceType = reflect.TypeOf((*context.Context)(nil)).Elem()
)

func (g *HttpGroup) mkerr(v string, args ...any) error {
	return internal.ErrNamespace{Namespace: fmt.Sprintf("Gin.HttpGroup(%s)", g.dir)}.Errorf(v, args...)
}

var (
	badArgsMsg = "the arguments of `fnc`, except the first one which must be ctx, all others must be pointer to struct"
)

func (g *HttpGroup) Register(methods string, path string, fnc gin.HandlerFunc, inputTypes []reflect.Type, outType reflect.Type) *HttpGroup {
	g.fncs = append(g.fncs, httpFunc{
		Methods:    methods,
		Path:       path,
		Fnc:        fnc,
		InputTypes: inputTypes,
		OutputType: outType,
	})
	return g
}

type STRING string

func (s STRING) ResponseTo(ctx *gin.Context) {
	ctx.String(http.StatusOK, string(s))
}

type FILE string

func (f FILE) ResponseTo(ctx *gin.Context) {
	ctx.File(string(f))
}

type IGinResponse interface {
	ResponseTo(ctx *gin.Context)
}

func ginRespAny(ctx *gin.Context, v reflect.Value) {
	vt := v.Type()
	if vt.Kind() == reflect.Ptr {
		vt = vt.Elem()
		v = v.Elem()
	}
	switch vt.Kind() {
	case reflect.Func, reflect.Pointer, reflect.Chan:
		{
			panic(fmt.Errorf("fv.gin: can not cast `%s` to response", v.Interface()))
		}
	default:
		{
		}
	}

	iv := v.Interface()

	switch tv := iv.(type) {
	case IGinResponse:
		{
			tv.ResponseTo(ctx)
			return
		}
	case error:
		{
			// goto the default recover
			panic(tv)
		}
	default:
		{
			ctx.JSON(http.StatusOK, tv)
			return
		}
	}
}
