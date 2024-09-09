package fakes

import (
	"github.com/NathMcBride/digest-authentication/src/authentication/model"
	"github.com/NathMcBride/digest-authentication/src/providers/credential"
)

type FakeDigest struct {
	callCount   int
	argsForCall []struct {
		credentials credential.Credentials
		authHeader  model.AuthHeader
		Method      string
	}
	calculateReturns struct {
		digest string
		err    error
	}
}

func (fd *FakeDigest) CalculateReturns(digest string, err error) {
	fd.calculateReturns.digest = digest
	fd.calculateReturns.err = err
}

func (fd *FakeDigest) CalculateArgsForCall(i int) (credential.Credentials, model.AuthHeader, string) {
	args := fd.argsForCall[i]
	return args.credentials, args.authHeader, args.Method
}

func (fd *FakeDigest) Calculate(credentials credential.Credentials, authHeader model.AuthHeader, Method string) (string, error) {
	fd.callCount++

	fd.argsForCall = append(fd.argsForCall,
		struct {
			credentials credential.Credentials
			authHeader  model.AuthHeader
			Method      string
		}{
			credentials: credentials,
			authHeader:  authHeader,
			Method:      Method})

	r := fd.calculateReturns
	return r.digest, r.err
}
