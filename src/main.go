package main

import (
	"log"
	"net/http"

	"github.com/NathMcBride/digest-authentication/src/authentication/authenticator"
	"github.com/NathMcBride/digest-authentication/src/authentication/contexts"
	"github.com/NathMcBride/digest-authentication/src/authentication/digest"
	"github.com/NathMcBride/digest-authentication/src/authentication/middleware"
)

func somethingProtected(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Something protected"))
}

func someThingWithSessionHandler() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			session := contexts.GetSession(ctx)
			someThingWithSession(w, r, session)
		},
	)
}

func someThingWithSession(w http.ResponseWriter, r *http.Request, session *authenticator.Session) {
	w.Write([]byte("Something that needs a user"))
}

func somethingPublic(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Something public"))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	key := digest.RandomKey{}
	authMiddleware := middleware.NewDigestAuth(
		"A-Realm",
		key.Create(),
		true)

	mux := http.NewServeMux()
	mux.Handle("/protected", authMiddleware(someThingWithSessionHandler()))
	mux.Handle("/public", http.HandlerFunc(somethingPublic))
	mux.Handle("/health", http.HandlerFunc(handleHealth))
	mux.Handle("/", http.NotFoundHandler())

	log.Fatal(http.ListenAndServe("localhost:8080", mux))
}
