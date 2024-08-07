package fv

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
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
	if ft.NumOut() < 1 || ft.Out(0) != ctxInterfaceType {
		panic("fv")
	}
	// todo get doc info

	handleFnc := func(c *gin.Context) {
		// todo bind args
		args := make([]reflect.Value, ft.NumIn())
		args[0] = reflect.ValueOf(c)
		outs := fv.Call(args)
		fmt.Println(outs)
		// todo write response
	}
	g.fncs = append(g.fncs, httpFunc{Methods: methods, Path: path, Fnc: handleFnc})
	return g
}
