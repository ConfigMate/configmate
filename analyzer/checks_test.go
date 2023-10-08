package analyzer

import (
	"reflect"
	"testing"
)

// TestParseCheckArg tests the parseCheckArg function.
func TestParseCheckArg(t *testing.T) {
	// Test cases
	type testCase struct {
		arg      string
		expected CheckArg
		err      bool
		errMsg   string
	}

	testCases := []testCase{
		{
			arg: "f:i:test.server.port",
			expected: CheckArg{
				v: FileValue{
					alias: "test",
					path:  "server.port",
				},
				s: File,
				t: Int,
			},
			err: false,
		},
		{
			arg: "l:s:hello",
			expected: CheckArg{
				v: "hello",
				s: Literal,
				t: String,
			},
			err: false,
		},
		{
			arg:      "f:z:test.server.port",
			expected: CheckArg{},
			err:      true,
			errMsg:   "invalid type in rule argument",
		},
		{
			arg:      "l:hello",
			expected: CheckArg{},
			err:      true,
			errMsg:   "invalid number of segments in rule argument",
		},
		{
			arg:      "z:i:test.server.port",
			expected: CheckArg{},
			err:      true,
			errMsg:   "invalid source in rule argument",
		},
		{
			arg: "l:s:",
			expected: CheckArg{
				v: "",
				s: Literal,
				t: String,
			},
			err: false,
		},
		{
			arg:      "l:i:",
			expected: CheckArg{},
			err:      true,
			errMsg:   "value cannot be empty in rule argument",
		},
		{
			arg: "l:f:3.14",
			expected: CheckArg{
				v: 3.14,
				s: Literal,
				t: Float,
			},
			err: false,
		},
		{
			arg: "l:b:true",
			expected: CheckArg{
				v: true,
				s: Literal,
				t: Bool,
			},
			err: false,
		},
		{
			arg: "l:s:false",
			expected: CheckArg{
				v: "false",
				s: Literal,
				t: String,
			},
			err: false,
		},
		{
			arg:      "l:i:3.14",
			expected: CheckArg{},
			err:      true,
			errMsg:   "failed to decode literal value",
		},
		{
			arg:      "f:i:test",
			expected: CheckArg{},
			err:      true,
			errMsg:   "invalid number of segments in file value",
		},
		{
			arg:      "f:i:.something",
			expected: CheckArg{},
			err:      true,
			errMsg:   "invalid file value",
		},
	}

	// Run tests
	for _, test := range testCases {
		actual, err := ParseCheckArg(test.arg)
		if test.err && err == nil {
			t.Errorf("parseCheckArg(%s) returned no error, expected error", test.arg)
			continue
		} else if !test.err && err != nil {
			t.Errorf("parseCheckArg(%s) returned error %s, expected no error", test.arg, err)
			continue
		}

		if err == nil {
			if !reflect.DeepEqual(actual.v, test.expected.v) {
				t.Errorf("parseCheckArg(%s) returned %+v, expected %+v", test.arg, actual, test.expected)
			}
			if actual.s != test.expected.s {
				t.Errorf("parseCheckArg(%s) returned %+v, expected %+v", test.arg, actual, test.expected)
			}
			if actual.t != test.expected.t {
				t.Errorf("parseCheckArg(%s) returned %+v, expected %+v", test.arg, actual, test.expected)
			}
		}
	}
}
