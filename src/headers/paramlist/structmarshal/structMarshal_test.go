package structmarshal_test

import (
	"bytes"
	"reflect"

	"github.com/NathMcBride/digest-authentication/src/headers/paramlist/structinfo"
	"github.com/NathMcBride/digest-authentication/src/headers/paramlist/structmarshal"
	. "github.com/NathMcBride/digest-authentication/src/headers/paramlist/support"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Marshal a struct", func() {
	var (
		structMarshal structmarshal.StructMarshal
		toMarshal     TestStruct
		buffer        *bytes.Buffer
	)

	BeforeEach(func() {
		buffer = new(bytes.Buffer)
		structMarshal = structmarshal.StructMarshal{}
		toMarshal = TestStruct{"some-value"}
	})

	It("can successfully marshal a struct", func() {
		info := NewMakeStructInfo().WithNoFlags().Build()

		err := structMarshal.Marshal(buffer, &info, reflect.ValueOf(toMarshal))

		Expect(err).NotTo(HaveOccurred())
		Expect(string(buffer.Bytes()[:])).To(Equal(`field="some-value"`))
	})

	When("info argument contains no fields", func() {
		It("returns nil & the buffer is empty", func() {
			info := NewMakeStructInfo().Build()

			err := structMarshal.Marshal(buffer, &info, reflect.ValueOf(toMarshal))

			Expect(err).NotTo(HaveOccurred())
			Expect(string(buffer.Bytes()[:])).To(BeEmpty())
		})
	})

	DescribeTable("Processing flags",
		func(i structinfo.Info, t TestStruct, expected string) {
			b := new(bytes.Buffer)
			err := structMarshal.Marshal(b, &i, reflect.ValueOf(t))

			Expect(err).NotTo(HaveOccurred())
			Expect(string(b.Bytes()[:])).To(Equal(expected))
		},
		Entry("FUnq",
			NewMakeStructInfo().WithFUnqFlag().Build(),
			TestStruct{"some-value"},
			`field=some-value`),
		Entry("FOmitEmpty",
			NewMakeStructInfo().WithFOmitEmptyFlag().Build(),
			TestStruct{""},
			""),
		Entry("FUnq | FOmitEmpty - with empty field",
			NewMakeStructInfo().WithAllFlags().Build(),
			TestStruct{""},
			""),
		Entry("FUnq | FOmitEmpty - with value in field",
			NewMakeStructInfo().WithAllFlags().Build(),
			TestStruct{"some-value"},
			"field=some-value"),
	)
})
