package authenticator

import "github.com/NathMcBride/web-authentication/digest/errors"

const (
	AuthenticationErrorCode errors.ErrorCode = "AUTHENTICATION_ERROR"
	HeaderNotFound          errors.ErrorCode = "HEADER_NOT_FOUND"
)

func IsAuthenticationError(err error) bool {
	return errors.Code(err) == AuthenticationErrorCode
}

func AuthenticationError(val string) error {
	return errors.NewDomainError(AuthenticationErrorCode, "failed to authenticate: %s", val)
}
