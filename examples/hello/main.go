package main

import (
	fv "github.com/zzztttkkk/faceless.void"
)

func main() {
	fv.RunHTTP(8080, main, []string{"./modules/**/*.go"})
}
