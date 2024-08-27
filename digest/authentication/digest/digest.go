package digest

import (
	"errors"
	"fmt"
	"strings"

	"github.com/NathMcBride/web-authentication/digest/authentication/hasher"
	"github.com/NathMcBride/web-authentication/digest/authentication/model"
	"github.com/NathMcBride/web-authentication/digest/constants"
	"github.com/NathMcBride/web-authentication/digest/headers/paramlist"
	"github.com/NathMcBride/web-authentication/digest/providers/credential"
)

type Digest struct{}

func (d *Digest) Calculate(credentials credential.Credentials, authHeader model.AuthHeader, Method string) (string, error) {
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

func MakeHeader(realm string, opaque string, nonce string, shouldHashUserName bool) (string, error) {
	dh := model.DigestHeader{
		Realm:     realm,
		Algorithm: constants.SHA256,
		Qop:       constants.Auth,
		Opaque:    opaque,
		Nonce:     nonce,
		UserHash:  shouldHashUserName,
	}

	result, error := paramlist.Marshal(dh)
	if error != nil {
		return "", errors.New("marshaling digest header failed")
	}

	return fmt.Sprintf(`Digest %s`, string(result[:])), nil
}
