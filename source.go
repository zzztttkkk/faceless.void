package fv

import (
	"fmt"
	"path"
	"runtime"
)

func SOURCE() string {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		panic(fmt.Errorf("could not get the caller information"))
	}
	return file
}

func SOURCE_DIR() string {
	return path.Dir(SOURCE())
}
