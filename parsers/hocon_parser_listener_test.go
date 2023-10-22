package parsers

import (
	"fmt"
	"testing"
)

// Test_HOCON_Parse tests the Parse function of a *JsonParser using a hocon config.
func Test_HOCON_Parse(t *testing.T) {
	// Input
	var hoconConfig = []byte(`
		# Server settings
		server {
			host = "localhost"
			port = 8080
			ssl {
				enabled = true
				cert = "/path/to/cert.pem"
				key = "/path/to/key.pem"
			}
		}

		# Database settings
		database {
			name = "mydb"
			user = "myuser"
			password = "mypassword"
			pool {
				max_connections = 10
				idle_timeout = 5m
			}
		}

		# Logging settings
		logging {
			level = "debug"
			format = "json"
			output {
				file {
					path = "/var/log/myapp.log"
					max_size = 10MB
					max_age = 7d
				}
				console {
					enabled = true
					color = true
				}
			}
		}
	`)

	// Test case
	type testCase struct {
		input    []byte
		expected error
		err      bool
	}

	// Mock Node result
	testCases := []testCase{
		{
			input:    hoconConfig,
			expected: nil,
			err:      false,
		},
	}

	// Run test cases
	for _, test := range testCases {
		parser := &HoconParser{}
		_, err := parser.Parse(test.input)

		if err != nil && !test.err {
			t.Errorf("Unexpected error: %s", err)
		} else if err == nil && test.err {
			t.Errorf("Expected error, got nil")
		} else if err != nil && test.err && err.Error() != test.expected.Error() {
			t.Errorf("Expected %v, got %v", test.expected, err)
		}
	}
}

// TestSingleError_Parse tests the Parse function of a *JsonParser using a wrong hocon config.
func TestSyntaxError_HOCON_Parse(t *testing.T) {
	// Input
	var hoconConfig = []byte(`{
		"name": "sample",
		"version" 1.3
		"active": true,
		"features": ["auth", "logs"]
	`)

	// Test case
	type testCase struct {
		input    []byte
		expected error
		err      bool
	}

	// Mock Node result
	testCases := []testCase{
		{
			input:    hoconConfig,
			expected: fmt.Errorf(`Syntax errors: [line 3:12 no viable alternative at input '"version"1.3' line 6:1 extraneous input '<EOF>' expecting {',', '}', STRING, PATHELEMENT}]`),
			err:      true,
		},
	}

	// Run test cases
	for _, test := range testCases {
		parser := &HoconParser{}
		_, err := parser.Parse(test.input)

		if err != nil && !test.err {
			t.Errorf("Unexpected error: %s", err)
		} else if err == nil && test.err {
			t.Errorf("Expected error, got nil")
		} else if err != nil && test.err && err.Error() != test.expected.Error() {
			t.Errorf("Expected %v, got %v", test.expected, err)
		}
	}
}
