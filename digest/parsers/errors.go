package parsers

import (
	"github.com/NathMcBride/web-authentication/digest/domainerror"
)

const (
	ParsingErrorCode domainerror.ErrorCode = "Parsing_ERROR"
)

func IsParsingError(err error) bool {
	return domainerror.Code(err) == ParsingErrorCode
}

func ParsingError() error {
	return domainerror.NewDomainError(ParsingErrorCode, "failed to parse")
}
