package fakes

import "github.com/NathMcBride/digest-authentication/src/providers/credential"

type FakeCredentialProvider struct {
	callCount   int
	argsForCall []struct {
		userID  string
		useHash bool
	}
	getCredentialsReturns struct {
		credentials *credential.Credentials
		found       bool
		err         error
	}
}

func (fcp *FakeCredentialProvider) GetCredentialsReturns(credentials *credential.Credentials, found bool, err error) {
	fcp.getCredentialsReturns.credentials = credentials
	fcp.getCredentialsReturns.found = found
	fcp.getCredentialsReturns.err = err
}

func (fcp *FakeCredentialProvider) GetCredentialsArgsForCall(i int) (string, bool) {
	args := fcp.argsForCall[i]
	return args.userID, args.useHash
}

func (fcp *FakeCredentialProvider) GetCredentials(userID string, useHash bool) (*credential.Credentials, bool, error) {
	fcp.callCount++

	fcp.argsForCall = append(fcp.argsForCall,
		struct {
			userID  string
			useHash bool
		}{
			userID:  userID,
			useHash: useHash})

	r := fcp.getCredentialsReturns
	return r.credentials, r.found, r.err
}
