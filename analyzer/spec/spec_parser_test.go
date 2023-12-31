package spec

import (
	"os"
	"reflect"
	"testing"

	"github.com/ConfigMate/configmate/parsers"
)

// TestParseSimple tests the parser's ability to parse the simple.cms file
func TestParseSimple(t *testing.T) {
	simpleCMS, err := os.ReadFile("./test_specs/simple.cms")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	expectedSpec := &Specification{
		File: "./examples/configurations/config0.json",
		FileLocation: parsers.TokenLocation{
			Start: parsers.CharLocation{Line: 0, Column: 8},
			End:   parsers.CharLocation{Line: 0, Column: 48},
		},
		FileFormat: "json",
		FileFormatLocation: parsers.TokenLocation{
			Start: parsers.CharLocation{Line: 0, Column: 49},
			End:   parsers.CharLocation{Line: 0, Column: 53},
		},
		Imports:              map[string]string{},
		ImportsAliasLocation: map[string]parsers.TokenLocation{},
		ImportsLocation:      map[string]parsers.TokenLocation{},
		Fields: []FieldSpec{
			{
				Field: &parsers.NodeKey{Segments: []string{"server"}},
				FieldLocation: parsers.TokenLocation{
					Start: parsers.CharLocation{Line: 3, Column: 4},
					End:   parsers.CharLocation{Line: 3, Column: 10},
				},
				Type: "object",
				TypeLocation: parsers.TokenLocation{
					Start: parsers.CharLocation{Line: 3, Column: 11},
					End:   parsers.CharLocation{Line: 3, Column: 17},
				},
				Optional: true,
				OptionalLocation: parsers.TokenLocation{
					Start: parsers.CharLocation{Line: 3, Column: 19},
					End:   parsers.CharLocation{Line: 3, Column: 27},
				},
				Checks: []CheckWithLocation{},
			},
			{
				Field: &parsers.NodeKey{Segments: []string{"server", "host"}},
				FieldLocation: parsers.TokenLocation{
					Start: parsers.CharLocation{Line: 4, Column: 8},
					End:   parsers.CharLocation{Line: 4, Column: 12},
				},
				Type: "string",
				TypeLocation: parsers.TokenLocation{
					Start: parsers.CharLocation{Line: 5, Column: 18},
					End:   parsers.CharLocation{Line: 5, Column: 24},
				},
				Default: "localhost",
				DefaultLocation: parsers.TokenLocation{
					Start: parsers.CharLocation{Line: 6, Column: 21},
					End:   parsers.CharLocation{Line: 6, Column: 32},
				},
				Notes: "This is the host that the server will listen on.",
				NotesLocation: parsers.TokenLocation{
					Start: parsers.CharLocation{Line: 7, Column: 19},
					End:   parsers.CharLocation{Line: 7, Column: 69},
				},
				Checks: []CheckWithLocation{
					{
						Check: "eq(\"localhost\")",
						Location: parsers.TokenLocation{
							Start: parsers.CharLocation{Line: 8, Column: 12},
							End:   parsers.CharLocation{Line: 8, Column: 27},
						},
					},
				},
			},
			{
				Field: &parsers.NodeKey{Segments: []string{"server", "port"}},
				FieldLocation: parsers.TokenLocation{
					Start: parsers.CharLocation{Line: 10, Column: 8},
					End:   parsers.CharLocation{Line: 10, Column: 12},
				},
				Type: "int",
				TypeLocation: parsers.TokenLocation{
					Start: parsers.CharLocation{Line: 11, Column: 18},
					End:   parsers.CharLocation{Line: 11, Column: 21},
				},
				Default: "80",
				DefaultLocation: parsers.TokenLocation{
					Start: parsers.CharLocation{Line: 12, Column: 21},
					End:   parsers.CharLocation{Line: 12, Column: 23},
				},
				Notes: "This is the port that the server will listen on. We are also testing multiline strings here.",
				NotesLocation: parsers.TokenLocation{
					Start: parsers.CharLocation{Line: 13, Column: 19},
					End:   parsers.CharLocation{Line: 16, Column: 15},
				},
				Checks: []CheckWithLocation{
					{
						Check: "range(25,100)",
						Location: parsers.TokenLocation{
							Start: parsers.CharLocation{Line: 17, Column: 12},
							End:   parsers.CharLocation{Line: 17, Column: 26},
						},
					},
				},
			},
			{
				Field: &parsers.NodeKey{Segments: []string{"server", "ssl_enabled"}},
				FieldLocation: parsers.TokenLocation{
					Start: parsers.CharLocation{Line: 19, Column: 8},
					End:   parsers.CharLocation{Line: 19, Column: 19},
				},
				Type: "bool",
				TypeLocation: parsers.TokenLocation{
					Start: parsers.CharLocation{Line: 20, Column: 18},
					End:   parsers.CharLocation{Line: 20, Column: 22},
				},
				Default: "false",
				DefaultLocation: parsers.TokenLocation{
					Start: parsers.CharLocation{Line: 21, Column: 21},
					End:   parsers.CharLocation{Line: 21, Column: 26},
				},
				Notes: "This is whether or not SSL is enabled.",
				NotesLocation: parsers.TokenLocation{
					Start: parsers.CharLocation{Line: 22, Column: 19},
					End:   parsers.CharLocation{Line: 22, Column: 59},
				},
				Checks: []CheckWithLocation{
					{
						Check: "eq(false)",
						Location: parsers.TokenLocation{
							Start: parsers.CharLocation{Line: 23, Column: 12},
							End:   parsers.CharLocation{Line: 23, Column: 21},
						},
					},
				},
			},
			{
				Field: &parsers.NodeKey{Segments: []string{"server", "dns. servers"}},
				FieldLocation: parsers.TokenLocation{
					Start: parsers.CharLocation{Line: 25, Column: 8},
					End:   parsers.CharLocation{Line: 25, Column: 22},
				},
				Type: "list<string>",
				TypeLocation: parsers.TokenLocation{
					Start: parsers.CharLocation{Line: 26, Column: 18},
					End:   parsers.CharLocation{Line: 26, Column: 30},
				},
				Optional: true,
				OptionalLocation: parsers.TokenLocation{
					Start: parsers.CharLocation{Line: 27, Column: 22},
					End:   parsers.CharLocation{Line: 27, Column: 26},
				},
				Notes: "This is a list of DNS servers.",
				NotesLocation: parsers.TokenLocation{
					Start: parsers.CharLocation{Line: 28, Column: 19},
					End:   parsers.CharLocation{Line: 28, Column: 51},
				},
				Checks: []CheckWithLocation{
					{
						Check: "len().gte(3)",
						Location: parsers.TokenLocation{
							Start: parsers.CharLocation{Line: 29, Column: 12},
							End:   parsers.CharLocation{Line: 29, Column: 24},
						},
					},
				},
			},
			{
				Field: &parsers.NodeKey{Segments: []string{"server", "apis"}},
				FieldLocation: parsers.TokenLocation{
					Start: parsers.CharLocation{Line: 31, Column: 8},
					End:   parsers.CharLocation{Line: 31, Column: 12},
				},
				Type: "list<api_info>",
				TypeLocation: parsers.TokenLocation{
					Start: parsers.CharLocation{Line: 31, Column: 13},
					End:   parsers.CharLocation{Line: 31, Column: 27},
				},
				Checks: []CheckWithLocation{},
			},
		},
		Objects: []ObjectDef{
			{
				Name: "api_info",
				NameLocation: parsers.TokenLocation{
					Start: parsers.CharLocation{Line: 36, Column: 4},
					End:   parsers.CharLocation{Line: 36, Column: 12},
				},
				Properties: []ObjectPropertyDef{
					{
						Name: "endpoint",
						NameLocation: parsers.TokenLocation{
							Start: parsers.CharLocation{Line: 37, Column: 8},
							End:   parsers.CharLocation{Line: 37, Column: 16},
						},
						Type: "string",
						TypeLocation: parsers.TokenLocation{
							Start: parsers.CharLocation{Line: 37, Column: 17},
							End:   parsers.CharLocation{Line: 37, Column: 23},
						},
					},
					{
						Name: "timeout",
						NameLocation: parsers.TokenLocation{
							Start: parsers.CharLocation{Line: 38, Column: 8},
							End:   parsers.CharLocation{Line: 38, Column: 15},
						},
						Type: "int",
						TypeLocation: parsers.TokenLocation{
							Start: parsers.CharLocation{Line: 38, Column: 16},
							End:   parsers.CharLocation{Line: 38, Column: 19},
						},
						Optional: true,
						OptionalLocation: parsers.TokenLocation{
							Start: parsers.CharLocation{Line: 38, Column: 21},
							End:   parsers.CharLocation{Line: 38, Column: 29},
						},
					},
					{
						Name: "method",
						NameLocation: parsers.TokenLocation{
							Start: parsers.CharLocation{Line: 39, Column: 8},
							End:   parsers.CharLocation{Line: 39, Column: 14},
						},
						Type: "string",
						TypeLocation: parsers.TokenLocation{
							Start: parsers.CharLocation{Line: 39, Column: 15},
							End:   parsers.CharLocation{Line: 39, Column: 21},
						},
					},
				},
			},
		},
	}

	parser := NewSpecParser()
	result, errs := parser.Parse(simpleCMS)
	if len(errs) > 0 {
		t.Errorf("Unexpected errors: %#v", errs)
	}
	if !reflect.DeepEqual(result, expectedSpec) {
		t.Errorf("Expected: %#v\nGot: %#v", expectedSpec, result)
	}
}

