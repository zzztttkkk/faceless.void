package account

import (
	"context"
	"hello/sitea/modules/internal/evts"

	fv "github.com/zzztttkkk/faceless.void"
)

type RegisterParams struct {
	Name     string
	Email    string
	Password string
	ExtInfo  struct {
		A string `json:"a"`
		B string `json:"b"`
	}
}

type RegisterResult struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func register(ctx context.Context, params *RegisterParams) (*RegisterResult, error) {
	// skip logics
	evts.EmitOnUserCreated(evts.EvtOnUserCreated{Uid: "0.0"})
	return nil, nil
}

func init() {
	fv.Endpoint().Func(fv.MakeHandleFunc(register))
}
