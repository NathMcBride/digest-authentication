package authenticator

import (
	"crypto/subtle"
	"net/http"

	"github.com/NathMcBride/web-authentication/digest/authentication/model"
	"github.com/NathMcBride/web-authentication/digest/constants"
	"github.com/NathMcBride/web-authentication/digest/headers/paramlist"
	"github.com/NathMcBride/web-authentication/digest/providers/credential"
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

type Authenticator struct {
	Opaque             string
	HashUserName       bool
	CredentialProvider CredentialProvider
	Digest             Digest
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
	err := paramlist.Unmarshal([]byte(authorization), &authHeader)
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
	/*
		// Check if the requested URI matches auth header
		if r.RequestURI != auth["uri"] {
			// We allow auth["uri"] to be a full path prefix of request-uri
			// for some reason lost in history, which is probably wrong, but
			// used to be like that for quite some time
			// (https://tools.ietf.org/html/rfc2617#section-3.2.2 explicitly
			// says that auth["uri"] is the request-uri).
			//
			// TODO: make an option to allow only strict checking.
			switch u, err := url.Parse(auth["uri"]); {
			case err != nil:
				return "", nil
			case r.URL == nil:
				return "", nil
			case len(u.Path) > len(r.URL.Path):
				return "", nil
			case !strings.HasPrefix(r.URL.Path, u.Path):
				return "", nil
			}
		}*/

	digest, err := auth.Digest.Calculate(*credentials, authHeader, r.Method)
	if err != nil {
		return notAuthenticated, err
	}

	if subtle.ConstantTimeCompare([]byte(digest), []byte(authHeader.Response)) != 1 {
		return notAuthenticated, nil
	}

	// At this point crypto checks are completed and validated.
	// Now check if the session is valid.

	// nc, err := strconv.ParseUint(auth["nc"], 16, 64)
	// if err != nil {
	// 	return "", nil
	// }

	// client, ok := da.clients[auth["nonce"]]
	// if !ok {
	// 	return "", nil
	// }
	// if client.nc != 0 && client.nc >= nc && !da.IgnoreNonceCount {
	// 	return "", nil
	// }
	// client.nc = nc
	// client.lastSeen = time.Now().UnixNano()

	// respHA2 := H(":" + auth["uri"])
	// rspauth := H(strings.Join([]string{HA1, auth["nonce"], auth["nc"], auth["cnonce"], auth["qop"], respHA2}, ":"))

	// info := fmt.Sprintf(`qop="auth", rspauth="%s", cnonce="%s", nc="%s"`, rspauth, auth["cnonce"], auth["nc"])
	// return auth["username"], &info

	session := Session{
		User: User{
			UserID: authHeader.UserID,
		},
		IsAuthenticated: true,
	}

	return session, nil
}
