package paramlist_test

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/NathMcBride/digest-authentication/src/headers/paramlist"
	. "github.com/NathMcBride/digest-authentication/src/headers/paramlist/fakes"
	"github.com/NathMcBride/digest-authentication/src/headers/paramlist/structinfo"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

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
		fakeStructMarshaler.MarshalWrites("some-marshaled-value")

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
