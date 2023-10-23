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
		path        string
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
			path:        "server.port",
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
			path:        "server.port",
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
							"port": {
								Type:  String,
								Value: "8080",
							},
						},
					},
				},
			},
			path:        "server.port[0]",
			expected:    nil,
			expectedErr: fmt.Errorf("cannot traverse leaf node port in path server.port"),
		},
		{
			configFile: &Node{
				Type: Object,
				Value: map[string]*Node{
					"server": {
						Type: Object,
						Value: map[string]*Node{
							"dns_servers": {
								Type: Array,
								Value: []*Node{
									{
										Type:  String,
										Value: "some.dns.server",
									},
									{
										Type:  String,
										Value: "some.other.dns.server",
									},
								},
							},
						},
					},
				},
			},
			path:        "server.dns_servers[1]",
			expected:    &Node{Type: String, Value: "some.other.dns.server"},
			expectedErr: nil,
		},
		{
			configFile: &Node{
				Type: Object,
				Value: map[string]*Node{
					"dns_servers": {
						Type: Array,
						Value: []*Node{
							{
								Type:  String,
								Value: "some.dns.server",
							},
						},
					},
				},
			},
			path:        "dns_servers[3.14]",
			expected:    nil,
			expectedErr: fmt.Errorf("failed to convert [3.14] to int in path server.dns_servers[3.14]"),
		},
		{
			configFile: &Node{
				Type: Object,
				Value: map[string]*Node{
					"dns_servers": {
						Type: Array,
						Value: []*Node{
							{
								Type:  String,
								Value: "some.dns.server",
							},
						},
					},
				},
			},
			path:        "dns_servers[something]",
			expected:    nil,
			expectedErr: fmt.Errorf("failed to convert [something] to int in path server.dns_servers[3.14]"),
		},
	}

	// Run tests
	for _, test := range testCases {
		actual, err := test.configFile.Get(test.path)
		if err != nil && test.expectedErr == nil {
			t.Errorf("getValueFromConfigFile(%+v, %s) returned error %s, expected no error", test.configFile, test.path, err.Error())
		} else if err == nil && test.expectedErr != nil {
			t.Errorf("getValueFromConfigFile(%+v, %s) returned no error, expected error", test.configFile, test.path)
		} else if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("getValueFromConfigFile(%+v, %s) returned %+v, expected %+v", test.configFile, test.path, actual, test.expected)
		}
	}
}
