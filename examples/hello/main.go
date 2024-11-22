package main

import (
	_ "hello/sitea/modules/account"

	fv "github.com/zzztttkkk/faceless.void"
)

func main() {
	fv.RunHTTP(
		main,
		fv.HttpSite{
			Port: 8080, EndpointsGlob: "sitea/modules/**/*.go",
		},
		fv.HttpSite{
			Port: 8081, EndpointsGlob: "siteb/modules/**/*.go",
		},
	)
}
