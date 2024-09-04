package authenticator

import (
	"crypto/subtle"
	"net/http"

	"github.com/NathMcBride/digest-authentication/src/authentication/model"
	"github.com/NathMcBride/digest-authentication/src/constants"
	"github.com/NathMcBride/digest-authentication/src/providers/credential"
)

type User struct {
	UserID string
}

type Session struct {
	User            User
	IsAuthenticated bool
}

type CredentialProvider interface {
	GetCredentials(userID string, useHash bool) (*credential.Credentials, bool, error)
}

type Digest interface {
	Calculate(
		credentials credential.Credentials,
		authHeader model.AuthHeader,
		Method string,
	) (string, error)
}

type Unmarshaler interface {
	Unmarshal(data []byte, v any) error
}

type Authenticator struct {
	Opaque             string
	HashUserName       bool
	CredentialProvider CredentialProvider
	Digest             Digest
	Unmarshaller       Unmarshaler
}

func (auth *Authenticator) Authenticate(r *http.Request) (Session, error) {
	notAuthenticated := Session{
		User:            User{},
		IsAuthenticated: false,
	}

	authorization := r.Header.Get(constants.Authorization)
	if authorization == "" {
		return notAuthenticated, nil
	}

	authHeader := model.AuthHeader{}
	err := auth.Unmarshaller.Unmarshal([]byte(authorization), &authHeader)
	if err != nil {
		return notAuthenticated, nil
	}

	credentials, found, err := auth.CredentialProvider.GetCredentials(authHeader.UserID, auth.HashUserName)
	if err != nil || !found {
		return notAuthenticated, err
	}

	if authHeader.Algorithm != constants.SHA256 ||
		authHeader.Opaque != auth.Opaque ||
		authHeader.Qop != constants.Auth {
		return notAuthenticated, nil
	}

	//Add checks to validate uri

	digest, err := auth.Digest.Calculate(*credentials, authHeader, r.Method)
	if err != nil {
		return notAuthenticated, err
	}

	if subtle.ConstantTimeCompare([]byte(digest), []byte(authHeader.Response)) != 1 {
		return notAuthenticated, nil
	}

	//Add checks to validate session

	session := Session{
		User: User{
			UserID: authHeader.UserID,
		},
		IsAuthenticated: true,
	}

	return session, nil
}
