package internal

import (
	"reflect"
	"runtime"
)

func FuncName(rv reflect.Value) string {
	ptr := rv.Pointer()
	return runtime.FuncForPC(ptr).Name()
}
