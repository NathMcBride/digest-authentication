package handlers

import (
	"log"
	"net/http"

	"github.com/NathMcBride/web-authentication/digest/authentication/digest"
	"github.com/NathMcBride/web-authentication/digest/authentication/store"
	"github.com/NathMcBride/web-authentication/digest/constants"
)

type UnauthorizedHandler struct {
	Opaque       string
	Realm        string
	HashUserName bool
	ClientStore  *store.ClientStore
}

func (ua *UnauthorizedHandler) HandleUnauthorized(w http.ResponseWriter, r *http.Request) {
	nonce := digest.RandomKey()
	ua.ClientStore.Add(nonce)

	header, err := digest.MakeHeader(ua.Realm, ua.Opaque, nonce, ua.HashUserName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add(constants.Authenticate, header)
	w.WriteHeader(http.StatusUnauthorized)
	log.Default().Print("Digest authentication needed.")
}
