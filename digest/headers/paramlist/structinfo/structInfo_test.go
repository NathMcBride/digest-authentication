package structinfo_test

import (
	"reflect"

	"github.com/NathMcBride/web-authentication/digest/headers/paramlist/structinfo"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Struct information", func() {
	var (
		structInfo structinfo.StructInfo
	)

	BeforeEach(func() {
		structInfo = structinfo.StructInfo{}
	})

	It("gets type information", func() {
		val := struct {
			TestField1 string `httpparam:"testfield"`
			TestField2 string `httpparam:"testfield2"`
		}{}

		info := structInfo.GetTypeInfo(reflect.TypeOf(val))

		Expect(len(info.Fields)).To(Equal(2))

		f1 := info.Fields[0]
		Expect(f1.Idx[0]).To(Equal(0))
		Expect(f1.Name).To(Equal("testfield"))
		Expect(f1.Flags).To(Equal(structinfo.FElement))

		f2 := info.Fields[1]
		Expect(f2.Idx[0]).To(Equal(1))
		Expect(f2.Name).To(Equal("testfield2"))
		Expect(f2.Flags).To(Equal(structinfo.FElement))
	})

	Context("field name", func() {
		It("uses the field name provided in tag", func() {
			val := struct {
				TestField string `httpparam:"renamed"`
			}{}

			info := structInfo.GetTypeInfo(reflect.TypeOf(val))
			f1 := info.Fields[0]

			Expect(f1.Name).To(Equal("renamed"))
		})

		It("uses the field name if no new name provided", func() {
			val := struct {
				TestField1 string `httpparam`
				TestField2 string `httpparam:,unq`
			}{}

			info := structInfo.GetTypeInfo(reflect.TypeOf(val))

			f1 := info.Fields[0]
			Expect(f1.Name).To(Equal("TestField1"))

			f2 := info.Fields[1]
			Expect(f2.Name).To(Equal("TestField2"))
		})
	})

	Context("Excluding fields", func() {
		It("omits private fields", func() {
			val := struct {
				privateField string `httpparam:"privatefield"`
			}{}

			info := structInfo.GetTypeInfo(reflect.TypeOf(val))

			Expect(len(info.Fields)).To(Equal(0))
		})

		It("omits fields with hyphen parameter", func() {
			val := struct {
				PublicField string `httpparam:"-"`
			}{}

			info := structInfo.GetTypeInfo(reflect.TypeOf(val))

			Expect(len(info.Fields)).To(Equal(0))
		})
	})

	Context("Field flags", func() {
		It("sets the unquoted flag FUnq", func() {
			val := struct {
				PublicField string `httpparam:"publicfield,unq"`
			}{}

			info := structInfo.GetTypeInfo(reflect.TypeOf(val))

			finfo := &info.Fields[0]
			isFUnqSet := finfo.Flags&structinfo.FUnq != 0
			Expect(isFUnqSet).To(BeTrue(), "unq flag not set")
		})

		It("sets the omit empty flag FOmitEmpty", func() {
			val := struct {
				PublicField string `httpparam:"publicfield,omitempty"`
			}{}

			info := structInfo.GetTypeInfo(reflect.TypeOf(val))

			finfo := &info.Fields[0]
			isOmitEmptySet := finfo.Flags&structinfo.FOmitEmpty != 0
			Expect(isOmitEmptySet).To(BeTrue(), "omit empty flag not set")
		})

		It("sets multiple flags", func() {
			val := struct {
				PublicField string `httpparam:"publicfield,unq,omitempty"`
			}{}

			info := structInfo.GetTypeInfo(reflect.TypeOf(val))

			finfo := &info.Fields[0]
			isOmitEmptySet := finfo.Flags&structinfo.FOmitEmpty != 0
			Expect(isOmitEmptySet).To(BeTrue(), "omit empty flag not set")

			isFUnqSet := finfo.Flags&structinfo.FUnq != 0
			Expect(isFUnqSet).To(BeTrue(), "unq flag not set")
		})
	})

})
