package hasher

import (
	"crypto/sha256"
	"fmt"
	"io"
)

type CryptoHash interface {
	io.Writer
	Sum(b []byte) []byte
	Reset()
	Size() int
	BlockSize() int
}

type CryptoFactory interface {
	New() CryptoHash
}

type Sha256Factory struct{}

func (t *Sha256Factory) New() CryptoHash {
	return sha256.New()
}

type Hash struct {
	CryptoFactory CryptoFactory
}

func (h *Hash) Do(data string) (string, error) {
	cryptoHash := h.CryptoFactory.New()
	_, err := cryptoHash.Write([]byte(data))
	if err != nil {
		return "", HashingError()
	}

	result := fmt.Sprintf("%x", cryptoHash.Sum(nil))
	return result, nil
}
