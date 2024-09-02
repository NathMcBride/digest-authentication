package digest_test

import (
	"fmt"

	"github.com/NathMcBride/digest-authentication/src/authentication/digest"
	"github.com/NathMcBride/digest-authentication/src/authentication/model"
	"github.com/NathMcBride/digest-authentication/src/providers/credential"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type FakeHasher struct {
	callCount        int
	argsForCall      []struct{ data string }
	doReturnsForCall map[int]struct {
		hash string
		err  error
	}
}

func (fh *FakeHasher) DoCallCount() int {
	return fh.callCount
}

func (fh *FakeHasher) DoArgsForCall(i int) string {
	return fh.argsForCall[i].data
}

func (fh *FakeHasher) DoReturnsOnCall(i int, hash string, err error) {
	if fh.doReturnsForCall == nil {
		fh.doReturnsForCall = make(map[int]struct {
			hash string
			err  error
		})
	}

	fh.doReturnsForCall[i] = struct {
		hash string
		err  error
	}{hash: hash, err: err}
}

func (fh *FakeHasher) Do(data string) (string, error) {
	fh.callCount++
	fh.argsForCall = append(fh.argsForCall, struct{ data string }{data})

	ret, hasReturn := fh.doReturnsForCall[fh.callCount-1]
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
		fakeHasher.DoReturnsOnCall(0, "a-ha1-hash", nil)
		fakeHasher.DoReturnsOnCall(1, "a-ha2-hash", nil)
		fakeHasher.DoReturnsOnCall(2, "a-kd-hash", nil)

		theDigest = &digest.Digest{Hasher: fakeHasher}
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

			args := fakeHasher.DoArgsForCall(0)
			Expect(args).To(Equal("a-username:a-realm:a-password"))
		})

		It("returns an error", func() {
			fakeHasher.DoReturnsOnCall(0, "", fmt.Errorf("hashing failed"))

			_, err := theDigest.Calculate(credential.Credentials{}, model.AuthHeader{}, "")

			Expect(err).To(HaveOccurred())
		})
	})

	Context("hashing HA2", func() {
		It("is called with the expected arguments", func() {
			credentials := credential.Credentials{}
			authHeader := model.AuthHeader{Uri: "a-uri"}

			theDigest.Calculate(credentials, authHeader, "a-http-method")

			args := fakeHasher.DoArgsForCall(1)
			Expect(args).To(Equal("a-http-method:a-uri"))
		})

		It("returns an error", func() {
			fakeHasher.DoReturnsOnCall(1, "", fmt.Errorf("hashing failed"))

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

			args := fakeHasher.DoArgsForCall(2)
			expected := "a-ha1-hash:a-nonce-value:a-nonce-count:a-client-nonce:a-qop:a-ha2-hash"
			Expect(args).To(Equal(expected))
		})

		It("returns an error", func() {
			fakeHasher.DoReturnsOnCall(2, "", fmt.Errorf("hashing failed"))

			_, err := theDigest.Calculate(credential.Credentials{}, model.AuthHeader{}, "")

			Expect(err).To(HaveOccurred())
		})
	})
})
