package fv_test

import (
	"fmt"
	"testing"

	fv "github.com/zzztttkkk/faceless.void"
)

func TestSource(t *testing.T) {
	fmt.Println(fv.SOURCE(), fv.SOURCE_DIR())
}
