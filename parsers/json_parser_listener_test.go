package parsers

import (
	"reflect"
	"testing"
)

var JSONInput = []byte(`{
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


// TestNode_Get tests the Get function of a *Node.
func Test_Parse(t *testing.T) {
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
				NameLocation: struct {
					Line   int
					Column int
					Length int
				} {Line: 2, Column: 2, Length: 6},
				ValueLocation: struct {
					Line   int
					Column int
					Length int
				} {Line: 2, Column: 10, Length: 8},
			},
			"version": {
				Type: Float, 
				Value: 1.3,
				NameLocation: struct {
					Line   int
					Column int
					Length int
				} {Line: 3, Column: 2, Length: 9},
				ValueLocation: struct {
					Line   int
					Column int
					Length int
				} {Line: 3, Column: 13, Length: 3},
			},
			"active": {
				Type: Bool, 
				Value: true,
				NameLocation: struct {
					Line   int
					Column int
					Length int
				} {Line: 4, Column: 2, Length: 8},
				ValueLocation: struct {
					Line   int
					Column int
					Length int
				} {Line: 4, Column: 12, Length: 4},
			},
			"settings": {
				Type: Object,
				Value: map[string]*Node{
					"theme": {
						Type: String, 
						Value: "dark",
						NameLocation: struct {
							Line   int
							Column int
							Length int
						} {Line: 6, Column: 3, Length: 7},
						ValueLocation: struct {
							Line   int
							Column int
							Length int
						} {Line: 6, Column: 12, Length: 6},
					},
					"notifications": {
						Type: Null, 
						Value: nil,
						NameLocation: struct {
							Line   int
							Column int
							Length int
						} {Line: 7, Column: 3, Length: 15},
						ValueLocation: struct {
							Line   int
							Column int
							Length int
						} {Line: 7, Column: 20, Length: 4},
					},
					"retryCount": {
						Type: Int, 
						Value: 3,
						NameLocation: struct {
							Line   int
							Column int
							Length int
						} {Line: 8, Column: 3, Length: 12},
						ValueLocation: struct {
							Line   int
							Column int
							Length int
						} {Line: 8, Column: 17, Length: 1},
					},
				},
				NameLocation: struct {
					Line   int
					Column int
					Length int
				} {Line: 5, Column: 2, Length: 10},
				ValueLocation: struct {
					Line   int
					Column int
					Length int
				} {Line: 5, Column: 14, Length: 52},
			},
			"features": {
				Type:      Array,
				ArrayType: String,
				Value:     []*Node{
					{
						Type: String, 
						Value: "auth",
						NameLocation: struct {
							Line   int
							Column int
							Length int
							} {Line: 0, Column: 0, Length: 0},
						ValueLocation: struct {
							Line   int
							Column int
							Length int
						} {Line: 10, Column: 15, Length: 6},
					}, 
					{
							Type: String, 
							Value: "logs",
							NameLocation: struct {
								Line   int
								Column int
								Length int
							} {Line: 0, Column: 0, Length: 0},
							ValueLocation: struct {
								Line   int
								Column int
								Length int
							} {Line: 10, Column: 23, Length: 6},
					},
				},
				NameLocation: struct {
					Line   int
					Column int
					Length int
				} {Line: 10, Column: 2, Length: 10},
				ValueLocation: struct {
					Line   int
					Column int
					Length int
				} {Line: 10, Column: 14, Length: 15},
			},
		},
		NameLocation: struct {
			Line   int
			Column int
			Length int
		} {Line: 0, Column: 0, Length: 0},
		ValueLocation: struct {
			Line   int
			Column int
			Length int
		} {Line: 1, Column: 0, Length: 136},
	}

	testCases := []testCase{
		{
			input: JSONInput,
			expected: expectedNode,
			err: false,
		},
	}

	// Run tests
	for i, test := range testCases {
		parser := &JsonParser{}
		actual, err := parser.Parse(test.input)

		if err != nil && !test.err {
			t.Errorf("Test case %d: Unexpected error: %s", i, err)
		} else if err == nil && test.err {
			t.Errorf("Test case %d: Expected error, got nil", i)
		} else if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("Test case %d: Expected %v, got %v", i, test.expected, actual)
		}
	}
}
