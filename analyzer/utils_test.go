package analyzer

import (
	"errors"
	"testing"
)

func TestSplitFileAliasAndPath(t *testing.T) {
	tests := []struct {
		input         string
		expectedAlias string
		expectedPath  string
		expectedErr   error
	}{
		{
			input:         "alias:path.to.field",
			expectedAlias: "alias",
			expectedPath:  "path.to.field",
			expectedErr:   nil,
		},
		{
			input:         "anotherAlias:path.to.another.field",
			expectedAlias: "anotherAlias",
			expectedPath:  "path.to.another.field",
			expectedErr:   nil,
		},
		{
			input:         "noColonHere",
			expectedAlias: "",
			expectedPath:  "",
			expectedErr:   errors.New("invalid field format: noColonHere"),
		},
		{
			input:         "multiple:colons:here",
			expectedAlias: "",
			expectedPath:  "",
			expectedErr:   errors.New("invalid field format: multiple:colons:here"),
		},
	}

	for _, test := range tests {
		alias, path, err := splitFileAliasAndPath(test.input)
		if alias != test.expectedAlias || path != test.expectedPath || (err != nil && err.Error() != test.expectedErr.Error()) {
			t.Errorf("SplitFileAliasAndPath(%q) = %q, %q, %v, want %q, %q, %v", test.input, alias, path, err, test.expectedAlias, test.expectedPath, test.expectedErr)
		}
	}
}
