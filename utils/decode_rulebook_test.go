package utils

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ConfigMate/configmate/analyzer"
)

// TestSimple_DecodeRulebook tests the DecodeRulebook function.
func TestSimple_DecodeRulebook(t *testing.T) {
	// Input
	var ruleBookConfig = []byte(`
	# This is a TOML configuration for the Rulebook
	name = "Sample Rulebook"
	description = "This is a sample rulebook for experimentation."
	
	# Files to be checked
	[files.file1]
	path = "/examples/configurations/sample_config.json"
	format = "json"
	
	# List of rules to be checked
	[[rules]]
	field = "file1.console.isActive"
	description = "Determines if the console is active"
	type = "bool"
	checks = ["is(true)"]
	default = "false"
	notes = "This is a note for the rule. It can be used to provide additional information about the field, the rule, or the checks being applied."
	`)

	// Test cases
	type testCase struct {
		input    []byte
		expected *analyzer.Rulebook
		err      bool
	}

	// Mock Node result
	expectedRulebook := &analyzer.Rulebook{
		Name:        "Sample Rulebook",
		Description: "This is a sample rulebook for experimentation.",
		Files: map[string]analyzer.FileDetails{ // Changed this to match actual struct
			"file1": {
				Path:   "/examples/configurations/sample_config.json",
				Format: "json",
			},
		},
		Rules: []analyzer.Rule{
			{
				Field:       "file1.console.isActive",
				Description: "Determines if the console is active",
				Type:        "bool",
				Checks:      []string{"is(true)"},
				Default:     "false",
				Notes:       "This is a note for the rule. It can be used to provide additional information about the field, the rule, or the checks being applied.",
			},
		},
	}

	testCases := []testCase{
		{
			input:    ruleBookConfig,
			expected: expectedRulebook,
			err:      false,
		},
	}

	// Run tests
	for _, test := range testCases {
		actual, err := DecodeRulebook(test.input)

		if err != nil && !test.err {
			t.Errorf("Unexpected error: %s", err)
		} else if err == nil && test.err {
			t.Errorf("Expected error, got nil")
		} else if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("Expected %v, got %v", test.expected, actual)
		}
	}
}

// TestComplex_DecodeRulebook tests the DecodeRulebook function.
func TestComplex_DecodeRulebook(t *testing.T) {
	// Input
	var ruleBookConfig = []byte(`
	# This is a TOML configuration for the Rulebook
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
	notes = "This is a note for the rule. It can be used to provide additional information about the field, the rule, or the checks being applied."

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
	checks = ["foreach().", "at()"]
	`)

	// Test cases
	type testCase struct {
		input    []byte
		expected *analyzer.Rulebook
		err      bool
	}

	// Mock Node result
	expectedRulebook := &analyzer.Rulebook{
		Name:        "Sample Rulebook",
		Description: "This is a sample rulebook for experimentation.",
		Files: map[string]analyzer.FileDetails{
			"file1": {
				Path:   "/examples/configurations/sample_config.json",
				Format: "json",
			},
			"file2": {
				Path:   "/examples/configurations/sample_config.toml",
				Format: "toml",
			},
		},
		Rules: []analyzer.Rule{
			{
				Field:       "file1.console.isActive",
				Description: "Determines if the console is active",
				Type:        "bool",
				Checks:      []string{"is(true)"},
				Default:     "false",
				Notes:       "This is a note for the rule. It can be used to provide additional information about the field, the rule, or the checks being applied.",
			},
			{
				Field:       "file1.proxy.bindPort",
				Description: "Determines the port the proxy is bound to",
				Type:        "port",
				Checks:      []string{"range(1009, 3280)"},
				Default:     "1009",
			},
			{
				Field:       "file1.proxy.bindAddress",
				Description: "Determines the address the proxy is bound to",
				Type:        "host",
				Checks:      []string{"is(some.host.com)", "reachable()", "join(file1.proxy.bindPort).listening()"},
				Default:     "some.host.com",
			},
			{
				Field:       "file1.proxy.certification",
				Description: "Details of the certification",
				Type:        "object",
				Optional:    true,
			},
			{
				Field:       "file1.proxy.certification.cert",
				Description: "The certification file",
				Type:        "file",
				Checks:      []string{"exists()"},
				Default:     "/path/to/cert.pem",
			},
			{
				Field:       "file1.proxy.certification.key",
				Description: "The key file",
				Type:        "file",
				Checks:      []string{"exists()"},
				Default:     "/path/to/key.pem",
			},
			{
				Field:       "dns_servers",
				Description: "List of DNS servers",
				Type:        "list:host",
				Checks:      []string{"foreach().", "at()"},
			},
		},
	}

	testCases := []testCase{
		{
			input:    ruleBookConfig,
			expected: expectedRulebook,
			err:      false,
		},
	}

	// Run tests
	for _, test := range testCases {
		actual, err := DecodeRulebook(test.input)

		if err != nil && !test.err {
			t.Errorf("Unexpected error: %s", err)
		} else if err == nil && test.err {
			t.Errorf("Expected error, got nil")
		} else if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("Expected %v, got %v", test.expected, actual)
		}
	}
}

