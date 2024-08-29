package handlers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/NathMcBride/web-authentication/digest/authentication/handlers"

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

type FakeDigestCreator struct {
	callCount                  int
	createChallengeArgsForCall []struct {
		realm              string
		opaque             string
		nonce              string
		shouldHashUserName bool
	}
	createChallengeReturns struct {
		header string
		err    error
	}
}

func (f *FakeDigestCreator) CreateChallengeCallCount() int {
	return f.callCount
}

func (f *FakeDigestCreator) CreateChallengeArgsForCall(i int) (realm string, opaque string, nonce string, shouldHashUserName bool) {
	args := f.createChallengeArgsForCall[i]
	return args.realm, args.opaque, args.nonce, args.shouldHashUserName
}

func (f *FakeDigestCreator) CreateChallengeReturns(header string, err error) {
	f.createChallengeReturns = struct {
		header string
		err    error
	}{header, err}
}

func (f *FakeDigestCreator) CreateChallenge(realm string, opaque string, nonce string, shouldHashUserName bool) (string, error) {
	f.callCount++
	f.createChallengeArgsForCall = append(f.createChallengeArgsForCall, struct {
		realm              string
		opaque             string
		nonce              string
		shouldHashUserName bool
	}{
		realm,
		opaque,
		nonce,
		shouldHashUserName})

	return f.createChallengeReturns.header, f.createChallengeReturns.err
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
		fakeDigestCreator    *FakeDigestCreator
		fakeClientStore      *FakeClientStore
		recorder             *httptest.ResponseRecorder
		request              *http.Request
		unauthorizedHandler  *handlers.UnauthorizedHandler
	)

	BeforeEach(func() {
		fakeRandomKeyCreator = &FakeRandomKeyCreator{}
		fakeRandomKeyCreator.CreateReturns("a-random-value")

		fakeDigestCreator = &FakeDigestCreator{}
		fakeDigestCreator.CreateChallengeReturns("a-header-value", nil)
		fakeClientStore = &FakeClientStore{}

		request = httptest.NewRequest("GET", "http://valid", nil)
		recorder = httptest.NewRecorder()

		unauthorizedHandler = &handlers.UnauthorizedHandler{
			Opaque:        "an-opaque-value",
			Realm:         "a-realm",
			HashUserName:  true,
			ClientStore:   fakeClientStore,
			RandomKey:     fakeRandomKeyCreator,
			DigestCreator: fakeDigestCreator,
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

			Expect(fakeDigestCreator.CreateChallengeCallCount()).To(Equal(1))
			realm, opaque, nonce, hashUserName := fakeDigestCreator.CreateChallengeArgsForCall(0)
			Expect(realm).To(Equal("a-realm"))
			Expect(opaque).To(Equal("an-opaque-value"))
			Expect(nonce).To(Equal("a-random-value"))
			Expect(hashUserName).To(BeTrue())
		})

		When("creating a digest fails", func() {
			It("returns 500 StatusInternalServerError", func() {
				fakeDigestCreator.CreateChallengeReturns("", fmt.Errorf("an-error"))

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
