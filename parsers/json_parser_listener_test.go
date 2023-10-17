package parsers

import (
	"reflect"
	"testing"
	"fmt"
)

var simpleJSONInput = []byte(`{
		"name": "sample",
		"version": 1.3,
		"active": true,
		"settings": {
			"theme": "dark",
			"notifications": null,
			"retryCount": 3
		},
		"features": ["auth", "logs"]
}`)

var wrongJSONInput = []byte(`{
	"name": "sample",
	"version": 1.3,
	"active": true
	"features": ["auth", "logs"]
`)

type LocationRange struct {
    Start Location
    End   Location
}

// TestSimple_Parse tests the Parse function of a *JsonParser using a simple json config.
func TestSimple_Parse(t *testing.T) {
	// Test cases
	type testCase struct {
		input    []byte
		expected *Node
		err      bool
	}

	// Mock Node result
	expectedNode := &Node{
		Type: Object,
		Value: map[string]*Node{
			"name": {
				Type: String,
				Value: "sample",
				NameLocation: LocationRange{Start: Location{Line: 2, Column: 3}, End: Location{Line: 2, Column: 11}},
				ValueLocation: LocationRange{Start: Location{Line: 2, Column: 11}, End: Location{Line: 2, Column: 11}},
			},
			"version": {
				Type: Float,
				Value: 1.3,
				NameLocation: LocationRange{Start: Location{Line: 3, Column: 3}, End: Location{Line: 3, Column: 14}},
				ValueLocation: LocationRange{Start: Location{Line: 3, Column: 14}, End: Location{Line: 3, Column: 14}},
			},
			"active": {
				Type: Bool,
				Value: true,
				NameLocation: LocationRange{Start: Location{Line: 4, Column: 3}, End: Location{Line: 4, Column: 13}},
				ValueLocation: LocationRange{Start: Location{Line: 4, Column: 13}, End: Location{Line: 4, Column: 13}},
			},
			"settings": {
				Type: Object,
				Value: map[string]*Node{
					"theme": {
						Type: String,
						Value: "dark",
						NameLocation: LocationRange{Start: Location{Line: 6, Column: 4}, End: Location{Line: 6, Column: 13}},
						ValueLocation: LocationRange{Start: Location{Line: 6, Column: 13}, End: Location{Line: 6, Column: 13}},
					},
					"notifications": {
						Type: Null,
						Value: nil,
						NameLocation: LocationRange{Start: Location{Line: 7, Column: 4}, End: Location{Line: 7, Column: 21}},
						ValueLocation: LocationRange{Start: Location{Line: 7, Column: 21}, End: Location{Line: 7, Column: 21}},
					},
					"retryCount": {
						Type: Int,
						Value: 3,
						NameLocation: LocationRange{Start: Location{Line: 8, Column: 4}, End: Location{Line: 8, Column: 18}},
						ValueLocation: LocationRange{Start: Location{Line: 8, Column: 18}, End: Location{Line: 8, Column: 18}},
					},
				},
				NameLocation: LocationRange{Start: Location{Line: 5, Column: 3}, End: Location{Line: 9, Column: 3}},
				ValueLocation: LocationRange{Start: Location{Line: 5, Column: 15}, End: Location{Line: 9, Column: 3}},
			},
			"features": {
				Type:      Array,
				ArrayType: String,
				Value:     []*Node{
					{
						Type: String,
						Value: "auth",
						NameLocation: LocationRange{Start: Location{Line: 0, Column: 0}, End: Location{Line: 0, Column: 0}},
						ValueLocation: LocationRange{Start: Location{Line: 10, Column: 16}, End: Location{Line: 0, Column: 0}},
					},
					{
						Type: String,
						Value: "logs",
						NameLocation: LocationRange{Start: Location{Line: 0, Column: 0}, End: Location{Line: 0, Column: 0}},
						ValueLocation: LocationRange{Start: Location{Line: 10, Column: 24}, End: Location{Line: 0, Column: 0}},
					},
				},
				NameLocation: LocationRange{Start: Location{Line: 10, Column: 3}, End: Location{Line: 10, Column: 30}},
				ValueLocation: LocationRange{Start: Location{Line: 10, Column: 15}, End: Location{Line: 10, Column: 30}},
			},
		},
		NameLocation: LocationRange{Start: Location{Line: 0, Column: 0}, End: Location{Line: 0, Column: 0}},
		ValueLocation: LocationRange{Start: Location{Line: 1, Column: 1}, End: Location{Line: 11, Column: 1}},
	}	

	testCases := []testCase{
		{
			input: simpleJSONInput,
			expected: expectedNode,
			err: false,
		},
	}

	// Run tests
	for _, test := range testCases {
		parser := &JsonParser{}
		actual, err := parser.Parse(test.input)

		if err != nil && !test.err {
			t.Errorf("Unexpected error: %s", err)
		} else if err == nil && test.err {
			t.Errorf("Expected error, got nil")
		} else if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("Expected %v, got %v", test.expected, actual)
		}
	}
}

// TestSingleError_Parse tests the Parse function of a *JsonParser using a wrong json config.
func TestSingleError_Parse(t *testing.T) {
	type testCase struct {
		input    []byte
		expected error
		err      bool
	}

	testCases := []testCase{
		{
			input:    wrongJSONInput,
			expected: fmt.Errorf(`Syntax errors: [line 5:1 extraneous input '"features"' expecting {',', '}'} line 5:28 mismatched input ']' expecting ':']`),
			err:      true,
		},
	}

	for _, test := range testCases {
		parser := &JsonParser{}
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
