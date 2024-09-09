package fakes

type FakeParamListMarshaler struct {
	marshalCallCount   int
	marshalArgsForCall []struct{ v any }
	marshalReturns     struct {
		b   []byte
		err error
	}
}

func (f *FakeParamListMarshaler) MarshalCallCount() int {
	return f.marshalCallCount
}

func (f *FakeParamListMarshaler) MarshalArgsForCall(i int) any {
	return f.marshalArgsForCall[i].v
}

func (f *FakeParamListMarshaler) MarshalReturnsOnCall(b []byte, err error) {
	f.marshalReturns = struct {
		b   []byte
		err error
	}{b, err}
}

func (f *FakeParamListMarshaler) Marshal(v any) ([]byte, error) {
	f.marshalCallCount++
	f.marshalArgsForCall = append(f.marshalArgsForCall, struct{ v any }{v})

	return f.marshalReturns.b, f.marshalReturns.err
}
