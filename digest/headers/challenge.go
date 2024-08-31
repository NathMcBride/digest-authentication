package headers

import (
	"errors"
	"fmt"

	"github.com/NathMcBride/web-authentication/digest/authentication/model"
	"github.com/NathMcBride/web-authentication/digest/constants"
)

type ParamListMarshaler interface {
	Marshal(v any) ([]byte, error)
}

type DigestChallenge struct {
	Marshaler ParamListMarshaler
}

func (dc *DigestChallenge) Create(realm string, opaque string, nonce string, shouldHashUserName bool) (string, error) {
	dh := model.DigestHeader{
		Realm:     realm,
		Algorithm: constants.SHA256,
		Qop:       constants.Auth,
		Opaque:    opaque,
		Nonce:     nonce,
		UserHash:  shouldHashUserName,
	}

	result, error := dc.Marshaler.Marshal(dh)
	if error != nil {
		return "", errors.New("marshaling digest header failed")
	}

	return fmt.Sprintf(`Digest %s`, string(result[:])), nil
}
