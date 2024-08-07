package fv

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

type httpFunc struct {
	Methods string
	Path    string
	Fnc     gin.HandlerFunc
	Inputs  []reflect.Type
	Output  reflect.Type
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

func (g *HttpGroup) Register(methods string, path string, fnc any) *HttpGroup {
	fv := reflect.ValueOf(fnc)
	if fv.Kind() != reflect.Func {
		panic("fv")
	}

	ft := fv.Type()
	if ft.NumIn() < 1 || ft.In(0) != ctxInterfaceType {
		panic("fv")
	}
	if ft.NumOut() != 1 {
		panic("fv")
	}

	argTypes := make([]reflect.Type, 0)
	for i := 1; i < ft.NumIn(); i++ {
		at := ft.In(i)
		if at.Kind() != reflect.Ptr {
			panic("fv")
		}
		at = at.Elem()
		if at.Kind() != reflect.Struct {
			panic("fv")
		}
		argTypes = append(argTypes, at)
	}

	handleFnc := func(c *gin.Context) {
		args := make([]reflect.Value, ft.NumIn())
		args[0] = reflect.ValueOf(c)
		for idx, at := range argTypes {
			av := reflect.New(at)
			if e := c.Bind(av.Interface()); e != nil {
				panic(e)
			}
			args[idx+1] = av
		}
		out := fv.Call(args)[0]
		ginRespAny(c, out)
	}
	g.fncs = append(g.fncs, httpFunc{Methods: methods, Path: path, Fnc: handleFnc})
	return g
}

type STRING string
type HTML string
type FILE string
type STATUS_CODE int

func ginRespAny(ctx *gin.Context, v reflect.Value) {
	vt := v.Type()
	if vt.Kind() == reflect.Ptr {
		vt = vt.Elem()
		v = v.Elem()
	}
	switch vt.Kind() {
	case reflect.Func, reflect.Pointer, reflect.Chan:
		{
			panic("")
		}
	default:
		{
		}
	}

	iv := v.Interface()

	switch tv := iv.(type) {
	case FILE:
		{
			ctx.File((string)(tv))
			return
		}
	case HTML:
		{
			ctx.HTML(http.StatusOK, (string)(tv), nil)
			return
		}
	case STRING:
		{
			ctx.String(http.StatusOK, (string)(tv))
			return
		}
	case STATUS_CODE:
		{
			ctx.Status(int(tv))
			return
		}
	default:
		{
			ctx.JSON(http.StatusOK, tv)
			return
		}
	}
}
