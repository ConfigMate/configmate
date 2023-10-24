package parsers

import (
	"fmt"
)

type Parser interface {
	Parse(content []byte) (*Node, error)
}

func Parse(content []byte, format string) (*Node, error) {
	// Supported parsers
	parsers := map[string]Parser{
		"json":  &JsonParser{},
		"hocon": &HoconParser{},
		"toml":  &TomlParser{},
	}

	// Check if the format is supported
	if _, ok := parsers[format]; !ok {
		return nil, fmt.Errorf("unsupported format: %s", format)
	}

	// Return the parse content
	return parsers[format].Parse(content)
}
