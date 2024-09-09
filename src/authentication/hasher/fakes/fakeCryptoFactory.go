package fakes

import "github.com/NathMcBride/digest-authentication/src/authentication/hasher"

type FakeCryptoFactory struct {
	newCallCount int
	newReturns   struct {
		hash hasher.CryptoHash
	}
}

func (f *FakeCryptoFactory) NewReturns(hash hasher.CryptoHash) {
	f.newReturns = struct{ hash hasher.CryptoHash }{hash}
}

func (f *FakeCryptoFactory) NewCallCount() int {
	return f.newCallCount
}

func (f *FakeCryptoFactory) New() hasher.CryptoHash {
	f.newCallCount++
	return f.newReturns.hash
}
