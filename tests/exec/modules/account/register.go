package account

import "context"

type RegisterParams struct {
	Email    string `binding:"required,email"`
	Password string `binding:"required"`
}

// Register
// :fv: http post,put /register
func Register(ctx context.Context, params *RegisterParams) (*LoginResult, error) {
	return nil, nil
}
