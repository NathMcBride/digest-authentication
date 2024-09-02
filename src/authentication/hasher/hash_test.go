package hasher_test

import (
	"fmt"

	"github.com/NathMcBride/digest-authentication/src/authentication/hasher"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type FakeCryptoHash struct {
	writeCallCount int
	writeReturns   struct {
		n   int
		err error
	}
	writeArgsForCall []struct {
		p []byte
	}

	sumCallCount int
	sumReturns   struct {
		bytes []byte
	}
	sumArgsForCall []struct {
		b []byte
	}
}

func (f *FakeCryptoHash) WriteReturns(n int, err error) {
	f.writeReturns = struct {
		n   int
		err error
	}{n, err}
}

func (f *FakeCryptoHash) WriteArgsForCall(i int) []byte {
	args := f.writeArgsForCall[i]
	return args.p
}

func (f *FakeCryptoHash) WriteCallCount() int {
	return f.writeCallCount
}

func (f *FakeCryptoHash) Write(p []byte) (n int, err error) {
	f.writeCallCount++
	f.writeArgsForCall = append(f.writeArgsForCall, struct {
		p []byte
	}{p})

	return f.writeReturns.n, f.writeReturns.err
}

func (f *FakeCryptoHash) SumReturns(b []byte) {
	f.sumReturns = struct{ bytes []byte }{b}
}

func (f *FakeCryptoHash) SumArgsForCall(i int) []byte {
	args := f.sumArgsForCall[i]
	return args.b
}

func (f *FakeCryptoHash) SumCallCount() int {
	return f.sumCallCount
}

func (f *FakeCryptoHash) Sum(b []byte) []byte {
	f.sumCallCount++
	f.sumArgsForCall = append(f.sumArgsForCall, struct {
		b []byte
	}{b})

	return f.sumReturns.bytes
}

func (f *FakeCryptoHash) Reset()         {}
func (f *FakeCryptoHash) Size() int      { return 0 }
func (f *FakeCryptoHash) BlockSize() int { return 0 }

type FakeCryptoFactory struct {
	newCallCount int
	newReturns   struct {
		hash hasher.CryptoHash
	}
}

func (f *FakeCryptoFactory) NewReturns(hash hasher.CryptoHash) {
	f.newReturns = struct{ hash hasher.CryptoHash }{hash}
}

func (f *FakeCryptoFactory) NewCallCount() int {
	return f.newCallCount
}

func (f *FakeCryptoFactory) New() hasher.CryptoHash {
	f.newCallCount++
	return f.newReturns.hash
}

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