// TestParseWithImports tests the parser's ability to parse imort statements
func TestParseWithImports(t *testing.T) {
	withImportsCMS, err := os.ReadFile("./test_specs/with_imports.cms")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	expectedSpec := &Specification{
		File: "./some/file.json",
		FileLocation: parsers.TokenLocation{
			Start: parsers.CharLocation{Line: 0, Column: 8},
			End:   parsers.CharLocation{Line: 0, Column: 26},
		},
		FileFormat: "json",
		FileFormatLocation: parsers.TokenLocation{
			Start: parsers.CharLocation{Line: 0, Column: 27},
			End:   parsers.CharLocation{Line: 0, Column: 31},
		},
		Imports: map[string]string{
			"someSpec":      "./specs/someSpec.cms",
			"someOtherSpec": "./specs/someOtherSpec.cms",
		},
		ImportsAliasLocation: map[string]parsers.TokenLocation{
			"someSpec": {
				Start: parsers.CharLocation{Line: 3, Column: 4},
				End:   parsers.CharLocation{Line: 3, Column: 12},
			},
			"someOtherSpec": {
				Start: parsers.CharLocation{Line: 4, Column: 4},
				End:   parsers.CharLocation{Line: 4, Column: 17},
			},
		},
		ImportsLocation: map[string]parsers.TokenLocation{
			"someSpec": {
				Start: parsers.CharLocation{Line: 3, Column: 14},
				End:   parsers.CharLocation{Line: 3, Column: 36},
			},
			"someOtherSpec": {
				Start: parsers.CharLocation{Line: 4, Column: 19},
				End:   parsers.CharLocation{Line: 4, Column: 46},
			},
		},
		Fields: []FieldSpec{},
	}

	parser := NewSpecParser()
	result, errs := parser.Parse(withImportsCMS)
	if len(errs) > 0 {
		t.Errorf("Unexpected errors: %#v", errs)
	}
	if !reflect.DeepEqual(result, expectedSpec) {
		t.Errorf("Expected: %#v\nGot: %#v", expectedSpec, result)
	}
}

