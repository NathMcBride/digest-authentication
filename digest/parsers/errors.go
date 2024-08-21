package parsers

import "github.com/NathMcBride/web-authentication/digest/errors"

const (
	ParsingErrorCode errors.ErrorCode = "Parsing_ERROR"
)

func IsParsingError(err error) bool {
	return errors.Code(err) == ParsingErrorCode
}

func ParsingError() error {
	return errors.NewDomainError(ParsingErrorCode, "failed to parse")
}
