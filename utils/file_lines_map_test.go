package utils

import (
	"reflect"
	"testing"
)

// TestCreateLineMap tests the CreateLineMap function.
func TestCreateLineMap(t *testing.T) {
	// Test cases
	type testCase struct {
		input    []byte
		expected map[int]string
	}

	// Create test inputs
	tomlFile := []byte(`# This is a TOML configuration for the Rulebook
name = "Sample Rulebook"
description = "This is a sample rulebook for experimentation."

# Files to be checked
[files.file1]
path = "/examples/configurations/sample_config.json"
format = "json"`)

	jsonFile := []byte(`{
	"name": "Sample Rulebook",
	"description": "This is a sample rulebook for experimentation.",
	"files": {
		"file1": {
			"path": "/examples/configurations/sample_config.json",
			"format": "json"
		}
	}
}`)

	hoconFile := []byte(`name = "Sample Rulebook"
description = "This is a sample rulebook for experimentation."

# Files to be checked
files {
	file1 {
		path = "/examples/configurations/sample_config.json"
		format = "json"
	}
}`)

	testCases := []testCase{
		{ // TOML file
			input: tomlFile,
			expected: map[int]string{
				1: "# This is a TOML configuration for the Rulebook",
				2: "name = \"Sample Rulebook\"",
				3: "description = \"This is a sample rulebook for experimentation.\"",
				4: "",
				5: "# Files to be checked",
				6: "[files.file1]",
				7: "path = \"/examples/configurations/sample_config.json\"",
				8: "format = \"json\"",
			},
		},
		{ // JSON file
			input: jsonFile,
			expected: map[int]string{
				1:  "{",
				2:  "\t\"name\": \"Sample Rulebook\",",
				3:  "\t\"description\": \"This is a sample rulebook for experimentation.\",",
				4:  "\t\"files\": {",
				5:  "\t\t\"file1\": {",
				6:  "\t\t\t\"path\": \"/examples/configurations/sample_config.json\",",
				7:  "\t\t\t\"format\": \"json\"",
				8:  "\t\t}",
				9:  "\t}",
				10: "}",
			},
		},
		{ // HOCON file
			input: hoconFile,
			expected: map[int]string{
				1:  "name = \"Sample Rulebook\"",
				2:  "description = \"This is a sample rulebook for experimentation.\"",
				3:  "",
				4:  "# Files to be checked",
				5:  "files {",
				6:  "\tfile1 {",
				7:  "\t\tpath = \"/examples/configurations/sample_config.json\"",
				8:  "\t\tformat = \"json\"",
				9:  "\t}",
				10: "}",
			},
		},
		{ // Empty file
			input:    []byte{},
			expected: map[int]string{},
		},
	}

	// Run tests
	for _, test := range testCases {
		actual := createLineMap(test.input)

		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("Expected %v, got %v", test.expected, actual)
		}
	}
}

// CreateLinesMapForFiles tests the CreateLinesMapForFiles function.
func TestCreateLinesMapForFiles(t *testing.T) {
	// Test cases
	type testCase struct {
		input    map[string][]byte
		expected map[string]map[int]string
	}

	// Create test inputs
	tomlFile := []byte(`# This is a TOML configuration for the Rulebook
name = "Sample Rulebook"
description = "This is a sample rulebook for experimentation."

# Files to be checked
[files.file1]
path = "/examples/configurations/sample_config.json"
format = "json"`)

	jsonFile := []byte(`{
	"name": "Sample Rulebook",
	"description": "This is a sample rulebook for experimentation.",
	"files": {
		"file1": {
			"path": "/examples/configurations/sample_config.json",
			"format": "json"
		}
	}
}`)

	hoconFile := []byte(`name = "Sample Rulebook"
description = "This is a sample rulebook for experimentation."

# Files to be checked
files {
	file1 {
		path = "/examples/configurations/sample_config.json"
		format = "json"
	}
}`)

	testCases := []testCase{
		{ // Three files
			input: map[string][]byte{
				"toml":  tomlFile,
				"json":  jsonFile,
				"hocon": hoconFile,
			},
			expected: map[string]map[int]string{
				"toml": {
					1: "# This is a TOML configuration for the Rulebook",
					2: "name = \"Sample Rulebook\"",
					3: "description = \"This is a sample rulebook for experimentation.\"",
					4: "",
					5: "# Files to be checked",
					6: "[files.file1]",
					7: "path = \"/examples/configurations/sample_config.json\"",
					8: "format = \"json\"",
				},
				"json": {
					1:  "{",
					2:  "\t\"name\": \"Sample Rulebook\",",
					3:  "\t\"description\": \"This is a sample rulebook for experimentation.\",",
					4:  "\t\"files\": {",
					5:  "\t\t\"file1\": {",
					6:  "\t\t\t\"path\": \"/examples/configurations/sample_config.json\",",
					7:  "\t\t\t\"format\": \"json\"",
					8:  "\t\t}",
					9:  "\t}",
					10: "}",
				},
				"hocon": {
					1:  "name = \"Sample Rulebook\"",
					2:  "description = \"This is a sample rulebook for experimentation.\"",
					3:  "",
					4:  "# Files to be checked",
					5:  "files {",
					6:  "\tfile1 {",
					7:  "\t\tpath = \"/examples/configurations/sample_config.json\"",
					8:  "\t\tformat = \"json\"",
					9:  "\t}",
					10: "}",
				},
			},
		},
		{ // Empty input
			input:    map[string][]byte{},
			expected: map[string]map[int]string{},
		},
	}

	// Run tests
	for _, test := range testCases {
		actual := CreateLinesMapForFiles(test.input)

		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("Expected %v, got %v", test.expected, actual)
		}
	}
}
