package types

import (
	"fmt"
	"strconv"
)

type tBool struct {
	value bool
}

func boolFactory(value interface{}) (IType, error) {
	if value, ok := value.(bool); ok {
		return &tBool{value: value}, nil
	}

	return nil, fmt.Errorf("value %v is not a bool", value)
}

func (t *tBool) Value() interface{} {
	return t.value
}

func (t *tBool) Checks() map[string]Check {
	return map[string]Check{
		"eq": func(args ...interface{}) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 1 {
				return nil, fmt.Errorf("bool.eq expects 1 argument")
			}

			// Check that the argument is a bool
			if _, ok := args[0].(bool); !ok {
				return nil, fmt.Errorf("bool.eq expects a bool argument")
			}

			// Check that the argument is equal to the value
			if args[0].(bool) != t.value {
				return nil, fmt.Errorf("bool.eq failed: %v != %v", args[0].(bool), t.value)
			}

			return &tBool{value: true}, nil
		},
		"toString": func(args ...interface{}) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 0 {
				return nil, fmt.Errorf("bool.toString expects 0 arguments")
			}

			// Convert to string
			return &tString{value: strconv.FormatBool(t.value)}, nil
		},
	}
}

func (t *tBool) ChecksDescription() map[string]string {
	return map[string]string{
		"eq":       "bool.eq(arg bool) : Checks that the value is equal to the argument",
		"toString": "bool.toString() : Converts the value to a string",
	}
}
