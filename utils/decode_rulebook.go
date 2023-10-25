package utils

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/ConfigMate/configmate/analyzer"
)

func DecodeRulebook(ruleBookContent []byte) (*analyzer.Rulebook, error) {
	// Initialize Rulebook struct
	var rulebook analyzer.Rulebook

	// Decode the TOML data into the Rulebook struct
	if _, err := toml.Decode(string(ruleBookContent), &rulebook); err != nil {
		return nil, fmt.Errorf("error decoding file into a rulebook object: %v", err)
	}

	return &rulebook, nil
}
