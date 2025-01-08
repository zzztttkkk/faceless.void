package vld

import (
	"fmt"
	"strings"

	"github.com/zzztttkkk/lion"
	"github.com/zzztttkkk/lion/enums"
)

type ErrorKind int

const (
	ErrorKindMissingRequired ErrorKind = iota

	ErrorKindCustom

	ErrorKindIntLtMin
	ErrorKindIntGtMax
	ErrorKindIntNotInRange

	ErrorKindStringTooLong
	ErrorKindStringTooShort
	ErrorKindStringNotMatched
	ErrorKindStringNotInRanges

	ErrorKindContainerSizeTooLarge
	ErrorKindContainerSizeTooSmall

	ErrorKindNilPointer
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
	Meta     *VldFieldMeta
	Kind     ErrorKind
	BadValue any
	RawError error
}

var (
	_ error = (*Error)(nil)
)

func (e *Error) Error() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("fv.vld: %s.%s, Kind: %s", e.Field.Typeinfo().GoType.Name(), e.Field.StructField().Name, e.Kind))
	if e.BadValue != nil {
		sb.WriteString(fmt.Sprintf(", BadValue: %v", e.BadValue))
	}
	if e.RawError != nil {
		sb.WriteString(fmt.Sprintf(", RawErr: %e", e.RawError))
	}
	return sb.String()
}

func newerr(field *lion.Field[VldFieldMeta], meta *VldFieldMeta, kind ErrorKind) *Error {
	return &Error{
		Field: field,
		Meta:  meta,
		Kind:  kind,
	}
}

func (err *Error) with(bv any, re error) *Error {
	err.BadValue = bv
	err.RawError = re
	return err
}

func (err *Error) withbv(bv any) *Error {
	err.BadValue = bv
	return err
}

func (err *Error) withre(re error) *Error {
	err.RawError = re
	return err
}
