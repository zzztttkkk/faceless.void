package admin

import (
	"context"
	"hello/sitea/modules/internal"
)

func Create(ctx context.Context, email, pwd string) string {
	uid, err := internal.AccountDelegates.Register(ctx, email, pwd)
	if err != nil {
		panic(err)
	}
	return uid
}
