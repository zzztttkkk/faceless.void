package fv

import (
	"context"

	"github.com/zzztttkkk/faceless.void/internal"
)

type IHttpError interface {
	error
	StatusCode() int
	BodyMessage(ctx context.Context) []byte
}

type ErrorKind internal.ErrorKind
type InternalError internal.Error

const (
	ErrorKindBindingMissingRequired = ErrorKind(iota + internal.MaxVldErrorKind + 1)
	ErrorKindBindingParseFailed
	ErrorKindBindingUnmarshalFailed

	_MaxBindingErrorKind
)

func init() {
	if _MaxBindingErrorKind > ErrorKind(internal.MaxBindingErrorKind) {
		panic("binding error kind > internal.MaxBindingErrorKind")
	}
}

func NewError(ctx context.Context, kind ErrorKind, i18n *I18nString, args ...any) error {
	return internal.NewError(ctx, internal.ErrorKind(kind), (*internal.I18nString)(i18n), args...)
}
