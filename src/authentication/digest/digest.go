package digest

import (
	"strings"

	"github.com/NathMcBride/digest-authentication/src/authentication/model"
	"github.com/NathMcBride/digest-authentication/src/providers/credential"
)

type Hasher interface {
	Do(data string) (string, error)
}

type Digest struct {
	Hasher Hasher
}

func (d *Digest) Calculate(credentials credential.Credentials, authHeader model.AuthHeader, Method string) (string, error) {
	HA1, err := d.Hasher.Do(
		strings.Join([]string{
			credentials.Username,
			authHeader.Realm,
			credentials.Password}, ":"))
	if err != nil {
		return "", err
	}

	HA2, err := d.Hasher.Do(Method + ":" + authHeader.Uri)
	if err != nil {
		return "", err
	}

	KD, err := d.Hasher.Do(
		strings.Join([]string{
			HA1,
			authHeader.Nonce,
			authHeader.Nc,
			authHeader.Cnonce,
			authHeader.Qop,
			HA2}, ":"))
	if err != nil {
		return "", err
	}

	return KD, nil
}
