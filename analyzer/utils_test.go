package analyzer

import (
	"strings"
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

// splitRuleFileArgument splits the argument of a rule that
// references a field in a file into file and key.
// The argument looks like this: "file_alias:server.port", the
// returned values will be "file_alias" and "server.port".
func splitRuleFileArgument(arg string) (string, string) {
	// Split the argument based on the colon
	segments := strings.Split(arg, ":")
	if len(segments) == 1 {
		return "", segments[0]
	}

	return segments[0], segments[1]
}

// TestGetValueFromConfigFile tests the getValueFromConfigFile function.
func TestGetValueFromConfigFile(t *testing.T) {
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
								Type:  parsers.String,
								Value: "8080",
							},
						},
					},
				},
			},
			key:      "server.port",
			expected: "8080",
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
			expected: "some.other.dns.server",
			err:      false,
		},
	}

	// Run tests
	for _, test := range testCases {
		actual, err := getValueFromConfigFile(test.configFile, test.key)
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
		input  string
		isFile bool
		alias  string
		key    string
	}

	tests := []testCase{
		{"f:file_alias.server.port", true, "file_alias", "server.port"},
		{"f:file_alias.settings.users[0].name", true, "file_alias", "settings.users[0].name"},
		{"server.port", false, "", ""},
		{"f:file_alias", false, "", ""},   // Missing dot after alias
		{"f:.server.port", false, "", ""}, // Missing alias before dot
	}

	for _, test := range tests {
		isFile, alias, key := decodeFileValue(test.input)
		if isFile != test.isFile || alias != test.alias || key != test.key {
			t.Errorf("decodeFileValue(%q) = (%v, %q, %q), want (%v, %q, %q)", test.input, isFile, alias, key, test.isFile, test.alias, test.key)
		}
	}
}

// TestDecodeLiteralValue tests the decodeLiteralValue function.
func TestDecodeLiteralValue(t *testing.T) {
	type testCase struct {
		input     string
		isLiteral bool
		value     string
	}

	tests := []testCase{
		{"l:100", true, "100"},
		{"l:hello world", true, "hello world"},
		{"l:true", true, "true"},
		{"100", false, ""},
		{"l:", true, ""},           // Edge case: Literal with empty value
		{"l:l:100", true, "l:100"}, // Nested "l:"
	}

	for _, tt := range tests {
		isLiteral, value := decodeLiteralValue(tt.input)
		if isLiteral != tt.isLiteral || value != tt.value {
			t.Errorf("decodeLiteralValue(%q) = (%v, %q), want (%v, %q)", tt.input, isLiteral, value, tt.isLiteral, tt.value)
		}
	}
}
