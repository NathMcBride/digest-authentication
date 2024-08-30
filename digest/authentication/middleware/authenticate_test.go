package middleware_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/NathMcBride/web-authentication/digest/authentication/authenticator"
	"github.com/NathMcBride/web-authentication/digest/authentication/contexts"
	"github.com/NathMcBride/web-authentication/digest/authentication/middleware"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type ServeHTTPArgs struct {
	writer  http.ResponseWriter
	request *http.Request
}

type FakeHandler struct {
	callCount   int
	argsForCall []ServeHTTPArgs
}

func (fh *FakeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fh.callCount++
	fh.argsForCall = append(fh.argsForCall, ServeHTTPArgs{writer: w, request: r})
}

type FakeUnauthorizedHandler struct {
	callCount   int
	argsForCall []ServeHTTPArgs
}

func (fu *FakeUnauthorizedHandler) HandleUnauthorized(w http.ResponseWriter, r *http.Request) {
	fu.callCount++
	fu.argsForCall = append(fu.argsForCall, ServeHTTPArgs{writer: w, request: r})
}

type FakeAuthenticator struct {
	callCount           int
	argsForCall         []*http.Request
	authenticateReturns struct {
		session authenticator.Session
		err     error
	}
}

func (fa *FakeAuthenticator) AuthenticateReturns(session authenticator.Session, err error) {
	fa.authenticateReturns.session = session
	fa.authenticateReturns.err = err
}

func (fa *FakeAuthenticator) Authenticate(r *http.Request) (authenticator.Session, error) {
	fa.callCount++
	fa.argsForCall = append(fa.argsForCall, r)
	return fa.authenticateReturns.session, fa.authenticateReturns.err
}

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
				Expect(fakeNextHandler.callCount).To(Equal(1))
			})

			It("stores the session in the context", func() {
				fakeAuthenticator.AuthenticateReturns(
					authenticator.Session{
						User:            authenticator.User{UserID: "a-test-user"},
						IsAuthenticated: true},
					nil,
				)

				requireAuthHandler.ServeHTTP(recorder, request.WithContext(context.Background()))

				Expect(fakeNextHandler.callCount).To(Equal(1))
				r := fakeNextHandler.argsForCall[0].request

				session, ok := r.Context().Value(contexts.SessionCtxKey).(*authenticator.Session)
				Expect(ok).To(BeTrue(), "Session not found")

				Expect(session.User.UserID).To(Equal("a-test-user"))
			})
		})

		When("authentication fails", func() {
			It("calls the unauthorized handler", func() {
				fakeAuthenticator.AuthenticateReturns(authenticator.Session{IsAuthenticated: false}, nil)
				requireAuthHandler.ServeHTTP(recorder, request)

				Expect(fakeUnauthorizedHandler.callCount).To(Equal(1))
			})

			It("the next handler is not called", func() {
				fakeAuthenticator.AuthenticateReturns(authenticator.Session{IsAuthenticated: false}, nil)
				requireAuthHandler.ServeHTTP(recorder, request)

				Expect(fakeNextHandler.callCount).To(Equal(0))
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

				Expect(fakeNextHandler.callCount).To(Equal(0))
			})
		})
	})

})
