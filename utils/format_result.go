package utils

import (
	"fmt"

	"github.com/ConfigMate/configmate/analyzer"
)

func FormatResult(res analyzer.Result, fileLinesMap map[string]map[int]string) string {
	var passed, comment, check, fieldType, optional string

	// Format of the result
	format := "%s:\n\tField:\t%s %s %s \n\tCheck:\t%s %s \n\tDefault: %v \n\tNotes:  %s\n"

	// format values
	if res.Passed {
		passed = ColorText("PASSED", Green)
	} else {
		passed = ColorText("FAILED", Red)
		comment = ColorText(res.ResultComment, Red)
		comment = fmt.Sprintf("- %s", comment)
	}

	// get other values
	check = ColorText(res.Rule.Checks[res.CheckNum], Cyan)
	fieldType = ColorText(res.Rule.Type, Blue)
	if res.Rule.Optional {
		optional = ColorText("optional", Yellow)
	} else {
		optional = ColorText("required", Red)
	}

	// Format the values
	formatted := fmt.Sprintf(format, passed, res.Rule.Field, fieldType, optional, check, comment, res.Rule.Default, res.Rule.Notes)

	return formatted
}
