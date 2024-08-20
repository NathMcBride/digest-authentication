package hasher

import (
	"crypto/sha256"
	"fmt"
)

func H(data string) (string, error) {
	digest := sha256.New()

	_, err := digest.Write([]byte(data))
	if err != nil {
		return "", HashingError()
	}

	result := fmt.Sprintf("%x", digest.Sum(nil))
	return result, nil
}
