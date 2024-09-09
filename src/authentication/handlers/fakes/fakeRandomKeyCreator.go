package fakes

type FakeRandomKeyCreator struct {
	callCount     int
	createReturns struct {
		key string
	}
}

func (f *FakeRandomKeyCreator) CreateCallCount() int {
	return f.callCount
}

func (f *FakeRandomKeyCreator) CreateReturns(key string) {
	f.createReturns = struct{ key string }{key: key}
}

func (f *FakeRandomKeyCreator) Create() string {
	f.callCount++
	return f.createReturns.key
}
