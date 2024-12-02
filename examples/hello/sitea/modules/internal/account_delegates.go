package internal

import (
	"context"
)

var (
	AccountDelegates = struct {
		Register func(ctx context.Context, email string, pwd string) (string, error)
	}{}
)
