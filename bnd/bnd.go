package bnd

import (
	"context"
	"net/http"
)

func Bind[T any](ctx context.Context, dst *T, req *http.Request) error {
	return nil
}
