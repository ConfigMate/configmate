package analyzer

import (
	"fmt"

	"github.com/ConfigMate/configmate/parsers"
)

type Analyzer interface {
	Analyze(rb Rulebook) (res []Result, err error)
}

type Result struct {
	Passed        bool            `json:"passed"`         // true if the check passed, false if it failed
	ResultComment string          `json:"result_comment"` // an error msg or comment about the result
	TokenList     []TokenLocation `json:"token_list"`     // a list of tokens that were involved in the rule
}

type TokenLocation struct {
	File   string `json:"file"`
	Line   int    `json:"line"`
	Column int    `json:"column"`
	Length int    `json:"length"`
}

type Check interface {
	Check(values []interface{}) (bool, error)
}

type AnalyzerImpl struct {
	checks  map[string]Check              // map of check names to checks
	parsers map[FileFormat]parsers.Parser // map of file formats to parsers
}

func (a *AnalyzerImpl) Analyze(rb Rulebook) (res []Result, err error) {
	files := make(map[string]parsers.ConfigFile)

	// Parse files
	for key, filePath := range rb.Files {
		// Get parser for file format
		parser, ok := a.parsers[getFileFormat(filePath)]
		if !ok {
			return nil, fmt.Errorf("no parser found for file format of file %s", filePath)
		}

		// Parse file
		configFile, err := parser.Parse(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to parse file %s: %s", filePath, err.Error())
		}

		// Add file to map
		files[key] = configFile
	}

	// Check rules
	for _, rule := range rb.Rules {
		args := make([]interface{}, 0)
		// Get values for arguments
		for _, arg := range rule.Args {
			if isFile, alias, key := decodeFileValue(arg); isFile {
				// Get value from config file
				value, err := getValueFromConfigFile(files[alias], key)
				if err != nil {
					return nil, fmt.Errorf("failed to get value from config file: %s", err.Error())
				}
				args = append(args, value)

			} else if isLiteral, value := decodeLiteralValue(arg); isLiteral {
				args = append(args, value)

			} else {
				return nil, fmt.Errorf("failed to decode argument %s", arg)
			}
		}

		// Apply check
		passed, err := a.checks[rule.CheckName].Check(args)
		if err != nil {
			return nil, fmt.Errorf("failed to apply check %s: %s", rule.CheckName, err.Error())
		}

		// Add result
		res = append(res, Result{
			Passed:        passed,
			ResultComment: rule.Description,
			TokenList:     nil,
		})
	}

	return res, nil
}
