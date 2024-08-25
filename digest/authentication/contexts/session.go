package contexts

import (
	"context"

	"github.com/NathMcBride/web-authentication/digest/authentication/authenticator"
)

type sessionCtxKeyType string

const sessionCtxKey sessionCtxKeyType = "session"

func WithSession(ctx context.Context, session *authenticator.Session) context.Context {
	return context.WithValue(ctx, sessionCtxKey, session)
}

func GetSession(ctx context.Context) *authenticator.Session {
	session, ok := ctx.Value(sessionCtxKey).(*authenticator.Session)
	if !ok {
		return nil
	}

	return session
}
