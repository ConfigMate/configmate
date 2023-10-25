package types

import (
	"fmt"
)

type tFloat struct {
	value float64
}

func floatFactory(value interface{}) (IType, error) {
	if value, ok := value.(float64); ok {
		return &tFloat{value: value}, nil
	}

	return nil, fmt.Errorf("value %v is not a float", value)
}

func (t tFloat) TypeName() string {
	return "float"
}

func (t tFloat) Value() interface{} {
	return t.value
}

func (t tFloat) Methods() []string {
	return []string{
		"eq",
		"gt",
		"gte",
		"lt",
		"lte",
		"range",
		"toInt",
		"toString",
	}
}

func (t tFloat) MethodDescription(method string) string {
	tFloatMethodDescriptions := map[string]string{
		"eq":       "float.eq(arg float) : Checks that the value is equal to the argument",
		"gt":       "float.gt(arg float) : Checks that the value is greater than the argument",
		"gte":      "float.gte(arg float) : Checks that the value is greater than or equal to the argument",
		"lt":       "float.lt(arg float) : Checks that the value is less than the argument",
		"lte":      "float.lte(arg float) : Checks that the value is less than or equal to the argument",
		"range":    "float.range(min float, max float) : Checks that the value is within the range",
		"toInt":    "float.toInt() : Converts the value to an int",
		"toString": "float.toString() : Converts the value to a string",
	}

	return tFloatMethodDescriptions[method]
}

func (t tFloat) GetMethod(method string) Method {
	tFloatMethods := map[string]Method{
		"eq": func(args []IType) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 1 {
				return nil, fmt.Errorf("float.eq expects 1 argument")
			}

			// Cast argument to float type
			f, ok := args[0].(*tFloat)
			if !ok {
				return nil, fmt.Errorf("float.eq expects a float argument")
			}

			// Check that the argument is equal to the value
			if f.value != t.value {
				return &tBool{value: false}, fmt.Errorf("float.eq failed: %v != %v", f.value, t.value)
			}
			return &tBool{value: true}, nil
		},
		"gt": func(args []IType) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 1 {
				return nil, fmt.Errorf("float.gt expects 1 argument")
			}

			// Cast argument to float type
			f, ok := args[0].(*tFloat)
			if !ok {
				return nil, fmt.Errorf("float.gt expects a float argument")
			}

			// Check that the value is greater than the argument
			if t.value <= f.value {
				return &tBool{value: false}, fmt.Errorf("float.gt failed: %v <= %v", t.value, f.value)
			}
			return &tBool{value: true}, nil
		},
		"gte": func(args []IType) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 1 {
				return nil, fmt.Errorf("float.gte expects 1 argument")
			}

			// Cast argument to float type
			f, ok := args[0].(*tFloat)
			if !ok {
				return nil, fmt.Errorf("float.gte expects a float argument")
			}

			// Check that the value is greater than or equal to the argument
			if t.value < f.value {
				return &tBool{value: false}, fmt.Errorf("float.gte failed: %v < %v", t.value, f.value)
			}
			return &tBool{value: true}, nil
		},
		"lt": func(args []IType) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 1 {
				return nil, fmt.Errorf("float.lt expects 1 argument")
			}

			// Cast argument to float type
			f, ok := args[0].(*tFloat)
			if !ok {
				return nil, fmt.Errorf("float.lt expects a float argument")
			}

			// Check that the value is less than the argument
			if t.value >= f.value {
				return &tBool{value: false}, fmt.Errorf("float.lt failed: %v >= %v", t.value, f.value)
			}
			return &tBool{value: true}, nil
		},
		"lte": func(args []IType) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 1 {
				return nil, fmt.Errorf("float.lte expects 1 argument")
			}

			// Cast argument to float type
			f, ok := args[0].(*tFloat)
			if !ok {
				return nil, fmt.Errorf("float.lte expects a float argument")
			}

			// Check that the value is less than or equal to the argument
			if t.value > f.value {
				return &tBool{value: false}, fmt.Errorf("float.lte failed: %v > %v", t.value, f.value)
			}
			return &tBool{value: true}, nil
		},
		"range": func(args []IType) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 2 {
				return nil, fmt.Errorf("float.range expects 2 arguments")
			}

			// Cast arguments to float type
			f0, ok := args[0].(*tFloat)
			if !ok {
				return nil, fmt.Errorf("float.range expects a float argument")
			}
			f1, ok := args[1].(*tFloat)
			if !ok {
				return nil, fmt.Errorf("float.range expects a float argument")
			}

			// Check that the value is within the range
			if t.value < f0.value || t.value > f1.value {
				return &tBool{value: false}, fmt.Errorf("float.range failed: %v not in [%v, %v]", t.value, f0.value, f1.value)
			}

			return &tBool{value: true}, nil
		},
		"toInt": func(args []IType) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 0 {
				return nil, fmt.Errorf("float.toInt expects 0 arguments")
			}

			// Convert to int
			return &tInt{value: int(t.value)}, nil
		},
		"toString": func(args []IType) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 0 {
				return nil, fmt.Errorf("float.toString expects 0 arguments")
			}

			// Convert to string
			return &tString{value: fmt.Sprintf("%v", t.value)}, nil
		},
	}

	// Check if method does not exist
	if _, ok := tFloatMethods[method]; !ok {
		return func(args []IType) (IType, error) {
			return nil, fmt.Errorf("float does not have method %s", method)
		}
	}

	return tFloatMethods[method]
}
