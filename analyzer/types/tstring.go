package types

import (
	"fmt"
	"regexp"
)

type tString struct {
	value string
}

func stringFactory(value interface{}) (IType, error) {
	if value, ok := value.(string); ok {
		return &tString{value: value}, nil
	}

	return nil, fmt.Errorf("value %v is not a string", value)
}

func (t tString) TypeName() string {
	return "string"
}

func (t tString) Value() interface{} {
	return t.value
}

func (t tString) Checks() map[string]Check {
	return map[string]Check{
		"eq": func(args []IType) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 1 {
				return nil, fmt.Errorf("string.eq expects 1 argument")
			}

			// Cast argument to string type
			s, ok := args[0].(*tString)
			if !ok {
				return nil, fmt.Errorf("string.eq expects a string argument")
			}

			// Check that the argument is equal to the value
			if s.value != t.value {
				return &tBool{value: false}, fmt.Errorf("string.eq failed: %v != %v", s.value, t.value)
			}

			return &tBool{value: true}, nil
		},
		"regex": func(args []IType) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 1 {
				return nil, fmt.Errorf("string.regex expects 1 argument")
			}

			// Cast pattern argument to string type
			pattern, ok := args[0].(*tString)
			if !ok {
				return nil, fmt.Errorf("string.regex expects a string argument")
			}

			// Compile the regular expression
			re, err := regexp.Compile(pattern.value)
			if err != nil {
				return nil, fmt.Errorf("string.regex failed to compile pattern: %v", err)
			}

			// Check that the value matches the regular expression
			if !re.MatchString(t.value) {
				return &tBool{value: false}, fmt.Errorf("string.regex failed: %v does not match pattern %v", t.value, pattern)
			}

			return &tBool{value: true}, nil
		},
	}
}

func (t tString) ChecksDescription() map[string]string {
	return map[string]string{
		"eq":    "string.eq(s string) : Checks that the value is equal to s",
		"regex": "string.regex(pattern string) : Checks that the value matches the pattern",
	}
}
