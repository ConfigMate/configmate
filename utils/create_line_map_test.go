package utils

import (
	"reflect"
	"testing"
)

// TestTomlFile_CreateLineMap tests the CreateLineMap function of the TomlFile struct.
func TestTomlFile_CreateLineMap(t *testing.T) {
	// Input
	var tomlFile = []byte(`# This is a TOML configuration for the Rulebook
name = "Sample Rulebook"
description = "This is a sample rulebook for experimentation."

# Files to be checked
[files.file1]
path = "/examples/configurations/sample_config.json"
format = "json"`)

	// Test cases
	type testCase struct {
		input    []byte
		expected map[int]string
		err      bool
	}

	// Mock result
	expectedMap := map[int]string{
		1: "# This is a TOML configuration for the Rulebook",
		2: "name = \"Sample Rulebook\"",
		3: "description = \"This is a sample rulebook for experimentation.\"",
		4: "",
		5: "# Files to be checked",
		6: "[files.file1]",
		7: "path = \"/examples/configurations/sample_config.json\"",
		8: "format = \"json\"",
	}

	testCases := []testCase{
		{
			input:    tomlFile,
			expected: expectedMap,
			err:      false,
		},
	}

	// Run tests
	for _, test := range testCases {
		actual, err := CreateLineMap(test.input)

		if err != nil && !test.err {
			t.Errorf("Unexpected error: %s", err)
		} else if err == nil && test.err {
			t.Errorf("Expected error, got nil")
		} else if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("Expected %v, got %v", test.expected, actual)
		}
	}
}

// TestJsonFile_CreateLineMap tests the CreateLineMap function of the JsonFile struct.
func TestJsonFile_CreateLineMap(t *testing.T) {
	// Input
	var jsonFile = []byte(`{
	"name": "Sample Rulebook",
	"description": "This is a sample rulebook for experimentation.",
	"files": {
		"file1": {
			"path": "/examples/configurations/sample_config.json",
			"format": "json"
		}
	}
}`)

	// Test cases
	type testCase struct {
		input    []byte
		expected map[int]string
		err      bool
	}

	// Mock result
	expectedMap := map[int]string{
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
	}

	testCases := []testCase{
		{
			input:    jsonFile,
			expected: expectedMap,
			err:      false,
		},
	}

	// Run tests
	for _, test := range testCases {
		actual, err := CreateLineMap(test.input)

		if err != nil && !test.err {
			t.Errorf("Unexpected error: %s", err)
		} else if err == nil && test.err {
			t.Errorf("Expected error, got nil")
		} else if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("Expected %v, got %v", test.expected, actual)
		}
	}
}

func TestHoconFile_CreateLineMap(t *testing.T) {
	// Input
	var hoconFile = []byte(`name = "Sample Rulebook"
description = "This is a sample rulebook for experimentation."

# Files to be checked
files {
	file1 {
		path = "/examples/configurations/sample_config.json"
		format = "json"
	}
}`)

	// Test cases
	type testCase struct {
		input    []byte
		expected map[int]string
		err      bool
	}

	// Mock result
	expectedMap := map[int]string{
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
	}

	testCases := []testCase{
		{
			input:    hoconFile,
			expected: expectedMap,
			err:      false,
		},
	}

	// Run tests
	for _, test := range testCases {
		actual, err := CreateLineMap(test.input)

		if err != nil && !test.err {
			t.Errorf("Unexpected error: %s", err)
		} else if err == nil && test.err {
			t.Errorf("Expected error, got nil")
		} else if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("Expected %v, got %v", test.expected, actual)
		}
	}
}
