package handlers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/NathMcBride/digest-authentication/src/authentication/handlers"
	. "github.com/NathMcBride/digest-authentication/src/authentication/handlers/fakes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Unauthorized handler", func() {
	var (
		fakeRandomKeyCreator *FakeRandomKeyCreator
		fakeChallengeCreator *FakeChallengeCreator
		fakeClientStore      *FakeClientStore
		recorder             *httptest.ResponseRecorder
		request              *http.Request
		unauthorizedHandler  *handlers.UnauthorizedHandler
	)

	BeforeEach(func() {
		fakeRandomKeyCreator = &FakeRandomKeyCreator{}
		fakeRandomKeyCreator.CreateReturns("a-random-value")

		fakeChallengeCreator = &FakeChallengeCreator{}
		fakeChallengeCreator.CreateReturns("a-header-value", nil)
		fakeClientStore = &FakeClientStore{}

		request = httptest.NewRequest("GET", "http://valid", nil)
		recorder = httptest.NewRecorder()

		unauthorizedHandler = &handlers.UnauthorizedHandler{
			Opaque:           "an-opaque-value",
			Realm:            "a-realm",
			HashUserName:     true,
			ClientStore:      fakeClientStore,
			RandomKey:        fakeRandomKeyCreator,
			ChallengeCreator: fakeChallengeCreator,
		}
	})

	It("adds the nonce to client store", func() {
		unauthorizedHandler.HandleUnauthorized(recorder, request)

		Expect(fakeRandomKeyCreator.CreateCallCount()).To(Equal(1))
		Expect(fakeClientStore.AddCallCount()).To(Equal(1))
	})

	Context("creating a Digest challenge", func() {
		It("calls create with the correct arguments", func() {
			unauthorizedHandler.HandleUnauthorized(recorder, request)

			Expect(fakeChallengeCreator.CreateCallCount()).To(Equal(1))
			realm, opaque, nonce, hashUserName := fakeChallengeCreator.CreateArgsForCall(0)
			Expect(realm).To(Equal("a-realm"))
			Expect(opaque).To(Equal("an-opaque-value"))
			Expect(nonce).To(Equal("a-random-value"))
			Expect(hashUserName).To(BeTrue())
		})

		When("creating a digest fails", func() {
			It("returns 500 StatusInternalServerError", func() {
				fakeChallengeCreator.CreateReturns("", fmt.Errorf("an-error"))

				unauthorizedHandler.HandleUnauthorized(recorder, request)

				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
			})
		})
	})

	It("responds with the expected header values", func() {
		unauthorizedHandler.HandleUnauthorized(recorder, request)

		header := recorder.Header().Get("WWW-Authenticate")

		Expect(header).To(Equal("a-header-value"))
		Expect(recorder.Code).To(Equal(http.StatusUnauthorized))
	})
})
