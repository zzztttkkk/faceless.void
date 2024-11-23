package vld

import (
	"fmt"
	"testing"
)

func TestInts(t *testing.T) {
	err := Integer[uint64]().MaxValue(12).Finish()(34)
	fmt.Println(err)
}
