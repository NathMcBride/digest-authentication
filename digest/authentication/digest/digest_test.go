package digest_test

import (
	"fmt"

	"github.com/NathMcBride/web-authentication/digest/authentication/digest"
	"github.com/NathMcBride/web-authentication/digest/authentication/model"
	"github.com/NathMcBride/web-authentication/digest/providers/credential"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type FakeHasher struct {
	callCount          int
	argsForCall        []struct{ data string }
	hashReturnsForCall map[int]struct {
		hash string
		err  error
	}
}

func (fh *FakeHasher) HashCallCount() int {
	return fh.callCount
}

func (fh *FakeHasher) HashArgsForCall(i int) string {
	return fh.argsForCall[i].data
}

func (fh *FakeHasher) HashReturnsOnCall(i int, hash string, err error) {
	if fh.hashReturnsForCall == nil {
		fh.hashReturnsForCall = make(map[int]struct {
			hash string
			err  error
		})
	}

	fh.hashReturnsForCall[i] = struct {
		hash string
		err  error
	}{hash: hash, err: err}
}

func (fh *FakeHasher) Hash(data string) (string, error) {
	fh.callCount++
	fh.argsForCall = append(fh.argsForCall, struct{ data string }{data})

	ret, hasReturn := fh.hashReturnsForCall[fh.callCount-1]
	if hasReturn {
		return ret.hash, ret.err
	}

	return "", nil
}

var _ = Describe("Digest calculation RFC7616", func() {
	var (
		fakeHasher *FakeHasher
		theDigest  *digest.Digest
	)

	BeforeEach(func() {
		fakeHasher = &FakeHasher{}
		fakeHasher.HashReturnsOnCall(0, "a-ha1-hash", nil)
		fakeHasher.HashReturnsOnCall(1, "a-ha2-hash", nil)
		fakeHasher.HashReturnsOnCall(2, "a-kd-hash", nil)

		theDigest = &digest.Digest{Sha256: fakeHasher}
	})

	Context("on success", func() {
		It("returns the hash result", func() {
			hash, err := theDigest.Calculate(credential.Credentials{}, model.AuthHeader{}, "")

			Expect(err).NotTo(HaveOccurred())
			Expect(hash).To(Equal("a-kd-hash"))
		})
	})

	Context("hashing HA1", func() {
		It("is called with expected arguments", func() {
			credentials := credential.Credentials{Username: "a-username", Password: "a-password"}
			authHeader := model.AuthHeader{Realm: "a-realm"}

			theDigest.Calculate(credentials, authHeader, "")

			args := fakeHasher.HashArgsForCall(0)
			Expect(args).To(Equal("a-username:a-realm:a-password"))
		})

		It("returns an error", func() {
			fakeHasher.HashReturnsOnCall(0, "", fmt.Errorf("hashing failed"))

			_, err := theDigest.Calculate(credential.Credentials{}, model.AuthHeader{}, "")

			Expect(err).To(HaveOccurred())
		})
	})

	Context("hashing HA2", func() {
		It("is called with the expected arguments", func() {
			credentials := credential.Credentials{}
			authHeader := model.AuthHeader{Uri: "a-uri"}

			theDigest.Calculate(credentials, authHeader, "a-http-method")

			args := fakeHasher.HashArgsForCall(1)
			Expect(args).To(Equal("a-http-method:a-uri"))
		})

		It("returns an error", func() {
			fakeHasher.HashReturnsOnCall(1, "", fmt.Errorf("hashing failed"))

			_, err := theDigest.Calculate(credential.Credentials{}, model.AuthHeader{}, "")

			Expect(err).To(HaveOccurred())
		})
	})

	Context("hashing KD", func() {
		It("is called with the expected arguments", func() {
			credentials := credential.Credentials{}
			authHeader := model.AuthHeader{
				Nonce:  "a-nonce-value",
				Nc:     "a-nonce-count",
				Cnonce: "a-client-nonce",
				Qop:    "a-qop"}

			theDigest.Calculate(credentials, authHeader, "")

			args := fakeHasher.HashArgsForCall(2)
			expected := "a-ha1-hash:a-nonce-value:a-nonce-count:a-client-nonce:a-qop:a-ha2-hash"
			Expect(args).To(Equal(expected))
		})

		It("returns an error", func() {
			fakeHasher.HashReturnsOnCall(2, "", fmt.Errorf("hashing failed"))

			_, err := theDigest.Calculate(credential.Credentials{}, model.AuthHeader{}, "")

			Expect(err).To(HaveOccurred())
		})
	})
})
