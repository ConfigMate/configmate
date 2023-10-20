package parsers

import (
	"fmt"
	"reflect"
	"testing"
)

// TestSimple_Parse tests the Parse function of a *JsonParser using a simple json config.
func TestSimple_Parse(t *testing.T) {
	// Input
	var jsonConfig = []byte(`{
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
				Type:          String,
				Value:         "sample",
				NameLocation:  TokenLocation{Start: CharLocation{Line: 2, Column: 3}, End: CharLocation{Line: 2, Column: 11}},
				ValueLocation: TokenLocation{Start: CharLocation{Line: 2, Column: 11}, End: CharLocation{Line: 2, Column: 11}},
			},
			"version": {
				Type:          Float,
				Value:         1.3,
				NameLocation:  TokenLocation{Start: CharLocation{Line: 3, Column: 3}, End: CharLocation{Line: 3, Column: 14}},
				ValueLocation: TokenLocation{Start: CharLocation{Line: 3, Column: 14}, End: CharLocation{Line: 3, Column: 14}},
			},
			"active": {
				Type:          Bool,
				Value:         true,
				NameLocation:  TokenLocation{Start: CharLocation{Line: 4, Column: 3}, End: CharLocation{Line: 4, Column: 13}},
				ValueLocation: TokenLocation{Start: CharLocation{Line: 4, Column: 13}, End: CharLocation{Line: 4, Column: 13}},
			},
			"settings": {
				Type: Object,
				Value: map[string]*Node{
					"theme": {
						Type:          String,
						Value:         "dark",
						NameLocation:  TokenLocation{Start: CharLocation{Line: 6, Column: 4}, End: CharLocation{Line: 6, Column: 13}},
						ValueLocation: TokenLocation{Start: CharLocation{Line: 6, Column: 13}, End: CharLocation{Line: 6, Column: 13}},
					},
					"notifications": {
						Type:          Null,
						Value:         nil,
						NameLocation:  TokenLocation{Start: CharLocation{Line: 7, Column: 4}, End: CharLocation{Line: 7, Column: 21}},
						ValueLocation: TokenLocation{Start: CharLocation{Line: 7, Column: 21}, End: CharLocation{Line: 7, Column: 21}},
					},
					"retryCount": {
						Type:          Int,
						Value:         3,
						NameLocation:  TokenLocation{Start: CharLocation{Line: 8, Column: 4}, End: CharLocation{Line: 8, Column: 18}},
						ValueLocation: TokenLocation{Start: CharLocation{Line: 8, Column: 18}, End: CharLocation{Line: 8, Column: 18}},
					},
				},
				NameLocation:  TokenLocation{Start: CharLocation{Line: 5, Column: 3}, End: CharLocation{Line: 9, Column: 3}},
				ValueLocation: TokenLocation{Start: CharLocation{Line: 5, Column: 15}, End: CharLocation{Line: 9, Column: 3}},
			},
			"features": {
				Type:      Array,
				ArrayType: String,
				Value: []*Node{
					{
						Type:          String,
						Value:         "auth",
						NameLocation:  TokenLocation{Start: CharLocation{Line: 0, Column: 0}, End: CharLocation{Line: 0, Column: 0}},
						ValueLocation: TokenLocation{Start: CharLocation{Line: 10, Column: 16}, End: CharLocation{Line: 0, Column: 0}},
					},
					{
						Type:          String,
						Value:         "logs",
						NameLocation:  TokenLocation{Start: CharLocation{Line: 0, Column: 0}, End: CharLocation{Line: 0, Column: 0}},
						ValueLocation: TokenLocation{Start: CharLocation{Line: 10, Column: 24}, End: CharLocation{Line: 0, Column: 0}},
					},
				},
				NameLocation:  TokenLocation{Start: CharLocation{Line: 10, Column: 3}, End: CharLocation{Line: 10, Column: 30}},
				ValueLocation: TokenLocation{Start: CharLocation{Line: 10, Column: 15}, End: CharLocation{Line: 10, Column: 30}},
			},
		},
		NameLocation:  TokenLocation{Start: CharLocation{Line: 0, Column: 0}, End: CharLocation{Line: 0, Column: 0}},
		ValueLocation: TokenLocation{Start: CharLocation{Line: 1, Column: 1}, End: CharLocation{Line: 11, Column: 1}},
	}

	testCases := []testCase{
		{
			input:    jsonConfig,
			expected: expectedNode,
			err:      false,
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

// TestEmptyObj_Json_Parse tests the Parse function of a *JsonParser using an empty object json config.
func TestEmptyObj_Json_Parse(t *testing.T) {
	// Input
	var jsonConfig = []byte(`{}`)

	// Test cases
	type testCase struct {
		input    []byte
		expected *Node
		err      bool
	}

	// Mock Node result
	expectedNode := &Node{
		Type:          Object,
		Value:         map[string]*Node{},
		NameLocation:  TokenLocation{Start: CharLocation{Line: 0, Column: 0}, End: CharLocation{Line: 0, Column: 0}},
		ValueLocation: TokenLocation{Start: CharLocation{Line: 1, Column: 1}, End: CharLocation{Line: 1, Column: 2}},
	}

	testCases := []testCase{
		{
			input:    jsonConfig,
			expected: expectedNode,
			err:      false,
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

// TestEmptyArray_Json_Parse tests the Parse function of a *JsonParser using an empty array json config.
func TestEmptyArray_Json_Parse(t *testing.T) {
	// Input
	var jsonConfig = []byte(`[]`)

	// Test cases
	type testCase struct {
		input    []byte
		expected *Node
		err      bool
	}

	// Mock Node result
	expectedNode := &Node{
		Type:          Array,
		Value:         []*Node{},
		NameLocation:  TokenLocation{Start: CharLocation{Line: 0, Column: 0}, End: CharLocation{Line: 0, Column: 0}},
		ValueLocation: TokenLocation{Start: CharLocation{Line: 1, Column: 1}, End: CharLocation{Line: 1, Column: 2}},
	}

	testCases := []testCase{
		{
			input:    jsonConfig,
			expected: expectedNode,
			err:      false,
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

func TestSingleArray_Json_Parse(t *testing.T) {
	// Input
	var jsonConfig = []byte(`["sample"]`)

	// Test cases
	type testCase struct {
		input    []byte
		expected *Node
		err      bool
	}

	// Mock Node result
	expectedNode := &Node{
		Type:      Array,
		ArrayType: String,
		Value: []*Node{
			{
				Type:          String,
				Value:         "sample",
				NameLocation:  TokenLocation{Start: CharLocation{Line: 0, Column: 0}, End: CharLocation{Line: 0, Column: 0}},
				ValueLocation: TokenLocation{Start: CharLocation{Line: 1, Column: 2}, End: CharLocation{Line: 0, Column: 0}},
			},
		},
		NameLocation:  TokenLocation{Start: CharLocation{Line: 0, Column: 0}, End: CharLocation{Line: 0, Column: 0}},
		ValueLocation: TokenLocation{Start: CharLocation{Line: 1, Column: 1}, End: CharLocation{Line: 1, Column: 10}},
	}

	testCases := []testCase{
		{
			input:    jsonConfig,
			expected: expectedNode,
			err:      false,
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

// TestSingleString_JSON_Parse tests the Parse function of a *JsonParser using a single string json config.
func TestSingleString_Json_Parse(t *testing.T) {
	// Input
	var jsonConfig = []byte(`"sample"`)

	// Test cases
	type testCase struct {
		input    []byte
		expected *Node
		err      bool
	}

	// Mock Node result
	expectedNode := &Node{
		Type:          String,
		Value:         "sample",
		NameLocation:  TokenLocation{Start: CharLocation{Line: 1, Column: 1}, End: CharLocation{Line: 0, Column: 0}},
		ValueLocation: TokenLocation{Start: CharLocation{Line: 0, Column: 0}, End: CharLocation{Line: 1, Column: 1}},
	}

	testCases := []testCase{
		{
			input:    jsonConfig,
			expected: expectedNode,
			err:      false,
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

// TestSingleIntNumber_Json_Parse tests the Parse function of a *JsonParser using a single int number json config.
func TestSingleIntNumber_Json_Parse(t *testing.T) {
	// Input
	var jsonConfig = []byte(`12345`)

	// Test cases
	type testCase struct {
		input    []byte
		expected *Node
		err      bool
	}

	// Mock Node result
	expectedNode := &Node{
		Type:          Int,
		Value:         12345,
		NameLocation:  TokenLocation{Start: CharLocation{Line: 1, Column: 1}, End: CharLocation{Line: 0, Column: 0}},
		ValueLocation: TokenLocation{Start: CharLocation{Line: 0, Column: 0}, End: CharLocation{Line: 1, Column: 1}},
	}

	testCases := []testCase{
		{
			input:    jsonConfig,
			expected: expectedNode,
			err:      false,
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

// TestSingleFloatNumber_Json_Parse tests the Parse function of a *JsonParser using a single float json config.
func TestSingleFloatNumber_Json_Parse(t *testing.T) {
	// Input
	var jsonConfig = []byte(`123.45`)

	// Test cases
	type testCase struct {
		input    []byte
		expected *Node
		err      bool
	}

	// Mock Node result
	expectedNode := &Node{
		Type:          Float,
		Value:         123.45,
		NameLocation:  TokenLocation{Start: CharLocation{Line: 1, Column: 1}, End: CharLocation{Line: 0, Column: 0}},
		ValueLocation: TokenLocation{Start: CharLocation{Line: 0, Column: 0}, End: CharLocation{Line: 1, Column: 1}},
	}

	testCases := []testCase{
		{
			input:    jsonConfig,
			expected: expectedNode,
			err:      false,
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

// TestSingleTrueBool_Json_Parse tests the Parse function of a *JsonParser using a single true bool json config.
func TestSingleTrueBool_Json_Parse(t *testing.T) {
	// Input
	var jsonConfig = []byte(`true`)

	// Test cases
	type testCase struct {
		input    []byte
		expected *Node
		err      bool
	}

	// Mock Node result
	expectedNode := &Node{
		Type:          Bool,
		Value:         true,
		NameLocation:  TokenLocation{Start: CharLocation{Line: 1, Column: 1}, End: CharLocation{Line: 0, Column: 0}},
		ValueLocation: TokenLocation{Start: CharLocation{Line: 0, Column: 0}, End: CharLocation{Line: 1, Column: 1}},
	}

	testCases := []testCase{
		{
			input:    jsonConfig,
			expected: expectedNode,
			err:      false,
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

// TestSingleFalseBool_Json_Parse tests the Parse function of a *JsonParser using a single false bool json config.
func TestSingleFalseBool_Json_Parse(t *testing.T) {
	// Input
	var jsonConfig = []byte(`false`)

	// Test cases
	type testCase struct {
		input    []byte
		expected *Node
		err      bool
	}

	// Mock Node result
	expectedNode := &Node{
		Type:          Bool,
		Value:         false,
		NameLocation:  TokenLocation{Start: CharLocation{Line: 1, Column: 1}, End: CharLocation{Line: 0, Column: 0}},
		ValueLocation: TokenLocation{Start: CharLocation{Line: 0, Column: 0}, End: CharLocation{Line: 1, Column: 1}},
	}

	testCases := []testCase{
		{
			input:    jsonConfig,
			expected: expectedNode,
			err:      false,
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

// TestSingleNull_Json_Parse tests the Parse function of a *JsonParser using a single null json config.
func TestSingleNull_Json_Parse(t *testing.T) {
	// Input
	var jsonConfig = []byte(`null`)

	// Test cases
	type testCase struct {
		input    []byte
		expected *Node
		err      bool
	}

	// Mock Node result
	expectedNode := &Node{
		Type:          Null,
		Value:         nil,
		NameLocation:  TokenLocation{Start: CharLocation{Line: 1, Column: 1}, End: CharLocation{Line: 0, Column: 0}},
		ValueLocation: TokenLocation{Start: CharLocation{Line: 0, Column: 0}, End: CharLocation{Line: 1, Column: 1}},
	}

	testCases := []testCase{
		{
			input:    jsonConfig,
			expected: expectedNode,
			err:      false,
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
	// Input
	var jsonConfig = []byte(`{
		"name": "sample",
		"version": 1.3,
		"active": true
		"features": ["auth", "logs"]
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
			input:    jsonConfig,
			expected: fmt.Errorf(`syntax errors: [line 5:2 extraneous input '"features"' expecting {',', '}'} line 5:29 mismatched input ']' expecting ':']`),
			err:      true,
		},
	}

	// Run tests
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
