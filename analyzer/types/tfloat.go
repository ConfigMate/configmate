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

func (t *tFloat) Value() interface{} {
	return t.value
}

func (t *tFloat) Checks() map[string]Check {
	return map[string]Check{
		"eq": func(args ...interface{}) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 1 {
				return nil, fmt.Errorf("float.eq expects 1 argument")
			}

			// Check that the argument is a float
			if _, ok := args[0].(float64); !ok {
				return nil, fmt.Errorf("float.eq expects a float argument")
			}

			// Check that the argument is equal to the value
			if args[0].(float64) != t.value {
				return &tBool{value: false}, fmt.Errorf("float.eq failed: %v != %v", args[0].(float64), t.value)
			}

			return &tBool{value: true}, nil
		},
		"gt": func(args ...interface{}) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 1 {
				return nil, fmt.Errorf("float.gt expects 1 argument")
			}

			// Check that the argument is a float
			if _, ok := args[0].(float64); !ok {
				return nil, fmt.Errorf("float.gt expects a float argument")
			}

			// Check that the value is greater than the argument
			if t.value <= args[0].(float64) {
				return &tBool{value: false}, fmt.Errorf("float.gt failed: %v <= %v", t.value, args[0].(float64))
			}

			return &tBool{value: true}, nil
		},
		"gte": func(args ...interface{}) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 1 {
				return nil, fmt.Errorf("float.gte expects 1 argument")
			}

			// Check that the argument is a float
			if _, ok := args[0].(float64); !ok {
				return nil, fmt.Errorf("float.gte expects a float argument")
			}

			// Check that the value is greater than or equal to the argument
			if t.value < args[0].(float64) {
				return &tBool{value: false}, fmt.Errorf("float.gte failed: %v < %v", t.value, args[0].(float64))
			}

			return &tBool{value: true}, nil
		},
		"lt": func(args ...interface{}) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 1 {
				return nil, fmt.Errorf("float.lt expects 1 argument")
			}

			// Check that the argument is a float
			if _, ok := args[0].(float64); !ok {
				return nil, fmt.Errorf("float.lt expects a float argument")
			}

			// Check that the value is less than the argument
			if t.value >= args[0].(float64) {
				return &tBool{value: false}, fmt.Errorf("float.lt failed: %v >= %v", t.value, args[0].(float64))
			}

			return &tBool{value: true}, nil
		},
		"lte": func(args ...interface{}) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 1 {
				return nil, fmt.Errorf("float.lte expects 1 argument")
			}

			// Check that the argument is a float
			if _, ok := args[0].(float64); !ok {
				return nil, fmt.Errorf("float.lte expects a float argument")
			}

			// Check that the value is less than or equal to the argument
			if t.value > args[0].(float64) {
				return &tBool{value: false}, fmt.Errorf("float.lte failed: %v > %v", t.value, args[0].(float64))
			}

			return &tBool{value: true}, nil
		},
		"range": func(args ...interface{}) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 2 {
				return nil, fmt.Errorf("float.range expects 2 arguments")
			}

			// Check that the arguments are floats
			if _, ok := args[0].(float64); !ok {
				return nil, fmt.Errorf("float.range expects a float argument")
			}
			if _, ok := args[1].(float64); !ok {
				return nil, fmt.Errorf("float.range expects a float argument")
			}

			// Check that the value is within the range
			if t.value < args[0].(float64) || t.value > args[1].(float64) {
				return &tBool{value: false}, fmt.Errorf("float.range failed: %v not in [%v, %v]", t.value, args[0].(float64), args[1].(float64))
			}

			return &tBool{value: true}, nil
		},
		"toInt": func(args ...interface{}) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 0 {
				return nil, fmt.Errorf("float.toInt expects 0 arguments")
			}

			// Convert to int
			return &tInt{value: int(t.value)}, nil
		},
		"toString": func(args ...interface{}) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 0 {
				return nil, fmt.Errorf("float.toString expects 0 arguments")
			}

			// Convert to string
			return &tString{value: fmt.Sprintf("%v", t.value)}, nil
		},
	}
}

func (t *tFloat) ChecksDescription() map[string]string {
	return map[string]string{
		"eq":       "float.eq(arg float) : Checks that the value is equal to the argument",
		"gt":       "float.gt(arg float) : Checks that the value is greater than the argument",
		"gte":      "float.gte(arg float) : Checks that the value is greater than or equal to the argument",
		"lt":       "float.lt(arg float) : Checks that the value is less than the argument",
		"lte":      "float.lte(arg float) : Checks that the value is less than or equal to the argument",
		"range":    "float.range(min float, max float) : Checks that the value is within the range",
		"toInt":    "float.toInt() : Converts the value to an int",
		"toString": "float.toString() : Converts the value to a string",
	}
}
