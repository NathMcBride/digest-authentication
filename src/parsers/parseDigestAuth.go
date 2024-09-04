package parsers

import (
	"strings"

	"github.com/NathMcBride/digest-authentication/src/constants"
)

// Test
type Parser struct {
}

func (p *Parser) Parse(auth string) (map[string]string, error) {
	const prefix = constants.Digest + " "
	if len(auth) < len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		return nil, ParsingError()
	}

	return ParseHTTPPairs(auth[len(prefix):]), nil
}
