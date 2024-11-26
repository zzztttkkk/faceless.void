package fv

import (
	"context"
)

type IHttpError interface {
	error
	StatusCode() int
	BodyMessage(ctx context.Context) []byte
}
