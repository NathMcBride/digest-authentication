package fakes

import (
	"bytes"
	"reflect"

	"github.com/NathMcBride/digest-authentication/src/headers/paramlist/structinfo"
)

type FakeStructMarshaler struct {
	marshalCallCount   int
	marshalReturns     struct{ err error }
	marshalWrites      string
	marshalArgsForCall []struct {
		buffer *bytes.Buffer
		info   *structinfo.Info
		val    reflect.Value
	}
}

func (fs *FakeStructMarshaler) MarshalInfoCallCount() int {
	return fs.marshalCallCount
}

func (fs *FakeStructMarshaler) MarshalReturns(err error) {
	fs.marshalReturns = struct{ err error }{err}
}

func (fs *FakeStructMarshaler) MarshalWrites(s string) {
	fs.marshalWrites = s
}

func (fs *FakeStructMarshaler) MarshalArgsForCall(i int) (*bytes.Buffer, *structinfo.Info, reflect.Value) {
	args := fs.marshalArgsForCall[0]
	return args.buffer, args.info, args.val
}

func (fs *FakeStructMarshaler) Marshal(buffer *bytes.Buffer, info *structinfo.Info, val reflect.Value) error {
	fs.marshalCallCount++
	fs.marshalArgsForCall = append(fs.marshalArgsForCall, struct {
		buffer *bytes.Buffer
		info   *structinfo.Info
		val    reflect.Value
	}{buffer, info, val})

	buffer.WriteString(fs.marshalWrites)

	return fs.marshalReturns.err
}
