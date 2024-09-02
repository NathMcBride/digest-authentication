package handlers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/NathMcBride/digest-authentication/src/authentication/handlers"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type FakeRandomKeyCreator struct {
	callCount     int
	createReturns struct {
		key string
	}
}

func (f *FakeRandomKeyCreator) CreateReturns(key string) {
	f.createReturns = struct{ key string }{key: key}
}

func (f *FakeRandomKeyCreator) Create() string {
	f.callCount++
	return f.createReturns.key
}

type FakeChallengeCreator struct {
	callCount         int
	createArgsForCall []struct {
		realm              string
		opaque             string
		nonce              string
		shouldHashUserName bool
	}
	createReturns struct {
		header string
		err    error
	}
}

func (f *FakeChallengeCreator) CreateCallCount() int {
	return f.callCount
}

func (f *FakeChallengeCreator) CreateArgsForCall(i int) (realm string, opaque string, nonce string, shouldHashUserName bool) {
	args := f.createArgsForCall[i]
	return args.realm, args.opaque, args.nonce, args.shouldHashUserName
}

func (f *FakeChallengeCreator) CreateReturns(header string, err error) {
	f.createReturns = struct {
		header string
		err    error
	}{header, err}
}

func (f *FakeChallengeCreator) Create(realm string, opaque string, nonce string, shouldHashUserName bool) (string, error) {
	f.callCount++
	f.createArgsForCall = append(f.createArgsForCall, struct {
		realm              string
		opaque             string
		nonce              string
		shouldHashUserName bool
	}{
		realm,
		opaque,
		nonce,
		shouldHashUserName})

	return f.createReturns.header, f.createReturns.err
}

type FakeClientStore struct {
	addCallCount   int
	addArgsForCall []struct {
		entry string
	}
}

func (f *FakeClientStore) Add(entry string) {
	f.addCallCount++
	f.addArgsForCall = append(f.addArgsForCall, struct {
		entry string
	}{
		entry: entry})
}

func (f *FakeClientStore) Has(entry string) bool {
	return true
}

func (f *FakeClientStore) Delete(entry string) {}

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

		Expect(fakeRandomKeyCreator.callCount).To(Equal(1))
		Expect(fakeClientStore.addCallCount).To(Equal(1))
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
