package hasher

import (
	"github.com/NathMcBride/web-authentication/digest/domainerror"
)

const (
	HashingErrorCode domainerror.ErrorCode = "HASHING_ERROR"
)

func IsHashingError(err error) bool {
	return domainerror.Code(err) == HashingErrorCode
}

func HashingError() error {
	return domainerror.NewDomainError(HashingErrorCode, "failed to hash")
}
