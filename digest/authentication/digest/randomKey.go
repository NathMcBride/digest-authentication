package digest

import (
	"crypto/rand"
	"encoding/base64"
)

func RandomKey() string {
	b := make([]byte, 12)

	_, err := rand.Read(b)
	if err != nil {
		panic("rand.Read() failed")
	}

	return base64.StdEncoding.EncodeToString(b)
}
