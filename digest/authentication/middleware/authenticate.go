package middleware

import (
	"context"
	"net/http"

	"github.com/NathMcBride/web-authentication/digest/authentication/authenticator"
	"github.com/NathMcBride/web-authentication/digest/authentication/contexts"
	"github.com/NathMcBride/web-authentication/digest/authentication/digest"
	"github.com/NathMcBride/web-authentication/digest/authentication/handlers"
	"github.com/NathMcBride/web-authentication/digest/authentication/hasher"
	"github.com/NathMcBride/web-authentication/digest/authentication/store"
	"github.com/NathMcBride/web-authentication/digest/headers"
	"github.com/NathMcBride/web-authentication/digest/headers/paramlist"
	"github.com/NathMcBride/web-authentication/digest/providers/credential"
	"github.com/NathMcBride/web-authentication/digest/providers/secret"
	"github.com/NathMcBride/web-authentication/digest/providers/username"
)

type HandleUnauthorized interface {
	HandleUnauthorized(w http.ResponseWriter, r *http.Request)
}

type Authenticator interface {
	Authenticate(r *http.Request) (authenticator.Session, error)
}

type Authenticate struct {
	UnauthorizedHandler HandleUnauthorized
	Authenticator       Authenticator
}

func (a *Authenticate) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := a.Authenticator.Authenticate(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !session.IsAuthenticated {
			a.UnauthorizedHandler.HandleUnauthorized(w, r)
			return
		}

		ctx := contexts.WithSession(context.Background(), &session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func NewDigestAuth(Realm string, Opaque string, ShouldHashUsername bool) func(http.Handler) http.Handler {
	secretProvider := secret.SecretProviderProvider{}
	usernameProvider := username.UsernameProvider{Realm: Realm}
	clientStore := store.NewClientStore()
	randomKeyCreator := digest.RandomKey{}
	sha256Factory := hasher.Sha256Factory{}
	digest := digest.Digest{
		Hasher: &hasher.Hash{
			CryptoFactory: &sha256Factory,
		},
	}

	challenge := headers.DigestChallenge{
		Marshaler: &paramlist.Marshaler{},
	}
	unauthorizedHandler := handlers.UnauthorizedHandler{
		Opaque:           Opaque,
		Realm:            Realm,
		HashUserName:     ShouldHashUsername,
		ClientStore:      &clientStore,
		RandomKey:        &randomKeyCreator,
		ChallengeCreator: &challenge,
	}

	credentialProvider := credential.CredentialProvider{
		SecretProvider:   &secretProvider,
		UsernameProvider: &usernameProvider,
	}

	authenticator := authenticator.Authenticator{
		Opaque:             Opaque,
		HashUserName:       ShouldHashUsername,
		CredentialProvider: &credentialProvider,
		Digest:             &digest,
	}

	authenticate := Authenticate{
		UnauthorizedHandler: &unauthorizedHandler,
		Authenticator:       &authenticator,
	}

	return authenticate.RequireAuth
}
