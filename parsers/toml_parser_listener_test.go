package parsers

import (
	"fmt"
	"testing"
)

// TestSimple_Toml_Parse tests the Parse function of a *TomlParser using a simple toml config.
func TestSimple_Toml_Parse(t *testing.T) {
	// Input
	var simpleTOMLInput = []byte(`
		[server]
		host = "localhost"

		[database]
		name = "mydb"

		[database.user]
		name = "myuser"
		password = "mypassword"
	`)

	// Test cases
	type testCase struct {
		input    []byte
		expected error
		err      bool
	}

	// Mock Node result
	testCases := []testCase{
		{
			input:    simpleTOMLInput,
			expected: nil,
			err:      false,
		},
	}

	// Run tests
	for _, test := range testCases {
		parser := &TomlParser{}
		_, err := parser.Parse(test.input)

		if err != nil && !test.err {
			t.Errorf("Unexpected error: %s", err)
		} else if err == nil && test.err {
			t.Errorf("Expected error, got nil")
		} else if err != nil && test.err && err.Error() != test.expected.Error() {
			t.Errorf("Expected %v, got %v", test.expected, err)
		}
	}
}

// TestSingleError_Toml_Parse tests the Parse function of a *TomlParser using a wrong toml config.
func TestSyntaxError_Toml_Parse(t *testing.T) {
	// Input
	var syntaxErrorInput = []byte(`
		[server]
		host = "localhost"
		port = 8080

		[database]
		name = mydb
		user = "myuser"
		password = "mypassword"
	`)

	// Test cases
	type testCase struct {
		input    []byte
		expected error
		err      bool
	}

	// Mock Node result
	testCases := []testCase{
		{
			input:    syntaxErrorInput,
			expected: fmt.Errorf(`Syntax errors: [line 7:9 mismatched input 'mydb' expecting {'[', '{', BOOLEAN, BASIC_STRING, ML_BASIC_STRING, LITERAL_STRING, ML_LITERAL_STRING, FLOAT, INF, NAN, DEC_INT, HEX_INT, OCT_INT, BIN_INT, OFFSET_DATE_TIME, LOCAL_DATE_TIME, LOCAL_DATE, LOCAL_TIME}]`),
			err:      true,
		},
	}

	// Run tests
	for _, test := range testCases {
		parser := &TomlParser{}
		_, err := parser.Parse(test.input)

		if err != nil && !test.err {
			t.Errorf("Unexpected error: %s", err)
		} else if err == nil && test.err {
			t.Errorf("Expected error, got nil")
		} else if err != nil && test.err && err.Error() != test.expected.Error() {
			t.Errorf("Expected %v, got %v", test.expected, err)
		}
	}
}
