package fakes

import "net/http"

type FakeUnauthorizedHandler struct {
	callCount   int
	argsForCall []struct {
		writer  http.ResponseWriter
		request *http.Request
	}
}

func (fu *FakeUnauthorizedHandler) HandleUnauthorizedCallCount() int {
	return fu.callCount
}

func (fu *FakeUnauthorizedHandler) HandleUnauthorized(w http.ResponseWriter, r *http.Request) {
	fu.callCount++
	fu.argsForCall = append(fu.argsForCall, struct {
		writer  http.ResponseWriter
		request *http.Request
	}{writer: w, request: r})
}
