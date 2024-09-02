package headers_test

import (
	"fmt"

	"github.com/NathMcBride/digest-authentication/src/authentication/model"
	"github.com/NathMcBride/digest-authentication/src/constants"
	"github.com/NathMcBride/digest-authentication/src/headers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

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

var _ = Describe("Challenge", func() {
	var (
		fakeParamListMarshaler *FakeParamListMarshaler
		digestChallenge        headers.DigestChallenge
	)

	BeforeEach(func() {
		fakeParamListMarshaler = &FakeParamListMarshaler{}
		fakeParamListMarshaler.MarshalReturnsOnCall([]byte("a-marshaled-struct"), nil)
		digestChallenge = headers.DigestChallenge{Marshaler: fakeParamListMarshaler}
	})

	It("creates a digest challenge", func() {
		result, _ := digestChallenge.Create("a-realm", "an-opaque", "a-nonce-value", true)

		Expect(fakeParamListMarshaler.MarshalCallCount()).To(Equal(1))
		Expect(result).To(Equal("Digest a-marshaled-struct"))
	})

	Context("call to Marshal", func() {
		It("is called with expected arguments", func() {
			expected := model.DigestHeader{
				Realm:     "a-realm",
				Algorithm: constants.SHA256,
				Qop:       constants.Auth,
				Opaque:    "an-opaque",
				Nonce:     "a-nonce-value",
				UserHash:  true,
			}

			digestChallenge.Create("a-realm", "an-opaque", "a-nonce-value", true)

			Expect(fakeParamListMarshaler.MarshalArgsForCall(0)).To(Equal(expected))
		})

		It("returns an error on failure", func() {
			fakeParamListMarshaler.MarshalReturnsOnCall(nil, fmt.Errorf("an-error"))

			_, err := digestChallenge.Create("a-realm", "an-opaque", "a-nonce-value", true)

			Expect(err).To(HaveOccurred())
		})
	})

})
