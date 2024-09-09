package fakes

import "net/http"

type FakeHandler struct {
	callCount   int
	argsForCall []struct {
		writer  http.ResponseWriter
		request *http.Request
	}
}

func (fh *FakeHandler) ServeHTTPCallCount() int {
	return fh.callCount
}

func (fh *FakeHandler) ServeHTTPArgsForCall(i int) (writer http.ResponseWriter, request *http.Request) {
	args := fh.argsForCall[i]
	return args.writer, args.request
}

func (fh *FakeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fh.callCount++
	fh.argsForCall = append(fh.argsForCall, struct {
		writer  http.ResponseWriter
		request *http.Request
	}{writer: w, request: r})
}
