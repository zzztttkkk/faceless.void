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
	container := fv.NewContainer[string]()

	container.Register(func() (*C, F, []fv.TokenValue[string]) {
		return &C{}, F(12), []fv.TokenValue[string]{
			fv.ValueWithToken("12", &G{V: 34}),
			fv.ValueWithToken("13", &G{V: 45}),
			fv.ValueWithToken("14", &G{V: 56}),
		}
	}).Register(func(a *A, c *C) (*B, *D) {
		fmt.Println("require *A, *C, provide *B *D")
		return &B{}, &D{}
	}).Register(func() *A {
		fmt.Println("A")
		return &A{}
	}).Register(func(d *D, a *A) {
		fmt.Println("D 1")
	}).Register(func(d *D) {
		fmt.Println("D 2")
	}).Register(func(f F) {
		var g12, g13 *G
		container.
			GetByToken(&g12, "12").
			GetByToken(&g13, "13")

		fmt.Println("F", f, g12, g13)
	}).
		Run()
}
