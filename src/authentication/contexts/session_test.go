package contexts_test

import (
	"context"

	"github.com/NathMcBride/digest-authentication/src/authentication/authenticator"
	"github.com/NathMcBride/digest-authentication/src/authentication/contexts"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Session", func() {

	It("adds a session to the context", func() {
		ctx := context.Background()
		session := authenticator.Session{User: authenticator.User{UserID: "a-user-id"}, IsAuthenticated: true}

		newCtx := contexts.WithSession(ctx, &session)

		value, ok := newCtx.Value(contexts.SessionCtxKey).(*authenticator.Session)
		Expect(ok).To(BeTrue())
		Expect(*value).To(Equal(session))
	})

	It("gets a session from the context", func() {
		session := authenticator.Session{User: authenticator.User{UserID: "a-user-id"}, IsAuthenticated: true}
		ctx := context.WithValue(context.Background(), contexts.SessionCtxKey, &session)

		value := contexts.GetSession(ctx)

		Expect(value).ToNot(BeNil())
		Expect(*value).To(Equal(session))
	})
})
