package utils

import (
	"testing"

	"github.com/ConfigMate/configmate/analyzer"
	"github.com/ConfigMate/configmate/parsers"
)

func TestPrintResults(t *testing.T) {
	// TODO: put some real example
	results := analyzer.Result{
		Passed:        false,
		ResultComment: "Check passed",
		Rule: &analyzer.Rule{
			Field:       "field",
			Description: "description",
			Type:        "string",
			Optional:    false,
			Default:     "default",
			Checks:      []string{"check1", "check2"},
			Notes:       "notes",
		},
		TokenList: []analyzer.TokenLocationWithFileAlias{
			{
				File: "File1",
				Location: parsers.TokenLocation{
					Start: parsers.CharLocation{Line: 1, Column: 2},
					End:   parsers.CharLocation{Line: 1, Column: 4},
				},
			},
		},
        CheckNumber: 1,
	}

	// Just print for now
	formatted := FormatResult(results)
    t.Errorf("\n" + formatted)
}
