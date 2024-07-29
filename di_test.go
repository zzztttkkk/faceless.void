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

func TestDi(t *testing.T) {
	dic := fv.NewDIContainer()

	dic.Pre(func() (*C, F) {
		return &C{}, F(12)
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

	dic.Register(func(f F) {
		fmt.Println("F", f)
	})

	dic.Run()
}
