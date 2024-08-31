package paramlist_test

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/NathMcBride/web-authentication/digest/headers/paramlist"
	"github.com/NathMcBride/web-authentication/digest/headers/paramlist/structinfo"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type FakeStructInfoer struct {
	getTypeInfoCallCount   int
	getTypeInfoReturns     struct{ info structinfo.Info }
	getTypeInfoArgsForCall []struct{ typ reflect.Type }
}

func (fs *FakeStructInfoer) GetTypeInfoCallCount() int {
	return fs.getTypeInfoCallCount
}

func (fs *FakeStructInfoer) GetTypeInfoReturns(info structinfo.Info) {
	fs.getTypeInfoReturns = struct{ info structinfo.Info }{info}
}

func (fs *FakeStructInfoer) GetTypeInfoArgsForCall(i int) reflect.Type {
	return fs.getTypeInfoArgsForCall[0].typ
}

func (fs *FakeStructInfoer) GetTypeInfo(typ reflect.Type) *structinfo.Info {
	fs.getTypeInfoCallCount++
	fs.getTypeInfoArgsForCall = append(fs.getTypeInfoArgsForCall, struct{ typ reflect.Type }{typ})
	return &fs.getTypeInfoReturns.info
}

func (fs *FakeStructInfoer) FieldInfo(f *reflect.StructField) *structinfo.FieldInfo { return nil }

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

var _ = Describe("Marshal", func() {
	var (
		fakeStructInfoer    *FakeStructInfoer
		fakeStructMarshaler *FakeStructMarshaler
		marshaler           paramlist.Marshaler
	)

	BeforeEach(func() {
		fakeStructInfoer = &FakeStructInfoer{}
		fakeStructInfoer.GetTypeInfoReturns(
			structinfo.Info{
				Fields: []structinfo.FieldInfo{
					{Name: "a-field-name"},
				},
			})

		fakeStructMarshaler = &FakeStructMarshaler{}
		fakeStructMarshaler.marshalWrites = "some-marshaled-value"

		marshaler = paramlist.Marshaler{
			StructInfoer:    fakeStructInfoer,
			StructMarshaler: fakeStructMarshaler}
	})

	It("successfully marshals a struct", func() {
		result, err := marshaler.Marshal(struct{}{})

		Expect(err).NotTo(HaveOccurred())
		Expect(fakeStructInfoer.GetTypeInfoCallCount()).To(Equal(1))
		Expect(fakeStructMarshaler.MarshalInfoCallCount()).To(Equal(1))
		Expect(string(result[:])).To(Equal("some-marshaled-value"))
	})

	It("errors upon receiving an invalid argument", func() {
		_, err := marshaler.Marshal(nil)
		Expect(err).To(HaveOccurred())
	})

	It("passes the correct arguments to GetTypeInfo", func() {
		arg := struct{ something string }{"something-value"}

		marshaler.Marshal(arg)

		expected := reflect.ValueOf(arg).Type()
		Expect(fakeStructInfoer.GetTypeInfoArgsForCall(0)).To(Equal(expected))
	})

	Context("calling Marshal", func() {
		It("passes the correct arguments", func() {
			fakeStructMarshaler.MarshalWrites("")

			arg := struct{ something string }{"something-value"}
			marshaler.Marshal(arg)

			buffer, tinfo, val := fakeStructMarshaler.MarshalArgsForCall(0)
			Expect(buffer).To(Equal(new(bytes.Buffer)))
			Expect(tinfo).To(Equal(tinfo))
			Expect(val.Equal(reflect.ValueOf(arg))).To(BeTrue())
		})

		It("returns an error on failure", func() {
			fakeStructMarshaler.MarshalReturns(fmt.Errorf("an-error"))

			_, err := marshaler.Marshal(struct{}{})

			Expect(err).To(HaveOccurred())
		})
	})
})
