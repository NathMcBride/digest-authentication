package authenticator_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/NathMcBride/digest-authentication/src/authentication/authenticator"
	"github.com/NathMcBride/digest-authentication/src/authentication/model"
	"github.com/NathMcBride/digest-authentication/src/providers/credential"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type FakeCredentialProvider struct {
	callCount   int
	argsForCall []struct {
		userID  string
		useHash bool
	}
	getCredentialsReturns struct {
		credentials *credential.Credentials
		found       bool
		err         error
	}
}

func (fcp *FakeCredentialProvider) GetCredentialsReturns(credentials *credential.Credentials, found bool, err error) {
	fcp.getCredentialsReturns.credentials = credentials
	fcp.getCredentialsReturns.found = found
	fcp.getCredentialsReturns.err = err
}

func (fcp *FakeCredentialProvider) GetCredentialsArgsForCall(i int) (string, bool) {
	args := fcp.argsForCall[i]
	return args.userID, args.useHash
}

func (fcp *FakeCredentialProvider) GetCredentials(userID string, useHash bool) (*credential.Credentials, bool, error) {
	fcp.callCount++

	fcp.argsForCall = append(fcp.argsForCall,
		struct {
			userID  string
			useHash bool
		}{
			userID:  userID,
			useHash: useHash})

	r := fcp.getCredentialsReturns
	return r.credentials, r.found, r.err
}

type FakeDigest struct {
	callCount   int
	argsForCall []struct {
		credentials credential.Credentials
		authHeader  model.AuthHeader
		Method      string
	}
	calculateReturns struct {
		digest string
		err    error
	}
}

func (fd *FakeDigest) CalculateReturns(digest string, err error) {
	fd.calculateReturns.digest = digest
	fd.calculateReturns.err = err
}

func (fd *FakeDigest) CalculateArgsForCall(i int) (credential.Credentials, model.AuthHeader, string) {
	args := fd.argsForCall[i]
	return args.credentials, args.authHeader, args.Method
}

func (fd *FakeDigest) Calculate(credentials credential.Credentials, authHeader model.AuthHeader, Method string) (string, error) {
	fd.callCount++

	fd.argsForCall = append(fd.argsForCall,
		struct {
			credentials credential.Credentials
			authHeader  model.AuthHeader
			Method      string
		}{
			credentials: credentials,
			authHeader:  authHeader,
			Method:      Method})

	r := fd.calculateReturns
	return r.digest, r.err
}

var SuccessAuthorizationHeader = `Digest response="a-digest-response",
username="a-username",
realm="a-realm",
algorithm=SHA-256,
qop=auth,
cnonce="a-client-nonce",
nc=6,
opaque="an-opaque-value",
uri="a-uri",
nonce="a-nonce-value",
userhash=true`

