package middleware_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/NathMcBride/digest-authentication/src/authentication/authenticator"
	"github.com/NathMcBride/digest-authentication/src/authentication/contexts"
	"github.com/NathMcBride/digest-authentication/src/authentication/middleware"
	. "github.com/NathMcBride/digest-authentication/src/authentication/middleware/fakes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Authenticate middleware", func() {
	Describe("RequireAuth", func() {
		var (
			recorder                *httptest.ResponseRecorder
			request                 *http.Request
			fakeNextHandler         *FakeHandler
			fakeUnauthorizedHandler *FakeUnauthorizedHandler
			fakeAuthenticator       *FakeAuthenticator
			requireAuthHandler      http.Handler
		)

		BeforeEach(func() {
			fakeUnauthorizedHandler = &FakeUnauthorizedHandler{}
			fakeAuthenticator = &FakeAuthenticator{}
			session := authenticator.Session{IsAuthenticated: true}
			fakeAuthenticator.AuthenticateReturns(session, nil)

			middleware := middleware.Authenticate{
				UnauthorizedHandler: fakeUnauthorizedHandler,
				Authenticator:       fakeAuthenticator,
			}

			fakeNextHandler = &FakeHandler{}
			requireAuthHandler = middleware.RequireAuth(fakeNextHandler)
			request = httptest.NewRequest("GET", "http://valid", nil)
			recorder = httptest.NewRecorder()
		})

		When("authentication succeeds", func() {
			It("calls the next handler", func() {
				requireAuthHandler.ServeHTTP(recorder, request)
				Expect(fakeNextHandler.ServeHTTPCallCount()).To(Equal(1))
			})

			It("stores the session in the context", func() {
				fakeAuthenticator.AuthenticateReturns(
					authenticator.Session{
						User:            authenticator.User{UserID: "a-test-user"},
						IsAuthenticated: true},
					nil,
				)

				requireAuthHandler.ServeHTTP(recorder, request.WithContext(context.Background()))

				Expect(fakeNextHandler.ServeHTTPCallCount()).To(Equal(1))
				_, r := fakeNextHandler.ServeHTTPArgsForCall(0)

				session, ok := r.Context().Value(contexts.SessionCtxKey).(*authenticator.Session)
				Expect(ok).To(BeTrue(), "Session not found")

				Expect(session.User.UserID).To(Equal("a-test-user"))
			})
		})

		When("authentication fails", func() {
			It("calls the unauthorized handler", func() {
				fakeAuthenticator.AuthenticateReturns(authenticator.Session{IsAuthenticated: false}, nil)
				requireAuthHandler.ServeHTTP(recorder, request)

				Expect(fakeUnauthorizedHandler.HandleUnauthorizedCallCount()).To(Equal(1))
			})

			It("the next handler is not called", func() {
				fakeAuthenticator.AuthenticateReturns(authenticator.Session{IsAuthenticated: false}, nil)
				requireAuthHandler.ServeHTTP(recorder, request)

				Expect(fakeNextHandler.ServeHTTPCallCount()).To(Equal(0))
			})
		})

		When("authentication errors", func() {
			It("returns 500 StatusInternalServerError", func() {
				fakeAuthenticator.AuthenticateReturns(authenticator.Session{}, fmt.Errorf("an-error"))
				requireAuthHandler.ServeHTTP(recorder, request)

				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
			})

			It("the next handler is not called", func() {
				fakeAuthenticator.AuthenticateReturns(authenticator.Session{}, fmt.Errorf("an-error"))
				requireAuthHandler.ServeHTTP(recorder, request)

				Expect(fakeNextHandler.ServeHTTPCallCount()).To(Equal(0))
			})
		})
	})

})
