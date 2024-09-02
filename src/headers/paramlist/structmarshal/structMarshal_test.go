package structmarshal_test

import (
	"bytes"
	"reflect"

	"github.com/NathMcBride/digest-authentication/src/headers/paramlist/structinfo"
	"github.com/NathMcBride/digest-authentication/src/headers/paramlist/structmarshal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Marshal a struct", func() {
	var (
		structMarshal structmarshal.StructMarshal
		toMarshal     struct{ Field1 string }
		buffer        *bytes.Buffer
	)

	BeforeEach(func() {
		buffer = new(bytes.Buffer)
		structMarshal = structmarshal.StructMarshal{}
		toMarshal = struct {
			Field1 string
		}{"some-value"}
	})

	It("can successfully marshal a struct", func() {
		info := structinfo.Info{
			Fields: []structinfo.FieldInfo{
				{Idx: []int{0},
					Name:  "field1",
					Flags: 0},
			},
		}

		err := structMarshal.Marshal(buffer, &info, reflect.ValueOf(toMarshal))

		Expect(err).NotTo(HaveOccurred())
		Expect(string(buffer.Bytes()[:])).To(Equal(`field1="some-value"`))
	})

	When("info argument contains no fields", func() {
		It("returns nil & the buffer is empty", func() {
			info := structinfo.Info{Fields: []structinfo.FieldInfo{}}

			toMarshal := struct {
				Field1 string
			}{"some-value"}

			err := structMarshal.Marshal(buffer, &info, reflect.ValueOf(toMarshal))

			Expect(err).NotTo(HaveOccurred())
			Expect(string(buffer.Bytes()[:])).To(BeEmpty())
		})
	})

	DescribeTable("Processing flags",
		func(i structinfo.Info, t struct {
			Field string
		}, expected string) {
			b := new(bytes.Buffer)
			err := structMarshal.Marshal(b, &i, reflect.ValueOf(t))

			Expect(err).NotTo(HaveOccurred())
			Expect(string(b.Bytes()[:])).To(Equal(expected))
		},
		Entry("FUnq",
			structinfo.Info{
				Fields: []structinfo.FieldInfo{
					{Idx: []int{0},
						Name:  "field",
						Flags: structinfo.FUnq},
				}},
			struct {
				Field string
			}{"some-value"}, `field=some-value`),
		Entry("FOmitEmpty",
			structinfo.Info{
				Fields: []structinfo.FieldInfo{
					{Idx: []int{0},
						Name:  "field",
						Flags: structinfo.FOmitEmpty},
				}},
			struct {
				Field string
			}{""}, ""),
		Entry("FUnq | FOmitEmpty - with empty field",
			structinfo.Info{
				Fields: []structinfo.FieldInfo{
					{Idx: []int{0},
						Name:  "field",
						Flags: structinfo.FUnq | structinfo.FOmitEmpty},
				}},
			struct {
				Field string
			}{""}, ""),
		Entry("FUnq | FOmitEmpty - with value in field",
			structinfo.Info{
				Fields: []structinfo.FieldInfo{
					{Idx: []int{0},
						Name:  "field",
						Flags: structinfo.FUnq | structinfo.FOmitEmpty},
				}},
			struct {
				Field string
			}{"some-value"}, "field=some-value"),
	)
})
