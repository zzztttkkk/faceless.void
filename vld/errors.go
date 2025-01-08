package vld

import (
	"fmt"

	"github.com/zzztttkkk/lion"
	"github.com/zzztttkkk/lion/enums"
)

type ErrorKind int

const (
	ErrorKindMissingRequired ErrorKind = iota

	ErrorKindCustomFunc

	ErrorKindIntLtMin
	ErrorKindIntGtMax
	ErrorKindIntNotInRange

	ErrorKindStringTooLong
	ErrorKindStringTooShort
	ErrorKindStringNotMatched
	ErrorKindStringNotInRanges

	ErrorKindContainerSizeTooLarge
	ErrorKindContainerSizeTooSmall
)

var (
	AllErrorKinds []ErrorKind
)

func init() {
	enums.Generate(func() *enums.EnumOptions[ErrorKind] {
		return &enums.EnumOptions[ErrorKind]{
			NamePrefix:   "ErrorKind",
			GenAllSlice:  true,
			AllSliceName: "AllErrorKinds",
		}
	})
}

type Error struct {
	Field    *lion.Field[VldFieldMeta]
	Kind     ErrorKind
	BadValue any
	RawError error
}

var (
	_ error = (*Error)(nil)
)

func (e *Error) Error() string {
	return fmt.Sprintf(
		"fv.vld: %s.%s. %s, %v, %s",
		e.Field.Typeinfo().GoType.Name(), e.Field.StructField().Name,
		e.Kind, e.BadValue, e.RawError,
	)
}
