package fakes

type FakeClientStore struct {
	addCallCount   int
	addArgsForCall []struct {
		entry string
	}
}

func (f *FakeClientStore) AddCallCount() int {
	return f.addCallCount
}

func (f *FakeClientStore) Add(entry string) {
	f.addCallCount++
	f.addArgsForCall = append(f.addArgsForCall, struct {
		entry string
	}{
		entry: entry})
}

func (f *FakeClientStore) Has(entry string) bool {
	return true
}

func (f *FakeClientStore) Delete(entry string) {}
