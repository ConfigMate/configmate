package parsers

import (
	"fmt"
)

type Parser interface {
	Parse(content []byte) (*Node, []CMParserError)
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
			"json": &jsonParser{},
			"toml": &tomlParser{},
		},
	}
}

func (p *parserProviderImpl) GetParser(format string) (Parser, error) {
	parser, ok := p.parsers[format]
	if !ok {
		return nil, fmt.Errorf("format '%s' not supported", format)
	}

	return parser, nil
}
