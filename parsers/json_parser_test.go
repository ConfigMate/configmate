package parsers

import (
	"os"
	"reflect"
	"testing"
)

type jsonParserTestCase struct {
	input        []byte
	expected     *Node
	expectedErrs []CMParserError
}

func TestParseSimpleConfig_jsonParser(t *testing.T) {
	// Input
	testConfig, err := os.ReadFile("./test_configs/simple.json")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	// Test cases
	testCases := []jsonParserTestCase{
		{
			input: testConfig,
			expected: &Node{
				Type: Object,
				Value: map[string]*Node{
					"name": {
						Type:          String,
						Value:         "sample",
						NameLocation:  TokenLocation{Start: CharLocation{Line: 1, Column: 1}, End: CharLocation{Line: 1, Column: 7}},  // Indented with tabs instead of spaces
						ValueLocation: TokenLocation{Start: CharLocation{Line: 1, Column: 9}, End: CharLocation{Line: 1, Column: 17}}, // Indented with tabs instead of spaces
					},
					"version": {
						Type:          Float,
						Value:         1.3,
						NameLocation:  TokenLocation{Start: CharLocation{Line: 2, Column: 4}, End: CharLocation{Line: 2, Column: 13}},
						ValueLocation: TokenLocation{Start: CharLocation{Line: 2, Column: 15}, End: CharLocation{Line: 2, Column: 18}},
					},
					"active": {
						Type:          Bool,
						Value:         true,
						NameLocation:  TokenLocation{Start: CharLocation{Line: 3, Column: 4}, End: CharLocation{Line: 3, Column: 12}},
						ValueLocation: TokenLocation{Start: CharLocation{Line: 3, Column: 14}, End: CharLocation{Line: 3, Column: 18}},
					},
					"settings": {
						Type: Object,
						Value: map[string]*Node{
							"theme": {
								Type:          String,
								Value:         "dark",
								NameLocation:  TokenLocation{Start: CharLocation{Line: 5, Column: 8}, End: CharLocation{Line: 5, Column: 15}},
								ValueLocation: TokenLocation{Start: CharLocation{Line: 5, Column: 17}, End: CharLocation{Line: 5, Column: 23}},
							},
							"notifications": {
								Type:          Null,
								Value:         nil,
								NameLocation:  TokenLocation{Start: CharLocation{Line: 6, Column: 8}, End: CharLocation{Line: 6, Column: 23}},
								ValueLocation: TokenLocation{Start: CharLocation{Line: 6, Column: 25}, End: CharLocation{Line: 6, Column: 29}},
							},
							"retryCount": {
								Type:          Int,
								Value:         3,
								NameLocation:  TokenLocation{Start: CharLocation{Line: 7, Column: 8}, End: CharLocation{Line: 7, Column: 20}},
								ValueLocation: TokenLocation{Start: CharLocation{Line: 7, Column: 22}, End: CharLocation{Line: 7, Column: 23}},
							},
						},
						NameLocation:  TokenLocation{Start: CharLocation{Line: 4, Column: 4}, End: CharLocation{Line: 4, Column: 14}},
						ValueLocation: TokenLocation{Start: CharLocation{Line: 4, Column: 16}, End: CharLocation{Line: 8, Column: 5}},
					},
					"features": {
						Type: Array,
						Value: []*Node{
							{
								Type:          String,
								Value:         "auth",
								NameLocation:  TokenLocation{Start: CharLocation{Line: 0, Column: 0}, End: CharLocation{Line: 0, Column: 0}},
								ValueLocation: TokenLocation{Start: CharLocation{Line: 9, Column: 17}, End: CharLocation{Line: 9, Column: 23}},
							},
							{
								Type:          String,
								Value:         "logs",
								NameLocation:  TokenLocation{Start: CharLocation{Line: 0, Column: 0}, End: CharLocation{Line: 0, Column: 0}},
								ValueLocation: TokenLocation{Start: CharLocation{Line: 9, Column: 25}, End: CharLocation{Line: 9, Column: 31}},
							},
						},
						NameLocation:  TokenLocation{Start: CharLocation{Line: 9, Column: 4}, End: CharLocation{Line: 9, Column: 14}},
						ValueLocation: TokenLocation{Start: CharLocation{Line: 9, Column: 16}, End: CharLocation{Line: 9, Column: 32}},
					},
				},
				NameLocation:  TokenLocation{Start: CharLocation{Line: 0, Column: 0}, End: CharLocation{Line: 0, Column: 0}},
				ValueLocation: TokenLocation{Start: CharLocation{Line: 0, Column: 0}, End: CharLocation{Line: 10, Column: 1}},
			},
			expectedErrs: []CMParserError{},
		},
	}

	// Run tests
	for _, test := range testCases {
		parser := &jsonParser{}
		result, errs := parser.Parse(test.input)

		if len(errs) > 0 {
			t.Errorf("Unexpected errors: %#v", errs)
		} else if !reflect.DeepEqual(test.expected, result) {
			t.Errorf("Expected %v, got %v", test.expected, result)
		}
	}
}

