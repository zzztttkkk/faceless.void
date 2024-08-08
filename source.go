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

//goland:noinspection ALL
func SOURCE_DIR() string {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		panic(fmt.Errorf("could not get the caller information"))
	}
	return path.Dir(file)
}