// TestParserHighLevelErrors tests the parser's ability to report high level errors
func TestParserHighLevelErrors(t *testing.T) {
	cmsWithHighLevelErrors, err := os.ReadFile("./test_specs/with_highlevel_errors.cms")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	expectedErrors := []SpecParserError{
		{
			ErrorMessage: "duplicate default metadata for field server.host",
			Location: parsers.TokenLocation{
				Start: parsers.CharLocation{Line: 7, Column: 12},
				End:   parsers.CharLocation{Line: 7, Column: 32},
			},
		},
		{
			ErrorMessage: "duplicate notes metadata for field server.port",
			Location: parsers.TokenLocation{
				Start: parsers.CharLocation{Line: 17, Column: 12},
				End:   parsers.CharLocation{Line: 17, Column: 69},
			},
		},
		{
			ErrorMessage: "missing type metadata for field server.port",
			Location: parsers.TokenLocation{
				Start: parsers.CharLocation{Line: 11, Column: 8},
				End:   parsers.CharLocation{Line: 18, Column: 29},
			},
		},
		{
			ErrorMessage: "duplicate type metadata for field server.ssl_enabled",
			Location: parsers.TokenLocation{
				Start: parsers.CharLocation{Line: 22, Column: 12},
				End:   parsers.CharLocation{Line: 22, Column: 22},
			},
		},
		{
			ErrorMessage: "duplicate optional metadata for field server.dns_servers",
			Location: parsers.TokenLocation{
				Start: parsers.CharLocation{Line: 30, Column: 12},
				End:   parsers.CharLocation{Line: 30, Column: 26},
			},
		},
	}

	parser := NewSpecParser()
	_, errs := parser.Parse(cmsWithHighLevelErrors)
	if len(errs) == 0 {
		t.Errorf("Expecting errors, no errors where returned instead")
	} else if !reflect.DeepEqual(errs, expectedErrors) {
		t.Errorf("Expected: %#v\nGot: %#v", expectedErrors, errs)
	}
}

