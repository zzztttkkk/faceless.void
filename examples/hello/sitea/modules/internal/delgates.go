package internal

import (
	"context"

	fv "github.com/zzztttkkk/faceless.void"
)

var (
	Delgates = struct {
		Register fv.Dalgate[func(ctx context.Context, email string, pwd string) (string, error)]
	}{}
)
