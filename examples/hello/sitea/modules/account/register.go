package account

import (
	"context"

	fv "github.com/zzztttkkk/faceless.void"
)

type RegisterParams struct {
	Name     string `bind:"optional"`
	Email    string `vld:"regexp<email>"`
	Password string `vld:"regexp<pwd>"`
}

type RegisterResult struct {
	Id   string
	Name string
}

func Register(ctx context.Context, params *RegisterParams) (*RegisterResult, error) {
	return nil, nil
}

func init() {
	fv.RegisterHttpEndpoint(
		Register,
		fv.EndpointOptions.Pattern("/register"),
	)
}
