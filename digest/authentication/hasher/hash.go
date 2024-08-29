package hasher

import (
	"crypto/sha256"
	"fmt"
)

type Hasher struct {
}

func (h *Hasher) Hash(data string) (string, error) {
	digest := sha256.New()

	_, err := digest.Write([]byte(data))
	if err != nil {
		return "", HashingError()
	}

	result := fmt.Sprintf("%x", digest.Sum(nil))
	return result, nil
}
