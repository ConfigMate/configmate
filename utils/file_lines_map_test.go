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
				0: "# This is a TOML configuration for the Rulebook",
				1: "name = \"Sample Rulebook\"",
				2: "description = \"This is a sample rulebook for experimentation.\"",
				3: "",
				4: "# Files to be checked",
				5: "[files.file1]",
				6: "path = \"/examples/configurations/sample_config.json\"",
				7: "format = \"json\"",
			},
		},
		{ // JSON file
			input: jsonFile,
			expected: map[int]string{
				0: "{",
				1: "\t\"name\": \"Sample Rulebook\",",
				2: "\t\"description\": \"This is a sample rulebook for experimentation.\",",
				3: "\t\"files\": {",
				4: "\t\t\"file1\": {",
				5: "\t\t\t\"path\": \"/examples/configurations/sample_config.json\",",
				6: "\t\t\t\"format\": \"json\"",
				7: "\t\t}",
				8: "\t}",
				9: "}",
			},
		},
		{ // HOCON file
			input: hoconFile,
			expected: map[int]string{
				0: "name = \"Sample Rulebook\"",
				1: "description = \"This is a sample rulebook for experimentation.\"",
				2: "",
				3: "# Files to be checked",
				4: "files {",
				5: "\tfile1 {",
				6: "\t\tpath = \"/examples/configurations/sample_config.json\"",
				7: "\t\tformat = \"json\"",
				8: "\t}",
				9: "}",
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
					0: "# This is a TOML configuration for the Rulebook",
					1: "name = \"Sample Rulebook\"",
					2: "description = \"This is a sample rulebook for experimentation.\"",
					3: "",
					4: "# Files to be checked",
					5: "[files.file1]",
					6: "path = \"/examples/configurations/sample_config.json\"",
					7: "format = \"json\"",
				},
				"json": {
					0: "{",
					1: "\t\"name\": \"Sample Rulebook\",",
					2: "\t\"description\": \"This is a sample rulebook for experimentation.\",",
					3: "\t\"files\": {",
					4: "\t\t\"file1\": {",
					5: "\t\t\t\"path\": \"/examples/configurations/sample_config.json\",",
					6: "\t\t\t\"format\": \"json\"",
					7: "\t\t}",
					8: "\t}",
					9: "}",
				},
				"hocon": {
					0: "name = \"Sample Rulebook\"",
					1: "description = \"This is a sample rulebook for experimentation.\"",
					2: "",
					3: "# Files to be checked",
					4: "files {",
					5: "\tfile1 {",
					6: "\t\tpath = \"/examples/configurations/sample_config.json\"",
					7: "\t\tformat = \"json\"",
					8: "\t}",
					9: "}",
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
