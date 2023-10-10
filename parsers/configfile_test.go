package parsers

import (
	"reflect"
	"testing"
)

// TestNode_Get tests the Get function of a *Node.
func TestNode_Get(t *testing.T) {
	// Test cases
	type testCase struct {
		configFile *Node
		path       string
		expected   *Node
		err        bool
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
			path:     "server.port",
			expected: &Node{Type: Int, Value: 8080},
			err:      false,
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
			path:     "server.port",
			expected: nil,
			err:      true,
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
			path:     "server.port[0]",
			expected: nil,
			err:      true,
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
			path:     "server.dns_servers[1]",
			expected: &Node{Type: String, Value: "some.other.dns.server"},
			err:      false,
		},
	}

	// Run tests
	for _, test := range testCases {
		actual, err := test.configFile.Get(test.path)
		if err != nil && !test.err {
			t.Errorf("getValueFromConfigFile(%+v, %s) returned error %s, expected no error", test.configFile, test.path, err.Error())
		} else if err == nil && test.err {
			t.Errorf("getValueFromConfigFile(%+v, %s) returned no error, expected error", test.configFile, test.path)
		} else if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("getValueFromConfigFile(%+v, %s) returned %+v, expected %+v", test.configFile, test.path, actual, test.expected)
		}
	}
}

// TestNode_Get_SyntacticallyWrong tests the Get function of a *Node for syntactically incorrect content.
func TestNode_Get_SyntacticallyWrong(t *testing.T) {
	type testCase struct {
		configFile *Node
		path       string
		err        bool
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
			path: "server..port",  // Double dots should trigger an error
			err:  true,
		},
		{
			configFile: &Node{
				Type: Object,
				Value: map[string]*Node{},
			},
			path: "server[port]",  // Malformed index should trigger an error
			err:  true,
		},
		{
			configFile: &Node{
				Type: Object,
				Value: map[string]*Node{},
			},
			path: "",  // Empty path should trigger an error
			err:  true,
		},
	}

	for _, test := range testCases {
		_, err := test.configFile.Get(test.path)
		if err == nil && test.err {
			t.Errorf("Get(%s) returned no error, expected error", test.path)
		} else if err != nil && !test.err {
			t.Errorf("Get(%s) returned error %s, expected no error", test.path, err.Error())
		}
	}
}

func Test_SplitKey(t *testing.T) {
    // Test cases
    type testCase struct {
        path     string
        expected []string
    }

    testCases := []testCase{
        {
            path:     "server.port",
            expected: []string{"server", "port"},
        },
        {
            path:     "settings.users[0].name",
            expected: []string{"settings", "users", "0", "name"},
        },
        {
            path:     "logLevel",
            expected: []string{"logLevel"},
        },
    }

    // Run tests
    for _, test := range testCases {
        actual := splitKey(test.path)
        if !reflect.DeepEqual(actual, test.expected) {
            t.Errorf("splitKey(%s) returned %+v, expected %+v", test.path, actual, test.expected)
        }
    }
}