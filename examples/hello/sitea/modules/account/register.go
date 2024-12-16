package account

import (
	"context"
	"hello/sitea/modules/internal"
	"hello/sitea/modules/internal/evts"
	"reflect"

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

func init() {
	fv.RegisterTypes(reflect.TypeOf(RegisterParams{}))
}

type RegisterResult struct {
	Id   string
	Name string
}

func Register(ctx context.Context, params *RegisterParams) (*RegisterResult, error) {
	// skip logics
	evts.EmitOnUserCreated(evts.EvtOnUserCreated{Uid: "0.0"})
	return nil, nil
}

func init() {
	internal.AccountDelegates.Register = func(ctx context.Context, email, pwd string) (string, error) {
		var params RegisterParams
		params.Email = email
		params.Password = pwd
		r, e := Register(ctx, &params)
		if e != nil {
			return "", nil
		}
		return r.Id, nil
	}
}
