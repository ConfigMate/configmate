package parsers

import (
	"reflect"
	"testing"
)

var JSONInput = []byte(`
{
	"name": "sample",
	"version": 1.3,
	"active": true,
	"settings": {
	  "theme": "dark",
	  "notifications": null,
	  "retryCount": 3
	},
	"features": ["auth", "logs"]
}
`)


// type Node struct {
// 	Type      FieldType   // Type of field
// 	ArrayType FieldType   // Type of elements in array (if Type == Array)
// 	Value     interface{} // Value of field

// 	NameLocation struct { // Location of field name in configuration file
// 		Line   int
// 		Column int
// 		Length int
// 	}
// 	ValueLocation struct { // Location of field value in configuration file
// 		Line   int
// 		Column int
// 		Length int
// 	}
// }

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
			"name": {Type: String, Value: "sample"},
			"version": {Type: Float, Value: 1.3},
			"active": {Type: Bool, Value: true},
			"settings": {
				Type: Object,
				Value: map[string]*Node{
					"theme":         {Type: String, Value: "dark"},
					"notifications": {Type: Null, Value: nil},
					"retryCount":    {Type: Int, Value: 3},
				},
			},
			"features": {
				Type:      Array,
				ArrayType: String,
				Value:     []*Node{{Type: String, Value: "auth"}, {Type: String, Value: "logs"}},
			},
		},
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
