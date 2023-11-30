package types

import (
	"fmt"
	"strconv"
)

var tBoolMethodsDescriptions map[string]string = map[string]string{
	"eq":       "bool.eq(arg bool) bool : Checks that the value is equal to the argument",
	"toString": "bool.toString() string : Converts the value to a string",
}

type tBool struct {
	value bool
}

func boolFactory(value interface{}) (IType, error) {
	if value, ok := value.(bool); ok {
		return &tBool{value: value}, nil
	}

	return nil, fmt.Errorf("value is not a bool")
}

func (t tBool) TypeName() string {
	return "bool"
}

func (t tBool) Value() interface{} {
	return t.value
}

func (t tBool) GetMethod(method string) Method {
	tBoolMethods := map[string]Method{
		"eq": func(args []IType) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 1 {
				return nil, fmt.Errorf("bool.eq expects 1 argument")
			}

			// Cast argument to bool type
			b, ok := args[0].(*tBool)
			if !ok {
				return nil, fmt.Errorf("bool.eq expects a bool argument")
			}

			// Check that the argument is equal to the value
			if b.value != t.value {
				return &tBool{value: false}, fmt.Errorf("bool.eq failed: %v != %v", b.value, t.value)
			}

			return &tBool{value: true}, nil
		},
		"toString": func(args []IType) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 0 {
				return nil, fmt.Errorf("bool.toString expects 0 arguments")
			}

			// Convert to string
			return &tString{value: strconv.FormatBool(t.value)}, nil
		},
	}

	// Get requested method
	if _, ok := tBoolMethods[method]; !ok {
		return func(args []IType) (IType, error) {
			return nil, fmt.Errorf("bool does not have method %s", method)
		}
	}

	return tBoolMethods[method]
}
