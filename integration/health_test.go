package integration_test

import (
	"net/http"

	"github.com/NathMcBride/digest-authentication/integration/client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Health", func() {
	Describe("GET /health", func() {
		It("returns 200 StatusOk", func() {
			client := client.Client{
				Endpoint:           "http://localhost:8080",
				ShouldAuthenticate: false}

			resp, err := client.GetHealth()
			Expect(err).NotTo(HaveOccurred())

			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})
	})
})
