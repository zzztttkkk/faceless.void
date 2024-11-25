package vld

import (
	"fmt"
	"testing"
)

func TestInts(t *testing.T) {
	err := Integer[uint64]().MaxValue(12).Func()(34)
	fmt.Println(err)
}
