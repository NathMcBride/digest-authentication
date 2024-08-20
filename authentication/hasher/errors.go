package hasher

import "github.com/NathMcBride/web-authentication/errors"

const (
	HashingErrorCode errors.ErrorCode = "HASHING_ERROR"
)

func IsHashingError(err error) bool {
	return errors.Code(err) == HashingErrorCode
}

func HashingError() error {
	return errors.NewDomainError(HashingErrorCode, "failed to hash")
}
