package paramlist

import "github.com/NathMcBride/web-authentication/digest/errors"

const (
	MarshalErrorCode   errors.ErrorCode = "MARSHAL_ERROR"
	UnmarshalErrorCode errors.ErrorCode = "UNMARSHAL_ERROR"
)

func IsMarshalError(err error) bool {
	return errors.Code(err) == MarshalErrorCode
}

func MarshalError(val string) error {
	return errors.NewDomainError(MarshalErrorCode, "failed to marshal: %s", val)
}

func IsUnmarshallError(err error) bool {
	return errors.Code(err) == UnmarshalErrorCode
}

func UnmarshalError(val string) error {
	return errors.NewDomainError(UnmarshalErrorCode, "failed to unmarshal: %s", val)
}
