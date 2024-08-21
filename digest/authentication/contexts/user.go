package contexts

import (
	"context"

	"github.com/NathMcBride/web-authentication/digest/authentication/authenticator"
)

type userCtxKeyType string

const userCtxKey userCtxKeyType = "user"

func WithUser(ctx context.Context, user *authenticator.User) context.Context {
	return context.WithValue(ctx, userCtxKey, user)
}

func GetUser(ctx context.Context) *authenticator.User {
	user, ok := ctx.Value(userCtxKey).(*authenticator.User)
	if !ok {
		return nil
	}

	return user
}
