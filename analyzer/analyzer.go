package analyzer

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ConfigMate/configmate/analyzer/types"
	"github.com/ConfigMate/configmate/parsers"
)

type Analyzer interface {
	AnalyzeConfigFiles(Files map[string]*parsers.Node, Rules []Rule) (res []Result, err error)
}

type Result struct {
	Passed        bool                         `json:"passed"`         // true if the check passed, false if it failed
	ResultComment string                       `json:"result_comment"` // an error msg or comment about the result
	TokenList     []TokenLocationWithFileAlias `json:"token_list"`     // a list of tokens that were involved in the rule
}

// TokenLocationWithFileAlias is a TokenLocation enhanced with a file alias;
// representing the file the token is in.
type TokenLocationWithFileAlias struct {
	File     string                `json:"file"`
	Location parsers.TokenLocation `json:"location"`
}

type AnalyzerImpl struct{}

func (a *AnalyzerImpl) AnalyzeConfigFiles(files map[string]*parsers.Node, rules []Rule) (res []Result, err error) {
	// Find all fields and parse them
	// optMissingFields is a map of optional fields that are missing
	_, _, err = a.findAndParseAllFields(files, rules, res)
	if err != nil {
		return nil, err
	}

	// // Check rules
	// for _, rule := range rules {
	// 	// Get the field to which the rule refers
	// 	field := fields[rule.Field]

	// 	// Evaluate checks
	// 	for _, check := range rule.Checks {

	// 	}
	// }

	return res, nil
}

func (a *AnalyzerImpl) findAndParseAllFields(files map[string]*parsers.Node, rules []Rule, res []Result) (fields map[string]types.IType, optMissingFields map[string]bool, err error) {
	// Sort rules by field lenght (shortest first)
	// This guarantees parent fields are checked before child fields
	sort.Slice(rules, func(i, j int) bool {
		return len(rules[i].Field) < len(rules[j].Field)
	})

	// Check rules and store all fields
	fields = make(map[string]types.IType)
	optMissingFields = make(map[string]bool)
	for _, rule := range rules {
		// Check if a parent field is an optional
		// field that is missing, which makes the
		// current field optional as well
		for optMissingField := range optMissingFields {
			if strings.HasPrefix(rule.Field, optMissingField) {
				optMissingFields[rule.Field] = true
				break
			}
		}
		if optMissingFields[rule.Field] {
			continue
		}

		// Separate field string into file alias and path
		fileAlias, fieldPath, err := splitFileAliasAndPath(rule.Field)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to parse field %s: %s", rule.Field, err.Error())
		}

		// Get field from file tree
		field, err := files[fileAlias].Get(fieldPath)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get field %s: %s", rule.Field, err.Error())
		} else if field == nil && rule.Optional { // Field not found and optional
			optMissingFields[rule.Field] = true
		} else { // Field found
			t, err := types.MakeType(rule.Type, field.Value)
			if err != nil {
				res = append(res, Result{
					Passed:        false,
					ResultComment: fmt.Sprintf("Field %s has incorrect type: %s", rule.Field, err.Error()),
					TokenList:     []TokenLocationWithFileAlias{makeValueTokenLocation(fileAlias, field)},
				})
			} else {
				fields[rule.Field] = t
			}
		}
	}

	return fields, optMissingFields, nil
}
