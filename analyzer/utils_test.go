package analyzer

import (
	"testing"

	"github.com/ConfigMate/configmate/parsers"
)

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
