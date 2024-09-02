package handlers

import (
	"net/http"

	"github.com/NathMcBride/digest-authentication/src/constants"
)

type RandomKeyCreator interface {
	Create() string
}

type ChallengeCreator interface {
	Create(realm string, opaque string, nonce string, shouldHashUserName bool) (string, error)
}

type ClientStore interface {
	Add(entry string)
	Has(entry string) bool
	Delete(entry string)
}

type UnauthorizedHandler struct {
	Opaque           string
	Realm            string
	HashUserName     bool
	ClientStore      ClientStore
	RandomKey        RandomKeyCreator
	ChallengeCreator ChallengeCreator
}

func (ua *UnauthorizedHandler) HandleUnauthorized(w http.ResponseWriter, r *http.Request) {
	nonce := ua.RandomKey.Create()
	ua.ClientStore.Add(nonce)

	header, err := ua.ChallengeCreator.Create(ua.Realm, ua.Opaque, nonce, ua.HashUserName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add(constants.Authenticate, header)
	w.WriteHeader(http.StatusUnauthorized)
}
