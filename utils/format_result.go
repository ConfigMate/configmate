package utils

import (
	"fmt"

	"github.com/ConfigMate/configmate/analyzer"
)

const linesPaddingForErrors = 2

func FormatResult(res analyzer.Result, files map[string]analyzer.FileDetails, fileLinesMap map[string]map[int]string) string {
	var passed, comment, check, fieldType, optional string

	// Format of the result
	format := "%s:\n\tField: %s %s %s \n\tCheck: %s %s \n"

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
		optional = ColorText("required", Gray)
	}

	// Format the values
	formatted := fmt.Sprintf(format, passed, res.Rule.Field, fieldType, optional, check, comment)

	if res.Rule.Default != nil {
		formatted = fmt.Sprintf("%s\tDefault: %v\n", formatted, res.Rule.Default)
	}

	if res.Rule.Notes != "" {
		formatted = fmt.Sprintf("%s\tNotes: %s\n", formatted, res.Rule.Notes)
	}

	// If check failed, print the problematic line
	if !res.Passed {
		for _, token := range res.TokenList {
			// Get the file alias
			fileAlias := token.File

			// Add the file alias to the formatted string
			formatted = fmt.Sprintf("%s\tFile: %s\n\n", formatted, ColorText(files[fileAlias].Path, Red))

			// Get the line number
			startLineNum := token.Location.Start.Line
			startColNum := token.Location.Start.Column
			endLineNum := token.Location.End.Line
			endColNum := token.Location.End.Column

			var output string
			// Case where the token is in one line
			if startLineNum == endLineNum {
				// Get the line
				line := fileLinesMap[fileAlias][startLineNum]

				// Get the content of the line up to the token
				preTokenContent := line[:startColNum-1]

				// Get the content of the token and color it red
				tokenContent := ColorText(line[startColNum-1:endColNum], Red)
				tokenLength := endColNum - startColNum + 1

				// Get the content of the line after the token
				postTokenContent := line[endColNum:]

				// Format the line
				output = fmt.Sprintf("\t  Line %d: %s%s%s\n", startLineNum, preTokenContent, tokenContent, postTokenContent)
				// Add arrows below the token
				output = fmt.Sprintf("%s\t          ", output)
				for i := 0; i < startColNum-1; i++ {
					output = fmt.Sprintf("%s ", output)
				}
				for i := 0; i < tokenLength; i++ {
					output = fmt.Sprintf("%s%s", output, ColorText("^", Red))
				}
				output = fmt.Sprintf("%s\n", output)
			} else { // Case where the token is in multiple lines
				// Get the start line
				startLine := fileLinesMap[fileAlias][startLineNum]

				// Get the content of the start line up to the token
				preTokenContent := startLine[:startColNum-1]

				// Get the rest of the line and color it red
				startLineTokenContent := ColorText(startLine[startColNum-1:], Red)

				// Format the start line
				output = fmt.Sprintf("\t  Line %d: %s%s\n", startLineNum, preTokenContent, startLineTokenContent)

				// Get the middle lines
				for i := startLineNum + 1; i < endLineNum; i++ {
					// Get the line
					line := fileLinesMap[fileAlias][i]

					// Format the line
					output = fmt.Sprintf("%s\t  Line %d: %s\n", output, i, ColorText(line, Red))
				}

				// Get the end line
				endLine := fileLinesMap[fileAlias][endLineNum]

				// Get the content of the end line up to the end of the token
				endLineTokenContent := ColorText(endLine[:endColNum], Red)

				// Get the rest of the line
				postTokenContent := endLine[endColNum:]

				// Format the end line
				output = fmt.Sprintf("%s\t  Line %d: %s%s\n", output, endLineNum, endLineTokenContent, postTokenContent)
			}

			// Add padding top
			for i := 1; i <= linesPaddingForErrors; i++ {
				// Get line number
				lineNum := startLineNum - i

				// Check if line exists
				if line, ok := fileLinesMap[fileAlias][lineNum]; ok {
					// Format the line
					output = fmt.Sprintf("\t  Line %d: %s\n%s", lineNum, line, output)
				}
			}

			// Add padding bottom
			for i := 1; i <= linesPaddingForErrors; i++ {
				// Get line number
				lineNum := endLineNum + i

				// Check if line exists
				if line, ok := fileLinesMap[fileAlias][lineNum]; ok {
					// Format the line
					output = fmt.Sprintf("%s\t  Line %d: %s\n", output, lineNum, line)
				}
			}

			// Add output to formatted
			formatted = fmt.Sprintf("%s%s\n", formatted, output)
		}
	}

	return formatted
}
