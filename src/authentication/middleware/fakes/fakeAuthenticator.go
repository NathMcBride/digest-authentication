package fakes

import (
	"net/http"

	"github.com/NathMcBride/digest-authentication/src/authentication/authenticator"
)

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
