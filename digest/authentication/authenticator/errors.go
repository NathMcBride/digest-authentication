package authenticator

import "github.com/NathMcBride/web-authentication/digest/domainerror"

const (
	AuthenticationErrorCode domainerror.ErrorCode = "AUTHENTICATION_ERROR"
	HeaderNotFound          domainerror.ErrorCode = "HEADER_NOT_FOUND"
)

func IsAuthenticationError(err error) bool {
	return domainerror.Code(err) == AuthenticationErrorCode
}

func AuthenticationError(val string) error {
	return domainerror.NewDomainError(AuthenticationErrorCode, "failed to authenticate: %s", val)
}
