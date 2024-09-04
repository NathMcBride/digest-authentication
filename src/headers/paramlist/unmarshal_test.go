package paramlist_test

import (
	"fmt"
	"reflect"

	"github.com/NathMcBride/digest-authentication/src/headers/paramlist"
	"github.com/NathMcBride/digest-authentication/src/headers/paramlist/structinfo"
	. "github.com/NathMcBride/digest-authentication/src/headers/paramlist/support"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type FakeParser struct {
	callCount    int
	parseReturns struct {
		parsed map[string]string
		err    error
	}
	parseArgsForCall []struct{ auth string }
}

func (p *FakeParser) ParseCallCount() int {
	return p.callCount
}

func (p *FakeParser) ParseReturns(parsed map[string]string, err error) {
	p.parseReturns = struct {
		parsed map[string]string
		err    error
	}{
		parsed,
		err,
	}
}

func (p *FakeParser) ParseArgsForCall(i int) string {
	args := p.parseArgsForCall[i]
	return args.auth
}

func (p *FakeParser) Parse(auth string) (map[string]string, error) {
	p.callCount++
	p.parseArgsForCall = append(p.parseArgsForCall, struct{ auth string }{auth})
	return p.parseReturns.parsed, p.parseReturns.err
}

var _ = Describe("Unmarshal", func() {
	var (
		into struct {
			Field  string `httpparam:"field"`
			Field2 bool   `httpparam:"field2"`
		}
		fakeParser       *FakeParser
		fakeStructInfoer *FakeStructInfoer
		unmarshaler      paramlist.UnMarshaler
	)

	BeforeEach(func() {
		into = struct {
			Field  string `httpparam:"field"`
			Field2 bool   `httpparam:"field2"`
		}{}

		fakeParser = &FakeParser{}
		fakeParser.ParseReturns(
			map[string]string{
				"field":  "a-parsed-value",
				"field2": "true",
			},
			nil)

		fakeStructInfoer = &FakeStructInfoer{}
		fakeStructInfoer.GetTypeInfoReturns(
			NewMakeStructInfo().
				WithNoFlags().
				AddField("field2", 0).
				Build())

		unmarshaler = paramlist.UnMarshaler{
			StructInfoer: fakeStructInfoer,
			Parser:       fakeParser,
		}
	})

	It("successfully unmarshal's data to a struct", func() {
		err := unmarshaler.Unmarshal([]byte{}, &into)

		Expect(err).NotTo(HaveOccurred())
		Expect(into.Field).To(Equal("a-parsed-value"))
		Expect(into.Field2).To(BeTrue())
	})

	Context("argument validation", func() {
		When("arg v is not a pointer", func() {
			It("returns an error", func() {
				err := unmarshaler.Unmarshal([]byte{}, struct{}{})

				Expect(err).To(HaveOccurred())
			})
		})

		When("arg v is a nil pointer", func() {
			It("returns an error", func() {
				into := &struct{}{}
				into = nil

				err := unmarshaler.Unmarshal([]byte{}, into)

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("calling Parse", func() {
		It("is called with the expected arguments", func() {
			fakeParser.ParseReturns(map[string]string{}, nil)
			fakeStructInfoer.GetTypeInfoReturns(structinfo.Info{})

			unmarshaler.Unmarshal([]byte("some data"), &struct{}{})

			Expect(fakeParser.ParseCallCount()).To(Equal(1))
			data := fakeParser.ParseArgsForCall(0)
			Expect(data).To(Equal("some data"))
		})

		Context("on failure", func() {
			It("returns an error", func() {
				fakeParser.ParseReturns(map[string]string{}, fmt.Errorf("an-error"))

				err := unmarshaler.Unmarshal([]byte("some data"), &struct{}{})

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("calling GetTypeInfo", func() {
		It("is called with the expected arguments", func() {
			unmarshaler.Unmarshal([]byte{}, &into)

			Expect(fakeStructInfoer.GetTypeInfoCallCount()).To(Equal(1))
			expected := reflect.TypeOf(into)
			Expect(fakeStructInfoer.GetTypeInfoArgsForCall(0)).To(Equal(expected))
		})

		Context("on empty type info", func() {
			It("returns nil", func() {
				fakeStructInfoer.GetTypeInfoReturns(structinfo.Info{})

				err := unmarshaler.Unmarshal([]byte{}, &struct{}{})

				Expect(err).To(BeNil())
			})
		})
	})
})
