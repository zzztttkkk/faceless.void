package fv

import (
	"context"
	"fmt"
	"reflect"
	"testing"
)

func TestEnv(t *testing.T) {
	v, e := parseOneEnvValue(context.Background(), reflect.TypeOf([]int{}), "[1, 23]", false)
	if e != nil {
		t.Fatal(e)
	}
	ptr, ok := v.Interface().([]int)
	if !ok {
		t.Fatal(v)
	}
	fmt.Println(ptr)
}
