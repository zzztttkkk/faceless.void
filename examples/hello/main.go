package main

import (
	fv "github.com/zzztttkkk/faceless.void"
)

func main() {
	fv.RunHTTP(
		main,
		fv.HttpSite{
			Port: 8080, EndpointsGlob: "sitea/modules/**/*.go",
		},
	)
}
