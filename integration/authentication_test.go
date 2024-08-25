package integration_test

import (
	"net/http"

	"github.com/NathMcBride/web-authentication/digest/authentication/model"
	"github.com/NathMcBride/web-authentication/digest/headers/paramlist"
	"github.com/NathMcBride/web-authentication/integration/client"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Authentication", func() {

	Describe("GET /protected", func() {
		Context("with no credentials", func() {
			client := client.Client{
				Endpoint:           "http://localhost:8080",
				ShouldAuthenticate: false}

			It("responds with 401 StatusUnauthorized", func() {
				resp, err := client.GetProtected()

				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
			})

			It("responds with a digest authentication challenge", func() {
				resp, err := client.GetProtected()
				Expect(err).NotTo(HaveOccurred())

				digestHeader := model.DigestHeader{}
				buffer := []byte(resp.Header.Get("WWW-Authenticate"))
				err = paramlist.Unmarshal(buffer, &digestHeader)
				Expect(err).NotTo(HaveOccurred())

				Expect(digestHeader.Realm).To(Equal("A-Realm"))
				Expect(digestHeader.Algorithm).To(Equal("SHA-256"))
				Expect(digestHeader.Qop).To(Equal("auth"))
				Expect(digestHeader.UserHash).To(BeTrue())
				Expect(digestHeader.Nonce).ToNot(BeEmpty())
				Expect(digestHeader.Opaque).ToNot(BeEmpty())
			})
		})

		Context("with an invalid username", func() {
			It("responds with 401 StatusUnauthorized", func() {
				client := client.Client{
					Endpoint:           "http://localhost:8080",
					ShouldAuthenticate: true,
					Username:           "invalid",
					Password:           "password"}

				resp, err := client.GetProtected()

				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
			})
		})

		Context("with an invalid password", func() {
			It("responds with 401 StatusUnauthorized", func() {
				client := client.Client{
					Endpoint:           "http://localhost:8080",
					ShouldAuthenticate: true,
					Username:           "jim",
					Password:           "invalid"}

				resp, err := client.GetProtected()

				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
			})
		})

		Context("with valid credentials", func() {
			It("responds with 200 StatusOk", func() {
				client := client.Client{
					Endpoint:           "http://localhost:8080",
					ShouldAuthenticate: true,
					Username:           "jim",
					Password:           "password"}

				resp, err := client.GetProtected()

				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})
		})
	})
})
