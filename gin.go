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
	Methods string
	Path    string
	Fnc     gin.HandlerFunc
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

func (g *HttpGroup) Register(methods string, path string, fnc any) *HttpGroup {
	fv := reflect.ValueOf(fnc)
	if fv.Kind() != reflect.Func {
		panic(g.mkerr("`fnc` is not a function, %v", fnc))
	}

	ft := fv.Type()
	if ft.NumIn() < 1 || ft.In(0) != ctxInterfaceType {
		panic(g.mkerr("`fnc` must accept at least one argument, and the first one must be `context.Context`"))
	}
	if ft.NumOut() < 1 || ft.NumOut() > 2 {
		panic(g.mkerr("`fnc` must return one or two values, and the second return value must be `error`"))
	}

	argTypes := make([]reflect.Type, 0)
	for i := 1; i < ft.NumIn(); i++ {
		at := ft.In(i)
		if at.Kind() != reflect.Ptr {
			panic(g.mkerr(badArgsMsg))
		}
		at = at.Elem()
		if at.Kind() != reflect.Struct {
			panic(g.mkerr(badArgsMsg))
		}
		argTypes = append(argTypes, at)
	}

	exec := func(c *gin.Context) []reflect.Value {
		args := make([]reflect.Value, ft.NumIn())
		args[0] = reflect.ValueOf(c)
		for idx, at := range argTypes {
			av := reflect.New(at)
			if e := c.Bind(av.Interface()); e != nil {
				panic(e)
			}
			args[idx+1] = av
		}
		return fv.Call(args)
	}

	var handleFnc gin.HandlerFunc
	if ft.NumOut() == 2 {
		handleFnc = func(c *gin.Context) {
			rvs := exec(c)
			outv, errv := rvs[0], rvs[1]
			if errv.IsNil() {
				ginRespAny(c, errv)
			} else {
				ginRespAny(c, outv)
			}
		}
	} else {
		handleFnc = func(c *gin.Context) {
			out := exec(c)[0]
			ginRespAny(c, out)
		}
	}

	g.fncs = append(g.fncs, httpFunc{Methods: methods, Path: path, Fnc: handleFnc})
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
