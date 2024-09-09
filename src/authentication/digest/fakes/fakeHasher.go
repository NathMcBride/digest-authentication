package fakes

type FakeHasher struct {
	callCount        int
	argsForCall      []struct{ data string }
	doReturnsForCall map[int]struct {
		hash string
		err  error
	}
}

func (fh *FakeHasher) DoCallCount() int {
	return fh.callCount
}

func (fh *FakeHasher) DoArgsForCall(i int) string {
	return fh.argsForCall[i].data
}

func (fh *FakeHasher) DoReturnsOnCall(i int, hash string, err error) {
	if fh.doReturnsForCall == nil {
		fh.doReturnsForCall = make(map[int]struct {
			hash string
			err  error
		})
	}

	fh.doReturnsForCall[i] = struct {
		hash string
		err  error
	}{hash: hash, err: err}
}

func (fh *FakeHasher) Do(data string) (string, error) {
	fh.callCount++
	fh.argsForCall = append(fh.argsForCall, struct{ data string }{data})

	ret, hasReturn := fh.doReturnsForCall[fh.callCount-1]
	if hasReturn {
		return ret.hash, ret.err
	}

	return "", nil
}
