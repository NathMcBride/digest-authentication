package contexts

import (
	"context"

	"github.com/NathMcBride/digest-authentication/src/authentication/authenticator"
)

type SessionContextType string

const SessionCtxKey SessionContextType = "session"

func WithSession(ctx context.Context, session *authenticator.Session) context.Context {
	return context.WithValue(ctx, SessionCtxKey, session)
}

func GetSession(ctx context.Context) *authenticator.Session {
	session, ok := ctx.Value(SessionCtxKey).(*authenticator.Session)
	if !ok {
		return nil
	}

	return session
}
