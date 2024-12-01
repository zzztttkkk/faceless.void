package fv

import (
	"fmt"
	"testing"

	"github.com/zzztttkkk/faceless.void/internal"
)

func TestEnv(t *testing.T) {
	kvs := internal.Must(LoadEnv("./a.env"))
	for k, v := range kvs {
		fmt.Println("--------------------BEGIN----------------")
		fmt.Println(k)
		fmt.Println("--------------------")
		fmt.Println(v)
		fmt.Println("--------------------END-----------------")
	}
}
