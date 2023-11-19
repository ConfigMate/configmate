package parsers

import (
	"fmt"
)

type Parser interface {
	Parse(content []byte) (*Node, error)
}

type ParserProvider interface {
	GetParser(format string) (Parser, error)
}

type parserProviderImpl struct {
	parsers map[string]Parser
}

func NewParserProvider() ParserProvider {
	return &parserProviderImpl{
		parsers: map[string]Parser{
			"json": &JsonParser{},
		},
	}
}

func (p *parserProviderImpl) GetParser(format string) (Parser, error) {
	parser, ok := p.parsers[format]
	if !ok {
		return nil, fmt.Errorf("parser for format %s not found", format)
	}

	return parser, nil
}