// TestParserLexerSyntaxErrors tests the parser's ability to report errors in the lexer stage
func TestParserLexerSyntaxErrors(t *testing.T) {
	cmsWithLexerErrors, err := os.ReadFile("./test_specs/with_lexer_errors.cms")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	expectedErrors := []SpecParserError{
		{
			ErrorMessage: "mismatched input '\"./examples/configurations/config0.json' json\\n\\nspec {\\n    server <type: object> {\\n        host <\\n            type: string,\\n            default: \"' expecting SHORT_STRING",
			Location: parsers.TokenLocation{
				Start: parsers.CharLocation{Line: 0, Column: 8},
				End:   parsers.CharLocation{Line: 0, Column: 9},
			},
		},
		{
			ErrorMessage: "token recognition error at: '\"\n        > ( len().gte(3-5); )\n    }\n}\n'",
			Location: parsers.TokenLocation{
				Start: parsers.CharLocation{Line: 28, Column: 50},
				End:   parsers.CharLocation{Line: 28, Column: 51},
			},
		},
	}

	parser := NewSpecParser()
	_, errs := parser.Parse(cmsWithLexerErrors)
	if len(errs) == 0 {
		t.Errorf("Expecting errors, no errors where returned instead")
	} else if !reflect.DeepEqual(errs, expectedErrors) {
		t.Errorf("Expected: %#v\nGot: %#v", expectedErrors, errs)
	}
}

// TestParserSyntaxErrors tests the parser's ability to report syntax errors
func TestParserSyntaxErrors(t *testing.T) {
	cmsWithParserErrors, err := os.ReadFile("./test_specs/with_parser_errors.cms")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	expectedErrors := []SpecParserError{
		{
			ErrorMessage: "extraneous input ':' expecting IDENTIFIER",
			Location: parsers.TokenLocation{
				Start: parsers.CharLocation{Line: 0, Column: 49},
				End:   parsers.CharLocation{Line: 0, Column: 50},
			},
		},
		{
			ErrorMessage: "missing ';' at ')'",
			Location: parsers.TokenLocation{
				Start: parsers.CharLocation{Line: 17, Column: 27},
				End:   parsers.CharLocation{Line: 17, Column: 28},
			},
		},
		{
			ErrorMessage: "missing '>' at ','",
			Location: parsers.TokenLocation{
				Start: parsers.CharLocation{Line: 26, Column: 29},
				End:   parsers.CharLocation{Line: 26, Column: 30},
			},
		},
		{
			ErrorMessage: "mismatched input '\"true\"' expecting BOOL",
			Location: parsers.TokenLocation{
				Start: parsers.CharLocation{Line: 27, Column: 22},
				End:   parsers.CharLocation{Line: 27, Column: 23},
			},
		},
	}

	parser := NewSpecParser()
	_, errs := parser.Parse(cmsWithParserErrors)
	if len(errs) == 0 {
		t.Errorf("Expecting errors, no errors where returned instead")
	} else if !reflect.DeepEqual(errs, expectedErrors) {
		t.Errorf("Expected: %#v\nGot: %#v", expectedErrors, errs)
	}
}
