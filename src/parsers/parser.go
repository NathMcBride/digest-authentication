package parsers

import (
	"strings"
)

type Parser struct {
}

func (p *Parser) ParseList(toParse string, prefix string) (map[string]string, error) {
	if prefix == "" {
		return HTTPPairs(toParse[:]), nil
	}

	if len(toParse) < len(prefix) || !strings.EqualFold(toParse[:len(prefix)], prefix) {
		return nil, ParsingError()
	}

	return HTTPPairs(toParse[len(prefix):]), nil
}
