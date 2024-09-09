package hasher_test

import (
	"fmt"

	"github.com/NathMcBride/digest-authentication/src/authentication/hasher"
	. "github.com/NathMcBride/digest-authentication/src/authentication/hasher/fakes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Hasher", func() {
	var (
		fakeCryptoHash    *FakeCryptoHash
		fakeSha256Factory *FakeCryptoFactory
		theHasher         *hasher.Hash
	)

	BeforeEach(func() {
		fakeCryptoHash = &FakeCryptoHash{}
		fakeCryptoHash.WriteReturns(1, nil)
		fakeCryptoHash.SumReturns([]byte("a-sum-result"))

		fakeSha256Factory = &FakeCryptoFactory{}
		fakeSha256Factory.NewReturns(fakeCryptoHash)

		theHasher = &hasher.Hash{
			CryptoFactory: fakeSha256Factory,
		}
	})

	It("creates a new crypto instance", func() {
		_, err := theHasher.Do("something-to-encrypt")

		Expect(err).NotTo(HaveOccurred())
		Expect(fakeSha256Factory.NewCallCount()).To(Equal(1))
	})

	Context("writing data to encrypt", func() {
		It("calls Write with the expected arguments", func() {
			theHasher.Do("something-to-encrypt")

			value := fakeCryptoHash.WriteArgsForCall(0)

			Expect(string(value[:])).To(Equal("something-to-encrypt"))
		})

		It("returns an error when call to Write fails", func() {
			fakeCryptoHash.WriteReturns(0, fmt.Errorf("an-error"))

			r, err := theHasher.Do("something-to-encrypt")

			Expect(r).To(BeEmpty())
			Expect(err).To(HaveOccurred())
		})
	})

	It("calls Sum and returns the result", func() {
		r, _ := theHasher.Do("something-to-encrypt")

		Expect(fakeCryptoHash.SumCallCount()).To(Equal(1))
		Expect(r).To(Equal(fmt.Sprintf("%x", []byte("a-sum-result"))))
	})
})
