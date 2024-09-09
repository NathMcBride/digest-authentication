package authenticator_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/NathMcBride/digest-authentication/src/authentication/authenticator"
	. "github.com/NathMcBride/digest-authentication/src/authentication/authenticator/fakes"
	"github.com/NathMcBride/digest-authentication/src/authentication/model"
	"github.com/NathMcBride/digest-authentication/src/providers/credential"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Authenticator", func() {
	Describe("Authenticate", func() {
		var (
			fakeCredentialsProvider *FakeCredentialProvider
			fakeDigest              *FakeDigest
			fakeUnmarshaler         *FakeUnmarshaler
			successAuthHeader       model.AuthHeader
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

			fakeUnmarshaler = &FakeUnmarshaler{}
			successAuthHeader = model.AuthHeader{
				UserID:    "a-username",
				Algorithm: "SHA-256",
				Opaque:    "an-opaque-value",
				Qop:       "auth",
				Response:  "a-digest-response",
			}

			fakeUnmarshaler.UnmarshalUnmarshaledValue(successAuthHeader)

			theAuthenticator = authenticator.Authenticator{
				Opaque:             "an-opaque-value",
				HashUserName:       true,
				CredentialProvider: fakeCredentialsProvider,
				Digest:             fakeDigest,
				Unmarshaller:       fakeUnmarshaler,
			}

			request = httptest.NewRequest("GET", "http://valid", nil)
			request.Header.Add("Authorization", "some-authorization")
		})

		It("successfully authenticates a request", func() {
			session, err := theAuthenticator.Authenticate(request)

			Expect(err).NotTo(HaveOccurred())
			Expect(session.IsAuthenticated).To(BeTrue())
			Expect(session.User.UserID).To(Equal("a-username"))
		})

		When("the authorization header is missing", func() {
			It("returns an unauthorized session", func() {
				request.Header.Del("Authorization")

				session, err := theAuthenticator.Authenticate(request)

				Expect(err).NotTo(HaveOccurred())
				Expect(session.IsAuthenticated).To(BeFalse())
			})
		})

		Describe("Call to Unmarshal", func() {
			It("is called with the expected arguments", func() {
				expected := model.AuthHeader{UserID: "a-userido"}
				fakeUnmarshaler.UnmarshalUnmarshaledValue(expected)

				theAuthenticator.Authenticate(request)

				Expect(fakeUnmarshaler.UnmarshalCallCount()).To(Equal(1))
				data, value := fakeUnmarshaler.UnmarshalArgsForCall(0)

				Expect(string(data[:])).To(Equal("some-authorization"))
				Expect(*value.(*model.AuthHeader)).To(Equal(expected))
			})

			When("unmarshaling fails", func() {
				It("returns an unauthorized session", func() {
					fakeUnmarshaler.UnmarshalReturns(fmt.Errorf("an-error"))

					session, err := theAuthenticator.Authenticate(request)

					Expect(err).NotTo(HaveOccurred())
					Expect(session.IsAuthenticated).To(BeFalse())
				})
			})
		})

		Describe("Call to GetCredentials", func() {
			Context("called with", func() {
				It("was called correctly", func() {
					theAuthenticator.Authenticate(request)

					userId, shouldHash := fakeCredentialsProvider.GetCredentialsArgsForCall(0)
					Expect(userId).To(Equal("a-username"))
					Expect(shouldHash).To(BeTrue())
				})
			})

			Context("on error", func() {
				It("returns an error and an unauthorized session", func() {
					fakeCredentialsProvider.GetCredentialsReturns(nil, false, fmt.Errorf("an-error-occurred"))

					session, err := theAuthenticator.Authenticate(request)

					Expect(err).To(HaveOccurred())
					Expect(session.IsAuthenticated).To(BeFalse())
				})
			})

			Context("on not found", func() {
				It("returns an unauthorized session", func() {
					fakeCredentialsProvider.GetCredentialsReturns(nil, false, nil)

					session, err := theAuthenticator.Authenticate(request)

					Expect(err).NotTo(HaveOccurred())
					Expect(session.IsAuthenticated).To(BeFalse())
				})
			})
		})

		DescribeTable("authorization header validation",
			func(header model.AuthHeader, expected bool) {
				fakeUnmarshaler.UnmarshalUnmarshaledValue(header)

				session, err := theAuthenticator.Authenticate(request)

				Expect(err).NotTo(HaveOccurred())
				Expect(session.IsAuthenticated).To(Equal(expected))
			},
			Entry("Should authenticate",
				model.AuthHeader{
					Algorithm: "SHA-256",
					Opaque:    "an-opaque-value",
					Qop:       "auth",
					Response:  "a-digest-response",
				}, true),
			Entry("invalid algorithm",
				model.AuthHeader{
					Algorithm: "not-a-real-algorithm",
					Qop:       "auth",
					Opaque:    "an-opaque-value",
					Response:  "a-digest-response",
				}, false),
			Entry("opaque value does not match",
				model.AuthHeader{
					Algorithm: "SHA-256",
					Qop:       "auth",
					Opaque:    "opaque-value-mismatch",
					Response:  "a-digest-response",
				}, false),
			Entry("unsupported qop",
				model.AuthHeader{
					Algorithm: "SHA-256",
					Qop:       "not-supported",
					Opaque:    "an-opaque-value",
					Response:  "a-digest-response",
				}, false),
		)

		Describe("Call to Calculate digest", func() {
			Context("called with", func() {
				It("was called correctly", func() {
					theAuthenticator.Authenticate(request)

					credentials, authHeader, method := fakeDigest.CalculateArgsForCall(0)
					Expect(credentials.Username).To(Equal("a-plain-username"))
					Expect(credentials.Password).To(Equal("a-password"))
					Expect(authHeader).To(Equal(successAuthHeader))
					Expect(method).To(Equal("GET"))
				})
			})

			Context("on error", func() {
				It("returns an error and an unauthorized session", func() {
					fakeDigest.CalculateReturns("", fmt.Errorf("calculate-digest-error"))

					session, err := theAuthenticator.Authenticate(request)

					Expect(err).To(HaveOccurred())
					Expect(session.IsAuthenticated).To(BeFalse())
				})
			})

			Context("calculated digest does not match received response", func() {
				It("returns and unauthorized session", func() {
					fakeDigest.CalculateReturns("a-mismatching-digest", nil)

					session, err := theAuthenticator.Authenticate(request)

					Expect(err).NotTo(HaveOccurred())
					Expect(session.IsAuthenticated).To(BeFalse())
				})
			})
		})
	})
})
