package vld

import (
	"github.com/zzztttkkk/faceless.void/internal"
)

const (
	ErrorKindVldIntLtMin = internal.ErrorKind(iota)
	ErrorKindVldIntGtMax
	ErrorKindVldCustomFunc
	ErrorKindVldStringLengthLtMin
	ErrorKindVldStringLengthGtMax
	ErrorKindVldStringNotMatchRegexp
	ErrorKindVldStringNotInEnums

	_MaxVldErrorKind
)

func init() {
	if _MaxVldErrorKind >= internal.MaxVldErrorKind {
		panic("vld error kind > internal.MaxVldErrorKind")
	}
}
