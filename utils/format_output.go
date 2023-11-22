package utils

import (
	"fmt"
	"strings"

	"github.com/ConfigMate/configmate/analyzer"
)

const linesPaddingForErrors = 2
const rightPaddingForMultilineErrorArrows = 2

func FormatCheckResult(res analyzer.CheckResult, fileLinesMap map[string]map[int]string) string {
	var status, comment, check, fieldType, optional string

	// Format of the result
	format := "%s:\n\tField: %s %s %s \n\tCheck: %s %s \n"

	// format values
	if res.Status == analyzer.CheckPassed {
		status = ColorText("PASSED", Green)
	} else if res.Status == analyzer.CheckSkipped {
		status = ColorText("SKIPPED", Yellow)
		comment = ColorText(res.ResultComment, Yellow)
		comment = fmt.Sprintf("- %s", comment)
	} else {
		status = ColorText("FAILED", Red)
		comment = ColorText(res.ResultComment, Red)
		comment = fmt.Sprintf("- %s", comment)
	}

	// get other values
	check = ColorText(res.Field.Checks[res.CheckNum].Check, Cyan)
	fieldType = ColorText(res.Field.Type, Blue)
	if res.Field.Optional {
		optional = ColorText("optional", Yellow)
	} else {
		optional = ColorText("required", Gray)
	}

	// Format the values
	formatted := fmt.Sprintf(format, status, res.Field.Field, fieldType, optional, check, comment)

	if res.Field.Default != "" {
		formatted = fmt.Sprintf("%s\tDefault: %v\n", formatted, res.Field.Default)
	}

	if res.Field.Notes != "" {
		formatted = fmt.Sprintf("%s\tNotes: %s\n", formatted, res.Field.Notes)
	}

	// If check failed, print the problematic line
	if res.Status == analyzer.CheckFailed {
		for _, token := range res.TokenList {
			// Create token view
			tokenView := createTokenErrorView(token, fileLinesMap)
			// Add tokenView to formatted
			formatted = fmt.Sprintf("%s%s", formatted, tokenView)
		}
	}

	formatted = fmt.Sprintf("%s\n", formatted) // Add extra new line for readability

	return formatted
}

func FormatSpecError(specError analyzer.SpecError, fileLinesMap map[string]map[int]string) string {
	// Specification Error header
	header := ColorText("Specification Error:", Red)

	// Format of the result
	format := "%s:\n\tAnalyzer Message: %s\n"

	// Format the values
	formatted := fmt.Sprintf(format, header, specError.AnalyzerMsg)

	// Add more info if error message available
	if len(specError.ErrorMsgs) > 0 {
		formatted = fmt.Sprintf("%s\tErrors:", formatted)
		for _, errorMsg := range specError.ErrorMsgs {
			formattedErrorMessage := "- " + errorMsg
			formattedErrorMessage = strings.ReplaceAll(formattedErrorMessage, "\n", "\\n")
			formatted = fmt.Sprintf("%s\n\t\t%s", formatted, formattedErrorMessage)
		}
		formatted = fmt.Sprintf("%s\n", formatted) // Add extra new line for readability
	}

	// Add token view for each token
	for _, token := range specError.TokenList {
		// Create token view
		tokenView := createTokenErrorView(token, fileLinesMap)
		// Add tokenView to formatted
		formatted = fmt.Sprintf("%s\n%s", formatted, tokenView)
	}

	formatted = fmt.Sprintf("%s\n", formatted) // Add extra new line for readability

	return formatted
}

