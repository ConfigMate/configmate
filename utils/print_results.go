package utils

import (
	"fmt"
	"github.com/ConfigMate/configmate/analyzer"
)

func FormatResult(res analyzer.Result) string {
	var passed, comment, check, fieldType, optional string

	// Format of the result
	format := "%s:\t%s \n\tCheck:\t%s \n\tField:\t%s %s %s \n\t\t%s"

	// format values
	if(res.Passed) {
		passed = ColorText("PASSED", Green)
		comment = ColorText(res.ResultComment, Green)
	} else {
		passed = ColorText("FAILED", Red)
		comment = ColorText(res.ResultComment, Red)
	}
	
	// get other values
	check = res.Rule.Checks[res.CheckNumber]
	fieldType = ColorText(res.Rule.Type, Blue)
	if(res.Rule.Optional) {
		optional = ColorText("optional", Yellow)
	} else {
		optional = ColorText("required", Red)
	}

	// Format the values
	formatted := fmt.Sprintf(format, passed, comment, check, res.Rule.Field, fieldType, optional, res.Rule.Notes)

	return formatted
}