package domainerror

import "fmt"

type ErrorCode string

type domainError struct {
	error
	errorCode ErrorCode
}

func (e domainError) Error() string {
	return fmt.Sprintf("%s: %s", e.errorCode, e.error.Error())
}

func Code(err error) ErrorCode {
	if err == nil {
		return ""
	}

	if e, ok := err.(domainError); ok {
		return e.errorCode
	}

	return ""
}

func NewDomainError(errorCode ErrorCode, format string, args ...interface{}) error {
	return domainError{
		error:     fmt.Errorf(format, args...),
		errorCode: errorCode,
	}
}