func createTokenErrorView(token analyzer.TokenLocationWithFile, fileLinesMap map[string]map[int]string) string {
	// Get the line number
	startLineNum := token.Location.Start.Line
	startColNum := token.Location.Start.Column
	endLineNum := token.Location.End.Line
	endColNum := token.Location.End.Column

	// Add the file path to the output
	fileLine := fmt.Sprintf("\tFile: %s\n\n", ColorText(fmt.Sprintf("%s:%d", token.File, startLineNum+1), Red))

	var output string
	// Case where the token is in one line
	if startLineNum == endLineNum {
		// Get the line
		line := fileLinesMap[token.File][startLineNum]

		// Get the content of the line up to the token
		preTokenContent := line[:startColNum]

		// Get the content of the token and color it red
		tokenContent := ColorText(line[startColNum:endColNum], Red)
		tokenLength := endColNum - startColNum

		// Get the content of the line after the token
		postTokenContent := line[endColNum:]

		// Format the line
		output = fmt.Sprintf("\t  Line %d: %s%s%s\n", startLineNum, preTokenContent, tokenContent, postTokenContent)

		// Create string of empty spaces the size of that startLineNum would take
		lineNumberSpace := ""
		for i := 0; i < len(fmt.Sprintf("%d", startLineNum)); i++ {
			lineNumberSpace = fmt.Sprintf("%s ", lineNumberSpace)
		}

		// Add arrows below the token
		output = fmt.Sprintf("%s\t         %s", output, lineNumberSpace)
		for i := 0; i < startColNum; i++ {
			output = fmt.Sprintf("%s ", output)
		}
		for i := 0; i < tokenLength; i++ {
			output = fmt.Sprintf("%s%s", output, ColorText("^", Red))
		}
		output = fmt.Sprintf("%s\n", output)
	} else { // Case where the token is in multiple lines
		// Find the max length of the lines
		maxLength := 0
		for i := startLineNum; i <= endLineNum; i++ {
			length := len(strings.TrimRight(fileLinesMap[token.File][i], " "))
			if length > maxLength {
				maxLength = length
			}
		}

		// Get the start line
		startLine := fileLinesMap[token.File][startLineNum]

		// Get the content of the start line up to the token
		preTokenContent := startLine[:startColNum]

		// Get the rest of the line and color it red
		startLineTokenContent := ColorText(startLine[startColNum:], Red)

		// Find offset for the start line
		offset := maxLength - len(strings.TrimRight(startLine, " "))

		// Create arrowOffset of that offset plus the padding
		arrowOffset := ""
		for i := 0; i < offset+rightPaddingForMultilineErrorArrows; i++ {
			arrowOffset = fmt.Sprintf("%s ", arrowOffset)
		}

		// Format the start line
		output = fmt.Sprintf("\t  Line %d: %s%s%s%s\n", startLineNum, preTokenContent, startLineTokenContent, arrowOffset, ColorText("<", Red))

		// Get the middle lines
		for i := startLineNum + 1; i < endLineNum; i++ {
			// Get the line
			line := fileLinesMap[token.File][i]

			// Find offset for this line
			offset = maxLength - len(strings.TrimRight(line, " "))

			// Create arrowOffset of that offset plus the padding
			arrowOffset = ""
			for i := 0; i < offset+rightPaddingForMultilineErrorArrows; i++ {
				arrowOffset = fmt.Sprintf("%s ", arrowOffset)
			}

			// Format the line
			output = fmt.Sprintf("%s\t  Line %d: %s%s%s\n", output, i, ColorText(line, Red), arrowOffset, ColorText("<", Red))
		}

		// Get the end line
		endLine := fileLinesMap[token.File][endLineNum]

		// Get the content of the end line up to the end of the token
		endLineTokenContent := ColorText(endLine[:endColNum], Red)

		// Get the rest of the line
		postTokenContent := endLine[endColNum:]

		// Find offset for the end line
		offset = maxLength - len(strings.TrimRight(endLine, " "))

		// Create arrowOffset of that offset plus the padding
		arrowOffset = ""
		for i := 0; i < offset+rightPaddingForMultilineErrorArrows; i++ {
			arrowOffset = fmt.Sprintf("%s ", arrowOffset)
		}

		// Format the end line
		output = fmt.Sprintf("%s\t  Line %d: %s%s%s%s\n", output, endLineNum, endLineTokenContent, postTokenContent, arrowOffset, ColorText("<", Red))
	}

	// Add padding top
	for i := 1; i <= linesPaddingForErrors; i++ {
		// Get line number
		lineNum := startLineNum - i

		// Check if line exists
		if line, ok := fileLinesMap[token.File][lineNum]; ok {
			// Format the line
			output = fmt.Sprintf("\t  Line %d: %s\n%s", lineNum, line, output)
		}
	}

	// Add padding bottom
	for i := 1; i <= linesPaddingForErrors; i++ {
		// Get line number
		lineNum := endLineNum + i

		// Check if line exists
		if line, ok := fileLinesMap[token.File][lineNum]; ok {
			// Format the line
			output = fmt.Sprintf("%s\t  Line %d: %s\n", output, lineNum, line)
		}
	}

	return fmt.Sprintf("%s%s", fileLine, output)
}
