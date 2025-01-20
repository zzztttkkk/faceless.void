package vld

import (
	"fmt"
	"strings"
	"unsafe"

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

	ErrorKindTimeTooEarly
	ErrorKindTimeTooLate

	ErrorKindStringTooLong
	ErrorKindStringTooShort
	ErrorKindStringNotMatched
	ErrorKindStringNotInRanges

	ErrorKindContainerSizeTooLarge
	ErrorKindContainerSizeTooSmall

	ErrorKindNilPointer
	ErrorKindNilSlice
	ErrorKindNilMap
)

var (
	AllErrorKinds []ErrorKind
)

func init() {
	enums.Generate(func() *enums.Options[ErrorKind] {
		return &enums.Options[ErrorKind]{
			RemoveCommonPrefix: true,
			AllSlice:           true,
			AllSliceName:       "AllErrorKinds",
			NameOverwrites:     map[ErrorKind]string{},
		}
	})
}

type Error struct {
	Fields   []unsafe.Pointer
	Meta     unsafe.Pointer
	Kind     ErrorKind
	BadValue any
	RawError error
}

var (
	_ error = (*Error)(nil)
)

func (e *Error) Error() string {
	var sb strings.Builder
	sb.WriteString("fv.vld: [")
	isfirst := true
	for i := len(e.Fields) - 1; i >= 0; i-- {
		fuptr := e.Fields[i]
		fptr := (*lion.Field)(fuptr)

		if isfirst {
			isfirst = false
			pkgpath := fptr.TypeInfo().GoType.PkgPath()
			idx := strings.LastIndexByte(pkgpath, '/')
			sb.WriteString(pkgpath[idx+1:])
			sb.WriteByte('.')
			sb.WriteString(fptr.TypeInfo().GoType.Name())
		} else {
			sb.WriteByte(' ')
			sb.WriteString(fptr.TypeInfo().GoType.Name())
		}
		sb.WriteByte('.')
		sb.WriteString(fptr.StructField().Name)
	}

	sb.WriteString("] Kind: ")
	sb.WriteString(e.Kind.String())

	if e.BadValue != nil {
		sb.WriteString(" BadValue: ")
		sb.WriteString(fmt.Sprintf("%v", e.BadValue))
	}

	if e.RawError != nil {
		sb.WriteString(" RawError: ")
		sb.WriteString(fmt.Sprintf("%e", e.RawError))
	}
	return sb.String()
}

func newerr(field *lion.Field, meta *VldFieldMeta, kind ErrorKind) *Error {
	fs := make([]unsafe.Pointer, 0, 4)
	fs = append(fs, unsafe.Pointer(field))
	return &Error{
		Fields: fs,
		Meta:   unsafe.Pointer(meta),
		Kind:   kind,
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

func (err *Error) appendfield(filed *lion.Field) *Error {
	err.Fields = append(err.Fields, unsafe.Pointer(filed))
	return err
}
