package parsers

import (
	"strings"

	"github.com/NathMcBride/web-authentication/digest/constants"
)

func ParseDigestAuth(auth string) (map[string]string, error) {
	const prefix = constants.Digest + " "
	if len(auth) < len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		return nil, ParsingError()
	}

	return ParseHTTPPairs(auth[len(prefix):]), nil
}
