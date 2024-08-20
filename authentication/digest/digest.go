package digest

import (
	"strings"

	"github.com/NathMcBride/web-authentication/authentication/hasher"
	"github.com/NathMcBride/web-authentication/authentication/model"
	"github.com/NathMcBride/web-authentication/providers/credential"
)

func CreateDigest(credentials credential.Credentials, authHeader model.AuthHeader, Method string) (string, error) {
	HA1, err := hasher.H(credentials.Username + ":" + authHeader.Realm + ":" + credentials.Password)
	if err != nil {
		return "", err
	}

	HA2, err := hasher.H(Method + ":" + authHeader.Uri)
	if err != nil {
		return "", err
	}

	list := []string{HA1, authHeader.Nonce, authHeader.Nc, authHeader.Cnonce, authHeader.Qop, HA2}
	digest := strings.Join(list, ":")

	KD, err := hasher.H(digest)
	if err != nil {
		return "", err
	}

	return KD, nil
}