func TestSimpleExpressions_jsonParser(t *testing.T) {
	// Test cases
	testCases := []jsonParserTestCase{
		{ // empty object json config
			input: []byte(`{}`),
			expected: &Node{
				Type:          Object,
				Value:         map[string]*Node{},
				NameLocation:  TokenLocation{Start: CharLocation{Line: 0, Column: 0}, End: CharLocation{Line: 0, Column: 0}},
				ValueLocation: TokenLocation{Start: CharLocation{Line: 0, Column: 0}, End: CharLocation{Line: 0, Column: 2}},
			},
			expectedErrs: nil,
		},
		{ // empty array json config
			input: []byte(`[]`),
			expected: &Node{
				Type:          Array,
				Value:         []*Node{},
				NameLocation:  TokenLocation{Start: CharLocation{Line: 0, Column: 0}, End: CharLocation{Line: 0, Column: 0}},
				ValueLocation: TokenLocation{Start: CharLocation{Line: 0, Column: 0}, End: CharLocation{Line: 0, Column: 2}},
			},
			expectedErrs: nil,
		},
		{ // simple array json config
			input: []byte(`["sample"]`),
			expected: &Node{
				Type: Array,
				Value: []*Node{
					{
						Type:          String,
						Value:         "sample",
						NameLocation:  TokenLocation{Start: CharLocation{Line: 0, Column: 0}, End: CharLocation{Line: 0, Column: 0}},
						ValueLocation: TokenLocation{Start: CharLocation{Line: 0, Column: 1}, End: CharLocation{Line: 0, Column: 9}},
					},
				},
				NameLocation:  TokenLocation{Start: CharLocation{Line: 0, Column: 0}, End: CharLocation{Line: 0, Column: 0}},
				ValueLocation: TokenLocation{Start: CharLocation{Line: 0, Column: 0}, End: CharLocation{Line: 0, Column: 10}},
			},
			expectedErrs: nil,
		},
		{ // single string json config
			input: []byte(`"sample"`),
			expected: &Node{
				Type:          String,
				Value:         "sample",
				NameLocation:  TokenLocation{Start: CharLocation{Line: 0, Column: 0}, End: CharLocation{Line: 0, Column: 0}},
				ValueLocation: TokenLocation{Start: CharLocation{Line: 0, Column: 0}, End: CharLocation{Line: 0, Column: 8}},
			},
			expectedErrs: nil,
		},
		{ // single int number json config
			input: []byte(`12345`),
			expected: &Node{
				Type:          Int,
				Value:         12345,
				NameLocation:  TokenLocation{Start: CharLocation{Line: 0, Column: 0}, End: CharLocation{Line: 0, Column: 0}},
				ValueLocation: TokenLocation{Start: CharLocation{Line: 0, Column: 0}, End: CharLocation{Line: 0, Column: 5}},
			},
			expectedErrs: nil,
		},
		{ // single float number json config
			input: []byte(`123.45`),
			expected: &Node{
				Type:          Float,
				Value:         123.45,
				NameLocation:  TokenLocation{Start: CharLocation{Line: 0, Column: 0}, End: CharLocation{Line: 0, Column: 0}},
				ValueLocation: TokenLocation{Start: CharLocation{Line: 0, Column: 0}, End: CharLocation{Line: 0, Column: 6}},
			},
			expectedErrs: nil,
		},
		{ // single true bool json config
			input: []byte(`true`),
			expected: &Node{
				Type:          Bool,
				Value:         true,
				NameLocation:  TokenLocation{Start: CharLocation{Line: 0, Column: 0}, End: CharLocation{Line: 0, Column: 0}},
				ValueLocation: TokenLocation{Start: CharLocation{Line: 0, Column: 0}, End: CharLocation{Line: 0, Column: 4}},
			},
			expectedErrs: nil,
		},
		{ // single false bool json config
			input: []byte(`false`),
			expected: &Node{
				Type:          Bool,
				Value:         false,
				NameLocation:  TokenLocation{Start: CharLocation{Line: 0, Column: 0}, End: CharLocation{Line: 0, Column: 0}},
				ValueLocation: TokenLocation{Start: CharLocation{Line: 0, Column: 0}, End: CharLocation{Line: 0, Column: 5}},
			},
			expectedErrs: nil,
		},
		{ // single null json config
			input: []byte(`null`),
			expected: &Node{
				Type:          Null,
				Value:         nil,
				NameLocation:  TokenLocation{Start: CharLocation{Line: 0, Column: 0}, End: CharLocation{Line: 0, Column: 0}},
				ValueLocation: TokenLocation{Start: CharLocation{Line: 0, Column: 0}, End: CharLocation{Line: 0, Column: 4}},
			},
			expectedErrs: nil,
		},
	}

	// Run tests
	for _, test := range testCases {
		parser := &jsonParser{}
		result, errs := parser.Parse(test.input)

		if len(errs) > 0 {
			t.Errorf("Unexpected errors: %#v", errs)
		} else if !reflect.DeepEqual(test.expected, result) {
			t.Errorf("Expected %v, got %v", test.expected, result)
		}
	}
}

func TestErrorConditions_jsonParser(t *testing.T) {
	// Input
	var errJsonConfig0 = []byte(`{
		"name": "sample",
		"version": 1.3,
		"active": true
		"features": ["auth", "logs"]
	`)

	testCases := []jsonParserTestCase{
		{
			input:    errJsonConfig0,
			expected: nil,
			expectedErrs: []CMParserError{
				{
					Message: "extraneous input '\"features\"' expecting {',', '}'}",
					Location: TokenLocation{
						Start: CharLocation{Line: 4, Column: 2},
						End:   CharLocation{Line: 4, Column: 3},
					},
				},
				{
					Message: "mismatched input ']' expecting ':'",
					Location: TokenLocation{
						Start: CharLocation{Line: 4, Column: 29},
						End:   CharLocation{Line: 4, Column: 30},
					},
				},
			},
		},
	}

	// Run tests
	for _, test := range testCases {
		parser := &jsonParser{}
		_, errs := parser.Parse(test.input)

		if len(errs) == 0 {
			t.Errorf("Expected errors, got none")
		} else if !reflect.DeepEqual(test.expectedErrs, errs) {
			t.Errorf("Expected %v, got %v", test.expectedErrs, errs)
		}
	}
}
