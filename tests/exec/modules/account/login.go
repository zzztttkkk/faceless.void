package account

import (
	"context"
)

type LoginParams struct {
	Email      string `binding:"required"`
	Password   string `binding:"required"`
	SecretCode string `binding:"required"`
}

type LoginResult struct {
	UserId   uint64
	Token    string
	ExpireAt int64
}

// Login
// :fv: http post,put /login
func Login(ctx context.Context, params *LoginParams) (*LoginResult, error) {
	return nil, nil
}
