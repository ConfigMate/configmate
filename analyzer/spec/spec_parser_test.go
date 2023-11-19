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
			Start: parsers.CharLocation{Line: 1, Column: 6},
			End:   parsers.CharLocation{Line: 1, Column: 46},
		},
		FileFormat: "json",
		FileFormatLocation: parsers.TokenLocation{
			Start: parsers.CharLocation{Line: 1, Column: 47},
			End:   parsers.CharLocation{Line: 1, Column: 51},
		},
		Imports:         map[string]string{},
		ImportsLocation: map[string]parsers.TokenLocation{},
		Fields: []FieldSpec{
			{
				Field: "server",
				FieldLocation: parsers.TokenLocation{
					Start: parsers.CharLocation{Line: 4, Column: 4},
					End:   parsers.CharLocation{Line: 4, Column: 10},
				},
				Type: "object",
				TypeLocation: parsers.TokenLocation{
					Start: parsers.CharLocation{Line: 4, Column: 18},
					End:   parsers.CharLocation{Line: 4, Column: 24},
				},
				Optional:         false,
				OptionalLocation: parsers.TokenLocation{},
				Default:          "",
				DefaultLocation:  parsers.TokenLocation{},
				Notes:            "",
				NotesLocation:    parsers.TokenLocation{},
				Checks:           []CheckWithLocation{},
			},
			{
				Field: "server.host",
				FieldLocation: parsers.TokenLocation{
					Start: parsers.CharLocation{Line: 5, Column: 8},
					End:   parsers.CharLocation{Line: 5, Column: 12},
				},
				Type: "string",
				TypeLocation: parsers.TokenLocation{
					Start: parsers.CharLocation{Line: 6, Column: 18},
					End:   parsers.CharLocation{Line: 6, Column: 24},
				},
				Optional:         false,
				OptionalLocation: parsers.TokenLocation{},
				Default:          "localhost",
				DefaultLocation: parsers.TokenLocation{
					Start: parsers.CharLocation{Line: 7, Column: 21},
					End:   parsers.CharLocation{Line: 7, Column: 32},
				},
				Notes: "This is the host that the server will listen on.",
				NotesLocation: parsers.TokenLocation{
					Start: parsers.CharLocation{Line: 8, Column: 19},
					End:   parsers.CharLocation{Line: 8, Column: 69},
				},
				Checks: []CheckWithLocation{
					{
						Check: "eq(\"localhost\")",
						Location: parsers.TokenLocation{
							Start: parsers.CharLocation{Line: 9, Column: 12},
							End:   parsers.CharLocation{Line: 9, Column: 27},
						},
					},
				},
			},
			{
				Field: "server.port",
				FieldLocation: parsers.TokenLocation{
					Start: parsers.CharLocation{Line: 11, Column: 8},
					End:   parsers.CharLocation{Line: 11, Column: 12},
				},
				Type: "int",
				TypeLocation: parsers.TokenLocation{
					Start: parsers.CharLocation{Line: 12, Column: 18},
					End:   parsers.CharLocation{Line: 12, Column: 21},
				},
				Optional:         false,
				OptionalLocation: parsers.TokenLocation{},
				Default:          "80",
				DefaultLocation: parsers.TokenLocation{
					Start: parsers.CharLocation{Line: 13, Column: 21},
					End:   parsers.CharLocation{Line: 13, Column: 23},
				},
				Notes: "This is the port that the server will listen on. We are also testing multiline strings here.",
				NotesLocation: parsers.TokenLocation{
					Start: parsers.CharLocation{Line: 14, Column: 19},
					End:   parsers.CharLocation{Line: 17, Column: 15},
				},
				Checks: []CheckWithLocation{
					{
						Check: "range(25,100)",
						Location: parsers.TokenLocation{
							Start: parsers.CharLocation{Line: 18, Column: 12},
							End:   parsers.CharLocation{Line: 18, Column: 26},
						},
					},
				},
			},
			{
				Field: "server.ssl_enabled",
				FieldLocation: parsers.TokenLocation{
					Start: parsers.CharLocation{Line: 20, Column: 8},
					End:   parsers.CharLocation{Line: 20, Column: 19},
				},
				Type: "bool",
				TypeLocation: parsers.TokenLocation{
					Start: parsers.CharLocation{Line: 21, Column: 18},
					End:   parsers.CharLocation{Line: 21, Column: 22},
				},
				Optional:         false,
				OptionalLocation: parsers.TokenLocation{},
				Default:          "false",
				DefaultLocation: parsers.TokenLocation{
					Start: parsers.CharLocation{Line: 22, Column: 21},
					End:   parsers.CharLocation{Line: 22, Column: 26},
				},
				Notes: "This is whether or not SSL is enabled.",
				NotesLocation: parsers.TokenLocation{
					Start: parsers.CharLocation{Line: 23, Column: 19},
					End:   parsers.CharLocation{Line: 23, Column: 59},
				},
				Checks: []CheckWithLocation{
					{
						Check: "eq(false)",
						Location: parsers.TokenLocation{
							Start: parsers.CharLocation{Line: 24, Column: 12},
							End:   parsers.CharLocation{Line: 24, Column: 21},
						},
					},
				},
			},
			{
				Field: "server.dns_servers",
				FieldLocation: parsers.TokenLocation{
					Start: parsers.CharLocation{Line: 26, Column: 8},
					End:   parsers.CharLocation{Line: 26, Column: 19},
				},
				Type: "list:string",
				TypeLocation: parsers.TokenLocation{
					Start: parsers.CharLocation{Line: 27, Column: 18},
					End:   parsers.CharLocation{Line: 27, Column: 30},
				},
				Optional: true,
				OptionalLocation: parsers.TokenLocation{
					Start: parsers.CharLocation{Line: 28, Column: 22},
					End:   parsers.CharLocation{Line: 28, Column: 26},
				},
				Default:         "",
				DefaultLocation: parsers.TokenLocation{},
				Notes:           "This is a list of DNS servers.",
				NotesLocation: parsers.TokenLocation{
					Start: parsers.CharLocation{Line: 29, Column: 19},
					End:   parsers.CharLocation{Line: 29, Column: 51},
				},
				Checks: []CheckWithLocation{
					{
						Check: "len().gte(3)",
						Location: parsers.TokenLocation{
							Start: parsers.CharLocation{Line: 30, Column: 12},
							End:   parsers.CharLocation{Line: 30, Column: 24},
						},
					},
				},
			},
		},
	}

	parser := NewSpecParser()
	result, err := parser.Parse(string(simpleCMS))
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if !reflect.DeepEqual(result, expectedSpec) {
		t.Errorf("Expected: %#v\nGot: %#v", expectedSpec, result)
	}
}
