package fakes

import "reflect"

type FakeUnmarshaler struct {
	unmarshaledValue   any
	unmarshalCallCount int
	unmarshalReturns   struct {
		err error
	}
	unmarshalArgsForCall []struct {
		data []byte
		v    any
	}
}

func (um *FakeUnmarshaler) UnmarshalArgsForCall(i int) ([]byte, any) {
	args := um.unmarshalArgsForCall[i]
	return args.data, args.v
}

func (um *FakeUnmarshaler) UnmarshalCallCount() int {
	return um.unmarshalCallCount
}

func (um *FakeUnmarshaler) UnmarshalReturns(err error) {
	um.unmarshalReturns = struct{ err error }{err}
}

func (um *FakeUnmarshaler) UnmarshalUnmarshaledValue(v any) {
	um.unmarshaledValue = v
}

func (um *FakeUnmarshaler) Unmarshal(data []byte, v any) error {
	um.unmarshalCallCount++
	um.unmarshalArgsForCall = append(um.unmarshalArgsForCall, struct {
		data []byte
		v    any
	}{data, v})

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Pointer {
		if um.unmarshaledValue != nil {
			val.Elem().Set(reflect.ValueOf(um.unmarshaledValue))
		}
	}

	return um.unmarshalReturns.err
}
