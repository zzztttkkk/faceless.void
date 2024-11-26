package vld

import (
	"context"

	"github.com/zzztttkkk/faceless.void/i18n"
)

type ErrorKind int

const (
	ErrorKindIntLtMin = ErrorKind(iota)
	ErrorKindIntGtMax
	ErrorKindCustomFunc
	ErrorKindStringLengthLtMin
	ErrorKindStringLengthGtMax
	ErrorKindStringNotMatchRegexp
	ErrorKindStringNotInEnums
)

type vldError struct {
	kind ErrorKind
	args []any
	msg  string
}

func (err vldError) Error() string {
	return err.msg
}

func newerror(ctx context.Context, kind ErrorKind, i18n *i18n.String, args ...any) vldError {
	return vldError{
		kind: kind,
		msg:  i18n.Format(ctx, args...),
		args: args,
	}
}
