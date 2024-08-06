package fv_test

import (
	"fmt"
	"testing"

	fv "github.com/zzztttkkk/faceless.void"
)

type A struct{}

type B struct{}

type C struct{}

type D struct{}

type E struct{}

type F int

type G struct {
	V int
}

func TestDi(t *testing.T) {
	dic := fv.NewDIContainer()

	dic.Prepare(func() (*C, F, []fv.TokenValue[*G]) {
		return &C{}, F(12), []fv.TokenValue[*G]{
			fv.NewTokenValue("12", &G{V: 34}),
			fv.NewTokenValue("13", &G{V: 45}),
			fv.NewTokenValue("14", &G{V: 56}),
		}
	})

	dic.Register(func(a *A, c *C) (*B, *D) {
		fmt.Println("B")
		return &B{}, &D{}
	})

	dic.Register(func() *A {
		fmt.Println("A")
		return &A{}
	})

	dic.Register(func(d *D, a *A) {
		fmt.Println("D 1")
	})

	dic.Register(func(d *D) {
		fmt.Println("D 2")
	})

	dic.Register(func(f F, gtvg fv.TokenValueGetter[*G]) {
		g12 := gtvg.Get("12")
		g13 := gtvg.Get("13")
		fmt.Println("F", f, g12, g13)
	})

	dic.Run()
}
