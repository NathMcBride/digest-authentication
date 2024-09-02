package middleware

import (
	"context"
	"net/http"

	"github.com/NathMcBride/digest-authentication/src/authentication/authenticator"
	"github.com/NathMcBride/digest-authentication/src/authentication/contexts"
	"github.com/NathMcBride/digest-authentication/src/authentication/digest"
	"github.com/NathMcBride/digest-authentication/src/authentication/handlers"
	"github.com/NathMcBride/digest-authentication/src/authentication/hasher"
	"github.com/NathMcBride/digest-authentication/src/authentication/store"
	"github.com/NathMcBride/digest-authentication/src/headers"
	"github.com/NathMcBride/digest-authentication/src/headers/paramlist"
	"github.com/NathMcBride/digest-authentication/src/headers/paramlist/structinfo"
	"github.com/NathMcBride/digest-authentication/src/headers/paramlist/structmarshal"
	"github.com/NathMcBride/digest-authentication/src/providers/credential"
	"github.com/NathMcBride/digest-authentication/src/providers/secret"
	"github.com/NathMcBride/digest-authentication/src/providers/username"
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

// refactor
func NewDigestAuth(Realm string, Opaque string, ShouldHashUsername bool) func(http.Handler) http.Handler {
	secretProvider := secret.SecretProviderProvider{}
	usernameProvider := username.UsernameProvider{Realm: Realm}

	structInfo := structinfo.StructInfo{}
	structMarshal := structmarshal.StructMarshal{}
	challenge := headers.DigestChallenge{
		Marshaler: &paramlist.Marshaler{
			StructInfoer:    &structInfo,
			StructMarshaler: &structMarshal,
		},
	}
	clientStore := store.NewClientStore()
	randomKeyCreator := digest.RandomKey{}
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
	sha256Factory := hasher.Sha256Factory{}
	digest := digest.Digest{
		Hasher: &hasher.Hash{
			CryptoFactory: &sha256Factory,
		},
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
