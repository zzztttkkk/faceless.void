package blogs

import (
	"context"

	fv "github.com/zzztttkkk/faceless.void"
)

type PostParams struct {
	author  string
	title   string
	summary string
	content string
}

type PostResult struct {
	id string
}

func Post(ctx context.Context, params *PostParams) (*PostParams, error) {
	return nil, nil
}

func init() {
	fv.RegisterHttpEndpoint(Post)
}
