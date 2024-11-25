package fv

import (
	"context"
	"fmt"
	"net/http"
)

type HttpError struct {
	code    int
	reason  string
	message any
}

// Error implements error.
func (err HttpError) Error() string {
	return fmt.Sprintf(`http error: (%d %s), %v`, err.code, err.reason, err.message)
}

var _ error = (*HttpError)(nil)

func bindError(ctx context.Context, fieldname string) error {
	return HttpError{
		code:    http.StatusBadRequest,
		message: ``,
	}
}
