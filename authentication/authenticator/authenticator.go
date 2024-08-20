package authenticator

import (
	"crypto/subtle"
	"net/http"

	"github.com/NathMcBride/web-authentication/authentication/digest"
	"github.com/NathMcBride/web-authentication/authentication/model"
	"github.com/NathMcBride/web-authentication/constants"
	"github.com/NathMcBride/web-authentication/headers/paramlist"
	"github.com/NathMcBride/web-authentication/providers/credential"
)

type User struct {
	UserID          string
	IsAuthenticated bool
}

type Authenticator struct {
	Opaque             string
	HashUserName       bool
	CredentialProvider *credential.CredentialProvider
}

func (auth *Authenticator) Authenticate(r *http.Request) (User, error) {
	notAuthenticated := func(err error) (User, error) {
		return User{UserID: "", IsAuthenticated: false}, err
	}

	authorization := r.Header.Get(constants.Authorization)
	if authorization == "" {
		err := AuthenticationError("authorization not found")
		return notAuthenticated(err)
	}

	authHeader := model.AuthHeader{}
	err := paramlist.Unmarshal([]byte(authorization), &authHeader)
	if err != nil {
		return notAuthenticated(err)
	}

	credentials, err := auth.CredentialProvider.GetCredentials(authHeader.UserID, auth.HashUserName)
	if err != nil {
		return notAuthenticated(err)
	}

	if authHeader.Algorithm != constants.SHA256 ||
		authHeader.Opaque != auth.Opaque ||
		authHeader.Qop != constants.Auth {
		err := AuthenticationError("header mismatch")
		return notAuthenticated(err)

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

	digest, err := digest.CreateDigest(*credentials, authHeader, r.Method)
	if err != nil {
		return notAuthenticated(err)
	}

	if subtle.ConstantTimeCompare([]byte(digest), []byte(authHeader.Response)) != 1 {
		err := AuthenticationError("calculated digest does not match response")
		return notAuthenticated(err)
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

	return User{UserID: authHeader.UserID, IsAuthenticated: true}, nil
}
