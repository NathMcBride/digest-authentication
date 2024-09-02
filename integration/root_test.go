package integration_test

import (
	"net/http"

	"github.com/NathMcBride/digest-authentication/integration/client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Root", func() {
	Describe("GET /", func() {
		It("returns 404 StatusNotFound", func() {
			client := client.Client{
				Endpoint:           "http://localhost:8080",
				ShouldAuthenticate: false}

			resp, err := client.GetRoot()
			Expect(err).NotTo(HaveOccurred())

			Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
		})
	})
})
