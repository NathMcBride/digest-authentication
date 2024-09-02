package errors

import (
	"github.com/NathMcBride/digest-authentication/src/domainerror"
)

const (
	MarshalErrorCode   domainerror.ErrorCode = "MARSHAL_ERROR"
	UnmarshalErrorCode domainerror.ErrorCode = "UNMARSHAL_ERROR"
)

func IsMarshalError(err error) bool {
	return domainerror.Code(err) == MarshalErrorCode
}

func MarshalError(val string) error {
	return domainerror.NewDomainError(MarshalErrorCode, "failed to marshal: %s", val)
}

func IsUnmarshallError(err error) bool {
	return domainerror.Code(err) == UnmarshalErrorCode
}

func UnmarshalError(val string) error {
	return domainerror.NewDomainError(UnmarshalErrorCode, "failed to unmarshal: %s", val)
}