var _ = Describe("Authenticator", func() {

	Describe("Authenticate", func() {
		var (
			fakeCredentialsProvider *FakeCredentialProvider
			fakeDigest              *FakeDigest
			request                 *http.Request
			theAuthenticator        authenticator.Authenticator
		)
		BeforeEach(func() {
			fakeCredentialsProvider = &FakeCredentialProvider{}
			fakeCredentialsProvider.GetCredentialsReturns(
				&credential.Credentials{
					Username: "a-plain-username",
					Password: "a-password",
				},
				true,
				nil)

			fakeDigest = &FakeDigest{}
			fakeDigest.CalculateReturns("a-digest-response", nil)
			theAuthenticator = authenticator.Authenticator{
				Opaque:             "an-opaque-value",
				HashUserName:       true,
				CredentialProvider: fakeCredentialsProvider,
				Digest:             fakeDigest,
			}

			request = httptest.NewRequest("GET", "http://valid", nil)
		})

		It("successfully authenticates a request", func() {
			request.Header.Add("Authorization", SuccessAuthorizationHeader)

			session, err := theAuthenticator.Authenticate(request)

			Expect(err).NotTo(HaveOccurred())
			Expect(session.IsAuthenticated).To(BeTrue())
			Expect(session.User.UserID).To(Equal("a-username"))
		})

		When("the authorization header is missing", func() {
			It("returns an unauthorized session", func() {
				session, err := theAuthenticator.Authenticate(request)

				Expect(err).NotTo(HaveOccurred())
				Expect(session.IsAuthenticated).To(BeFalse())
			})
		})

		When("the authorization header is malformed", func() {
			It("returns an unauthorized session", func() {
				request.Header.Add("Authorization", "cheese")

				session, err := theAuthenticator.Authenticate(request)

				Expect(err).NotTo(HaveOccurred())
				Expect(session.IsAuthenticated).To(BeFalse())
			})
		})

		Describe("Call to GetCredentials", func() {
			Context("called with", func() {
				It("was called correctly", func() {
					request.Header.Add("Authorization", SuccessAuthorizationHeader)

					theAuthenticator.Authenticate(request)

					userId, shouldHash := fakeCredentialsProvider.GetCredentialsArgsForCall(0)
					Expect(userId).To(Equal("a-username"))
					Expect(shouldHash).To(BeTrue())
				})
			})

			Context("on error", func() {
				It("returns an error and an unauthorized session", func() {
					request.Header.Add("Authorization", SuccessAuthorizationHeader)
					fakeCredentialsProvider.GetCredentialsReturns(nil, false, fmt.Errorf("an-error-occurred"))

					session, err := theAuthenticator.Authenticate(request)

					Expect(err).To(HaveOccurred())
					Expect(session.IsAuthenticated).To(BeFalse())
				})
			})

			Context("on not found", func() {
				It("returns an unauthorized session", func() {
					request.Header.Add("Authorization", SuccessAuthorizationHeader)
					fakeCredentialsProvider.GetCredentialsReturns(nil, false, nil)

					session, err := theAuthenticator.Authenticate(request)

					Expect(err).NotTo(HaveOccurred())
					Expect(session.IsAuthenticated).To(BeFalse())
				})
			})
		})
		//table test?
		Describe("authorization header validation", func() {
			When("the hash algorithm is not supported", func() {
				It("returns and unauthorized session", func() {
					unsupportedAlgorithm := `Digest ,
					algorithm=not-a-real-algorithm,
					qop=auth,
					opaque="an-opaque-value"`
					request.Header.Add("Authorization", unsupportedAlgorithm)

					session, err := theAuthenticator.Authenticate(request)

					Expect(err).NotTo(HaveOccurred())
					Expect(session.IsAuthenticated).To(BeFalse())
				})
			})

			When("the opaque value does not match", func() {
				It("returns and unauthorized session", func() {
					unsupportedAlgorithm := `Digest ,
					algorithm=SHA-256,
					qop=auth,
					opaque="opaque-value-mismatch"`
					request.Header.Add("Authorization", unsupportedAlgorithm)

					session, err := theAuthenticator.Authenticate(request)

					Expect(err).NotTo(HaveOccurred())
					Expect(session.IsAuthenticated).To(BeFalse())
				})
			})

			When("the qop value is not supported", func() {
				It("returns and unauthorized session", func() {
					unsupportedAlgorithm := `Digest ,
					algorithm=SHA-256,
					qop=not-supported,
					opaque="an-opaque-value"`
					request.Header.Add("Authorization", unsupportedAlgorithm)

					session, err := theAuthenticator.Authenticate(request)

					Expect(err).NotTo(HaveOccurred())
					Expect(session.IsAuthenticated).To(BeFalse())
				})
			})
		})

		Describe("Call to Calculate digest", func() {
			Context("called with", func() {
				It("was called correctly", func() {
					request.Header.Add("Authorization", SuccessAuthorizationHeader)

					theAuthenticator.Authenticate(request)

					credentials, authHeader, method := fakeDigest.CalculateArgsForCall(0)
					Expect(credentials.Username).To(Equal("a-plain-username"))
					Expect(credentials.Password).To(Equal("a-password"))
					Expect(authHeader.Realm).To(Equal("a-realm"))
					Expect(authHeader.Uri).To(Equal("a-uri"))
					Expect(authHeader.Nonce).To(Equal("a-nonce-value"))
					Expect(authHeader.Nc).To(Equal("6"))
					Expect(authHeader.Cnonce).To(Equal("a-client-nonce"))
					Expect(authHeader.Qop).To(Equal("auth"))
					Expect(method).To(Equal("GET"))
				})
			})

			Context("on error", func() {
				It("returns an error and an unauthorized session", func() {
					request.Header.Add("Authorization", SuccessAuthorizationHeader)
					fakeDigest.CalculateReturns("", fmt.Errorf("calculate-digest-error"))

					session, err := theAuthenticator.Authenticate(request)

					Expect(err).To(HaveOccurred())
					Expect(session.IsAuthenticated).To(BeFalse())
				})
			})

			Context("calculated digest does not match received response", func() {
				It("returns and unauthorized session", func() {
					request.Header.Add("Authorization", SuccessAuthorizationHeader)
					fakeDigest.CalculateReturns("a-mismatching-digest", nil)

					session, err := theAuthenticator.Authenticate(request)

					Expect(err).NotTo(HaveOccurred())
					Expect(session.IsAuthenticated).To(BeFalse())
				})
			})
		})
	})
})
