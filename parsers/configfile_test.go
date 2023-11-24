package parsers

import (
	"fmt"
	"reflect"
	"testing"
)

// TestNode_Get tests the Get function of a *Node.
func TestNode_Get(t *testing.T) {
	// Test cases
	type testCase struct {
		configFile  *Node
		key         *NodeKey
		expected    *Node
		expectedErr error
	}

	testCases := []testCase{
		{
			configFile: &Node{
				Type: Object,
				Value: map[string]*Node{
					"server": {
						Type: Object,
						Value: map[string]*Node{
							"port": {
								Type:  Int,
								Value: 8080,
							},
						},
					},
				},
			},
			key:         &NodeKey{Segments: []string{"server", "port"}},
			expected:    &Node{Type: Int, Value: 8080},
			expectedErr: nil,
		},
		{
			configFile: &Node{
				Type: Object,
				Value: map[string]*Node{
					"server": {
						Type:  Object,
						Value: map[string]*Node{},
					},
				},
			},
			key:         &NodeKey{Segments: []string{"server", "port"}},
			expected:    nil,
			expectedErr: nil,
		},
		{
			configFile: &Node{
				Type: Object,
				Value: map[string]*Node{
					"server": {
						Type: Object,
						Value: map[string]*Node{
							"port. test": {
								Type:  String,
								Value: "8080",
							},
						},
					},
				},
			},
			key:         &NodeKey{Segments: []string{"server", "port. test", "something"}},
			expected:    nil,
			expectedErr: fmt.Errorf("cannot traverse leaf node in path server.'port. test'.something"),
		},
	}

	// Run tests
	for _, test := range testCases {
		actual, err := test.configFile.Get(test.key)
		if err != nil && test.expectedErr == nil {
			t.Errorf("getValueFromConfigFile(%+v, %s) returned error %s, expected no error", test.configFile, test.key.String(), err.Error())
		} else if err == nil && test.expectedErr != nil {
			t.Errorf("getValueFromConfigFile(%+v, %s) returned no error, expected error", test.configFile, test.key.String())
		} else if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("getValueFromConfigFile(%+v, %s) returned %+v, expected %+v", test.configFile, test.key.String(), actual, test.expected)
		}
	}
}
