package utils

import (
	"reflect"
	"testing"
)

// Test_ReadFile tests the ReadFile function.
func Test_ReadFile(t *testing.T) {
	// Input
	var filePath = "../examples/rulebooks/sample0.cmrb"

	// Test cases
	type testCase struct {
		input    string
		expected string
		err      bool
	}

	// Mock Node result
	expectedData := `# This is a TOML configuration for the Rulebook
name = "Sample Rulebook"
description = "This is a sample rulebook for experimentation."

# Files to be checked
[files.file1]
path = "/examples/configurations/sample_config.json"
format = "json"

[files.file2]
path = "/examples/configurations/sample_config.toml"
format = "toml"

# List of rules to be checked
[[rules]]
field = "file1.console.isActive"
description = "Determines if the console is active"
type = "bool"
checks = ["is(true)"]
default = "false"
notes = """
This is a note for the rule. It can be used to provide additional information
about the field, the rule, or the checks being applied.
"""

[[rules]]
field = "file1.proxy.bindPort"
description = "Determines the port the proxy is bound to"
type = "port"
checks = ["range(1009, 3280)"]
default = "1009"

[[rules]]
field = "file1.proxy.bindAddress"
description = "Determines the address the proxy is bound to"
type = "host"
checks = ["is(some.host.com)", "reachable()", "join(file1.proxy.bindPort).listening()"]
default = "some.host.com"

[[rules]]
field = "file1.proxy.certification"
description = "Details of the certification"
type = "object"
optional = true

[[rules]]
field = "file1.proxy.certification.cert"
description = "The certification file"
type = "file"
checks = ["exists()"]
default = "/path/to/cert.pem"

[[rules]]
field = "file1.proxy.certification.key"
description = "The key file"
type = "file"
checks = ["exists()"]
default = "/path/to/key.pem"

[[rules]]
field = "dns_servers"
description = "List of DNS servers"
type = "list:host"
checks = ["foreach().", "at()"]`

	testCases := []testCase{
		{
			input:    filePath,
			expected: expectedData,
			err:      false,
		},
	}

	// Run tests
	for _, test := range testCases {
		actual, err := ReadFile(test.input)

		if err != nil && !test.err {
			t.Errorf("Unexpected error: %s", err)
		} else if err == nil && test.err {
			t.Errorf("Expected error, got nil")
		} else if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("Expected %v, got %v", test.expected, string(actual))
		}
	}
}
