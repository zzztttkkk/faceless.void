package vld

import (
	"context"
	"fmt"
	"testing"
)

func TestInts(t *testing.T) {
	err := Integer[uint64]().MaxValue(12).Func()(context.Background(), 34)
	fmt.Println(err)
}
