package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/NathMcBride/web-authentication/digest/authentication/authenticator"
	"github.com/NathMcBride/web-authentication/digest/authentication/contexts"
	"github.com/NathMcBride/web-authentication/digest/authentication/digest"
	"github.com/NathMcBride/web-authentication/digest/authentication/store"
	"github.com/NathMcBride/web-authentication/digest/constants"
	"github.com/NathMcBride/web-authentication/digest/providers/credential"
	"github.com/NathMcBride/web-authentication/digest/providers/secret"
	"github.com/NathMcBride/web-authentication/digest/providers/username"
)

type Authenticate struct {
	Opaque        string
	Realm         string
	HashUserName  bool
	ClientStore   *store.ClientStore
	Authenticator *authenticator.Authenticator
}

func (a *Authenticate) HandleUnauthorized(w http.ResponseWriter, r *http.Request) {
	nonce := digest.RandomKey()
	a.ClientStore.Add(nonce)

	header, err := digest.MakeHeader(a.Realm, a.Opaque, nonce, a.HashUserName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add(constants.Authenticate, header)
	w.WriteHeader(http.StatusUnauthorized)
	log.Default().Print("Digest authentication needed.")
}

func (a *Authenticate) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := a.Authenticator.Authenticate(r)
		if err != nil {
			if authenticator.IsAuthenticationError(err) {
				a.HandleUnauthorized(w, r)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !user.IsAuthenticated {
			a.HandleUnauthorized(w, r)
			return
		}

		ctx := contexts.WithUser(context.Background(), &user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type Options struct {
	Realm              string
	Opaque             string
	ShouldHashUsername bool
}

func NewDigestAuth(o Options) *Authenticate {
	secretProvider := secret.SecretProviderProvider{}
	usernameProvider := username.UsernameProvider{Realm: o.Realm}
	clientStore := store.NewClientStore()

	credentialProvider := credential.CredentialProvider{
		SecretProvider:   &secretProvider,
		UsernameProvider: &usernameProvider,
	}

	authenticator := authenticator.Authenticator{
		Opaque:             o.Opaque,
		HashUserName:       o.ShouldHashUsername,
		CredentialProvider: &credentialProvider,
	}

	authenticate := Authenticate{
		Opaque:        o.Opaque,
		Realm:         o.Realm,
		HashUserName:  o.ShouldHashUsername,
		ClientStore:   clientStore,
		Authenticator: &authenticator,
	}

	return &authenticate
}
