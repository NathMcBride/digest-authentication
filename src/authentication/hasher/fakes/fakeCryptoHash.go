package fakes

type FakeCryptoHash struct {
	writeCallCount int
	writeReturns   struct {
		n   int
		err error
	}
	writeArgsForCall []struct {
		p []byte
	}

	sumCallCount int
	sumReturns   struct {
		bytes []byte
	}
	sumArgsForCall []struct {
		b []byte
	}
}

func (f *FakeCryptoHash) WriteReturns(n int, err error) {
	f.writeReturns = struct {
		n   int
		err error
	}{n, err}
}

func (f *FakeCryptoHash) WriteArgsForCall(i int) []byte {
	args := f.writeArgsForCall[i]
	return args.p
}

func (f *FakeCryptoHash) WriteCallCount() int {
	return f.writeCallCount
}

func (f *FakeCryptoHash) Write(p []byte) (n int, err error) {
	f.writeCallCount++
	f.writeArgsForCall = append(f.writeArgsForCall, struct {
		p []byte
	}{p})

	return f.writeReturns.n, f.writeReturns.err
}

func (f *FakeCryptoHash) SumReturns(b []byte) {
	f.sumReturns = struct{ bytes []byte }{b}
}

func (f *FakeCryptoHash) SumArgsForCall(i int) []byte {
	args := f.sumArgsForCall[i]
	return args.b
}

func (f *FakeCryptoHash) SumCallCount() int {
	return f.sumCallCount
}

func (f *FakeCryptoHash) Sum(b []byte) []byte {
	f.sumCallCount++
	f.sumArgsForCall = append(f.sumArgsForCall, struct {
		b []byte
	}{b})

	return f.sumReturns.bytes
}

func (f *FakeCryptoHash) Reset()         {}
func (f *FakeCryptoHash) Size() int      { return 0 }
func (f *FakeCryptoHash) BlockSize() int { return 0 }
