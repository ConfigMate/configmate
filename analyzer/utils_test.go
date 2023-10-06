package analyzer

import (
	"testing"

	"github.com/ConfigMate/configmate/parsers"
)

// TestGetFileFormat tests the getFileFormat function.
func TestGetFileFormat(t *testing.T) {
	// Test cases
	type testCase struct {
		filename string
		expected FileFormat
	}

	testCases := []testCase{
		{
			filename: "test.hocon",
			expected: HOCON,
		},
		{
			filename: "test.json",
			expected: JSON,
		},
		{
			filename: "test.toml",
			expected: TOML,
		},
		{
			filename: "test.yaml",
			expected: YAML,
		},
		{
			filename: "test.yml",
			expected: YAML,
		},
		{
			filename: "test",
			expected: Unknown,
		},
	}

	// Run tests
	for _, test := range testCases {
		actual := getFileFormat(test.filename)
		if actual != test.expected {
			t.Errorf("getFileFormat(%s) returned %d, expected %d", test.filename, actual, test.expected)
		}
	}
}

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
