package store_test

import (
	"github.com/NathMcBride/digest-authentication/src/authentication/store"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client Store", func() {
	It("can add and remove clients from the store", func() {
		clientStore := store.NewClientStore()

		clientStore.Add("an-entry")
		Expect(clientStore.Has("an-entry")).To(BeTrue())

		clientStore.Delete("an-entry")
		Expect(clientStore.Has("an-entry")).To(BeFalse())
	})
})
