package analyzer

import (
	"testing"

	"github.com/ConfigMate/configmate/parsers"
)

// TestSplitKey tests the splitKey function.
func TestSplitKey(t *testing.T) {
	// Test cases
	type testCase struct {
		key      string
		expected []string
	}

	testCases := []testCase{
		{
			key:      "server.port",
			expected: []string{"server", "port"},
		},
		{
			key:      "settings.users[0].name",
			expected: []string{"settings", "users", "0", "name"},
		},
		{
			key:      "logLevel",
			expected: []string{"logLevel"},
		},
		{
			key:      "[3].name",
			expected: []string{"3", "name"},
		},
	}

	// Run tests
	for _, test := range testCases {
		actual := splitKey(test.key)
		if len(actual) != len(test.expected) {
			t.Errorf("splitKey(%s) returned %+v, expected %+v", test.key, actual, test.expected)
			continue
		}

		for i, segment := range actual {
			if segment != test.expected[i] {
				t.Errorf("splitKey(%s) returned %+v, expected %+v", test.key, actual, test.expected)
			}
		}
	}
}

// TestGetNodeFromConfigFile tests the getNodeFromConfigFile function.
func TestGetNodeFromConfigFile(t *testing.T) {
	// Test cases
	type testCase struct {
		configFile parsers.ConfigFile
		key        string
		expected   interface{}
		err        bool
	}

	testCases := []testCase{
		{
			configFile: &parsers.Node{
				Type: parsers.Object,
				Value: map[string]*parsers.Node{
					"server": {
						Type: parsers.Object,
						Value: map[string]*parsers.Node{
							"port": {
								Type:  parsers.Int,
								Value: 8080,
							},
						},
					},
				},
			},
			key:      "server.port",
			expected: &parsers.Node{Type: parsers.Int, Value: "8080"},
			err:      false,
		},
		{
			configFile: &parsers.Node{
				Type: parsers.Object,
				Value: map[string]*parsers.Node{
					"server": {
						Type:  parsers.Object,
						Value: map[string]*parsers.Node{},
					},
				},
			},
			key:      "server.port",
			expected: nil,
			err:      true,
		},
		{
			configFile: &parsers.Node{
				Type: parsers.Object,
				Value: map[string]*parsers.Node{
					"server": {
						Type: parsers.Object,
						Value: map[string]*parsers.Node{
							"port": {
								Type:  parsers.String,
								Value: "8080",
							},
						},
					},
				},
			},
			key:      "server.port[0]",
			expected: nil,
			err:      true,
		},
		{
			configFile: &parsers.Node{
				Type: parsers.Object,
				Value: map[string]*parsers.Node{
					"server": {
						Type: parsers.Object,
						Value: map[string]*parsers.Node{
							"dns_servers": {
								Type: parsers.Array,
								Value: []*parsers.Node{
									{
										Type:  parsers.String,
										Value: "some.dns.server",
									},
									{
										Type:  parsers.String,
										Value: "some.other.dns.server",
									},
								},
							},
						},
					},
				},
			},
			key:      "server.dns_servers[1]",
			expected: &parsers.Node{Type: parsers.String, Value: "some.other.dns.server"},
			err:      false,
		},
	}

	// Run tests
	for _, test := range testCases {
		actual, err := getNodeFromConfigFileNode(test.configFile, test.key)
		if err != nil && !test.err {
			t.Errorf("getValueFromConfigFile(%+v, %s) returned error %s, expected no error", test.configFile, test.key, err.Error())
		} else if err == nil && test.err {
			t.Errorf("getValueFromConfigFile(%+v, %s) returned no error, expected error", test.configFile, test.key)
		} else if actual != test.expected {
			t.Errorf("getValueFromConfigFile(%+v, %s) returned %+v, expected %+v", test.configFile, test.key, actual, test.expected)
		}
	}
}

// TestDecodeFileValue tests the decodeFileValue function.
func TestDecodeFileValue(t *testing.T) {
	type testCase struct {
		input string
		alias string
		key   string
	}

	tests := []testCase{
		{"file_alias.server.port", "file_alias", "server.port"},
		{"file_alias.settings.users[0].name", "file_alias", "settings.users[0].name"},
	}

	for _, test := range tests {
		alias, key := decodeFileValue(test.input)
		if alias != test.alias || key != test.key {
			t.Errorf("decodeFileValue(%q) = (%v, %q), want (%v, %q)", test.input, alias, key, test.alias, test.key)
		}
	}
}

// TestEqualType tests the equalType function.
func TestEqualType(t *testing.T) {
	// Test cases
	type testCase struct {
		nodeType parsers.FieldType
		argType  CheckArgType
		expected bool
	}

	testCases := []testCase{
		{
			nodeType: parsers.Int,
			argType:  Int,
			expected: true,
		},
		{
			nodeType: parsers.Float,
			argType:  Float,
			expected: true,
		},
		{
			nodeType: parsers.Bool,
			argType:  Bool,
			expected: true,
		},
		{
			nodeType: parsers.String,
			argType:  String,
			expected: true,
		},
		{
			nodeType: parsers.Int,
			argType:  Float,
			expected: false,
		},
		{
			nodeType: parsers.Float,
			argType:  Int,
			expected: false,
		},
		{
			nodeType: parsers.Bool,
			argType:  String,
			expected: false,
		},
		{
			nodeType: parsers.String,
			argType:  Bool,
			expected: false,
		},
	}

	// Run tests
	for _, test := range testCases {
		actual := equalType(test.nodeType, test.argType)
		if actual != test.expected {
			t.Errorf("equalType(%d, %d) returned %t, expected %t", test.nodeType, test.argType, actual, test.expected)
		}
	}
}

// TestInterpretLiteralOutput tests the interpretLiteralOutput function.
func TestInterpretLiteralOutput(t *testing.T) {
	// Test cases
	type testCase struct {
		argType  CheckArgType
		argValue string
		expected interface{}
	}

	testCases := []testCase{
		{
			argType:  Int,
			argValue: "8080",
			expected: 8080,
		},
		{
			argType:  Int,
			argValue: "8080.0",
			expected: nil,
		},
		{
			argType:  Float,
			argValue: "8080.5",
			expected: 8080.5,
		},
		{
			argType:  Float,
			argValue: "8080",
			expected: 8080.0,
		},
		{
			argType:  Bool,
			argValue: "true",
			expected: true,
		},
		{
			argType:  Bool,
			argValue: "false",
			expected: false,
		},
		{
			argType:  Bool,
			argValue: "8080",
			expected: nil,
		},
		{
			argType:  String,
			argValue: "8080",
			expected: "8080",
		},
	}

	// Run tests
	for _, test := range testCases {
		actual, _ := interpretLiteralOutput(test.argType, test.argValue)
		if actual != test.expected {
			t.Errorf("interpretLiteralOutput(%d, %s) returned %+v, expected %+v", test.argType, test.argValue, actual, test.expected)
		}
	}
}
