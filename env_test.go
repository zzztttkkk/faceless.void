package fv

import (
	"fmt"
	"testing"
)

func TestEnv(t *testing.T) {
	kvs := Must(LoadEnv("./a.env"))
	for k, v := range kvs {
		fmt.Println("--------------------BEGIN----------------")
		fmt.Println(k)
		fmt.Println("--------------------")
		fmt.Println(v)
		fmt.Println("--------------------END-----------------")
	}
}