// TestWrongFormatError_DecodeRulebook tests the DecodeRulebook function with rulebook using json format instead of toml.
func TestWrongFormatError_DecodeRulebook(t *testing.T) {
	// Input
	var ruleBookConfig = []byte(`{
		"name": "sample",
		"version": 1.3,
		"active": true,
		"features": ["auth", "logs"]
	}
	`)

	// Test cases
	type testCase struct {
		input    []byte
		expected error
		err      bool
	}

	// Mock Node result
	testCases := []testCase{
		{
			input:    ruleBookConfig,
			expected: fmt.Errorf(`error decoding file into a rulebook object: toml: line 1: expected '.' or '=', but got '{' instead`),
			err:      true,
		},
	}

	// Run tests
	for _, test := range testCases {
		_, err := DecodeRulebook(test.input)

		if err != nil && !test.err {
			t.Errorf("Unexpected error: %s", err)
		} else if err == nil && test.err {
			t.Errorf("Expected error, got nil")
		} else if err != nil && test.err && err.Error() != test.expected.Error() {
			t.Errorf("Expected %v, got %v", test.expected, err)
		}
	}
}

// TestWrongSyntaxError_DecodeRulebook tests the DecodeRulebook function with rulebook having wrong syntaxs.
func TestWrongSyntaxError_DecodeRulebook(t *testing.T) {
	// Input
	var ruleBookConfig = []byte(`
		# This is a TOML configuration for the Rulebook
		name = "Sample Rulebook"
		description = "This is a sample rulebook for experimentation."

		# Files to be checked
		[files.file1]
		path = "/examples/configurations/sample_config.json"
		format = "json

		[files.file2]
		path = "/examples/configurations/sample_config.toml"
		format = "toml"

		[[rules
		field = "file1.proxy.bindPort"
		description = "Determines the port the proxy is bound to"
		type = "port"
		checks = ["range(1009, 3280)"]
		default = "1009"
	`)

	// Test cases
	type testCase struct {
		input    []byte
		expected error
		err      bool
	}

	// Mock Node result
	testCases := []testCase{
		{
			input:    ruleBookConfig,
			expected: fmt.Errorf(`error decoding file into a rulebook object: toml: line 9 (last key "files.file1.format"): strings cannot contain newlines`),
			err:      true,
		},
	}

	// Run tests
	for _, test := range testCases {
		_, err := DecodeRulebook(test.input)

		if err != nil && !test.err {
			t.Errorf("Unexpected error: %s", err)
		} else if err == nil && test.err {
			t.Errorf("Expected error, got nil")
		} else if err != nil && test.err && err.Error() != test.expected.Error() {
			t.Errorf("Expected %v, got %v", test.expected, err)
		}
	}
}
