package internal

import (
	"reflect"
	"runtime"
)

func FuncName(fnc any) string {
	ptr := reflect.ValueOf(fnc).Pointer()
	return runtime.FuncForPC(ptr).Name()
}
