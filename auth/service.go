package auth

import "context"

type Service interface {
	SignIn(ctx context.Context, username string, password string) (string, error)
	SignUp(ctx context.Context, username string, password string) (string, error)
}
