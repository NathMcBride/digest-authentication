package fakes

type FakeChallengeCreator struct {
	callCount         int
	createArgsForCall []struct {
		realm              string
		opaque             string
		nonce              string
		shouldHashUserName bool
	}
	createReturns struct {
		header string
		err    error
	}
}

func (f *FakeChallengeCreator) CreateCallCount() int {
	return f.callCount
}

func (f *FakeChallengeCreator) CreateArgsForCall(i int) (realm string, opaque string, nonce string, shouldHashUserName bool) {
	args := f.createArgsForCall[i]
	return args.realm, args.opaque, args.nonce, args.shouldHashUserName
}

func (f *FakeChallengeCreator) CreateReturns(header string, err error) {
	f.createReturns = struct {
		header string
		err    error
	}{header, err}
}

func (f *FakeChallengeCreator) Create(realm string, opaque string, nonce string, shouldHashUserName bool) (string, error) {
	f.callCount++
	f.createArgsForCall = append(f.createArgsForCall, struct {
		realm              string
		opaque             string
		nonce              string
		shouldHashUserName bool
	}{
		realm,
		opaque,
		nonce,
		shouldHashUserName})

	return f.createReturns.header, f.createReturns.err
}
