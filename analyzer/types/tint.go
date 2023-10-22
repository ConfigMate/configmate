package types

import (
	"fmt"
	"strconv"
)

type tInt struct {
	value int
}

func intFactory(value interface{}) (IType, error) {
	if value, ok := value.(int); ok {
		return &tInt{value: value}, nil
	}

	return nil, fmt.Errorf("value %v is not an int", value)
}

func (t *tInt) Value() interface{} {
	return t.value
}

func (t *tInt) Checks() map[string]Check {
	return map[string]Check{
		"eq": func(args ...interface{}) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 1 {
				return nil, fmt.Errorf("int.eq expects 1 argument")
			}

			// Check that the argument is an int
			if _, ok := args[0].(int); !ok {
				return nil, fmt.Errorf("int.eq expects an int argument")
			}

			// Check that the argument is equal to the value
			if args[0].(int) != t.value {
				return &tBool{value: false}, fmt.Errorf("int.eq failed: %v != %v", args[0].(int), t.value)
			}

			return &tBool{value: true}, nil
		},
		"gt": func(args ...interface{}) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 1 {
				return nil, fmt.Errorf("int.gt expects 1 argument")
			}

			// Check that the argument is an int
			if _, ok := args[0].(int); !ok {
				return nil, fmt.Errorf("int.gt expects an int argument")
			}

			// Check that the value is greater than the argument
			if t.value <= args[0].(int) {
				return &tBool{value: false}, fmt.Errorf("int.gt failed: %v <= %v", t.value, args[0].(int))
			}

			return &tBool{value: true}, nil
		},
		"gte": func(args ...interface{}) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 1 {
				return nil, fmt.Errorf("int.gte expects 1 argument")
			}

			// Check that the argument is an int
			if _, ok := args[0].(int); !ok {
				return nil, fmt.Errorf("int.gte expects an int argument")
			}

			// Check that the value is greater than or equal to the argument
			if t.value < args[0].(int) {
				return &tBool{value: false}, fmt.Errorf("int.gte failed: %v < %v", t.value, args[0].(int))
			}

			return &tBool{value: true}, nil
		},
		"lt": func(args ...interface{}) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 1 {
				return nil, fmt.Errorf("int.lt expects 1 argument")
			}

			// Check that the argument is an int
			if _, ok := args[0].(int); !ok {
				return nil, fmt.Errorf("int.lt expects an int argument")
			}

			// Check that the value is less than the argument
			if t.value >= args[0].(int) {
				return &tBool{value: false}, fmt.Errorf("int.lt failed: %v >= %v", t.value, args[0].(int))
			}

			return &tBool{value: true}, nil
		},
		"lte": func(args ...interface{}) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 1 {
				return nil, fmt.Errorf("int.lte expects 1 argument")
			}

			// Check that the argument is an int
			if _, ok := args[0].(int); !ok {
				return nil, fmt.Errorf("int.lte expects an int argument")
			}

			// Check that the value is less than or equal to the argument
			if t.value > args[0].(int) {
				return &tBool{value: false}, fmt.Errorf("int.lte failed: %v > %v", t.value, args[0].(int))
			}

			return &tBool{value: true}, nil
		},
		"range": func(args ...interface{}) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 2 {
				return nil, fmt.Errorf("int.range expects 2 arguments")
			}

			// Check that the arguments are ints
			if _, ok := args[0].(int); !ok {
				return nil, fmt.Errorf("int.range expects an int argument")
			}

			if _, ok := args[1].(int); !ok {
				return nil, fmt.Errorf("int.range expects an int argument")
			}

			// Check that the value is in the range
			if t.value < args[0].(int) || t.value > args[1].(int) {
				return &tBool{value: false}, fmt.Errorf("int.range failed: %v not in range [%v, %v]", t.value, args[0].(int), args[1].(int))
			}

			return &tBool{value: true}, nil
		},
		"toFloat": func(args ...interface{}) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 0 {
				return nil, fmt.Errorf("int.toFloat expects 0 arguments")
			}

			// Convert to float
			return &tFloat{value: float64(t.value)}, nil
		},
		"toString": func(args ...interface{}) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 0 {
				return nil, fmt.Errorf("int.toString expects 0 arguments")
			}

			// Convert to string
			return &tString{value: strconv.Itoa(t.value)}, nil
		},
	}
}

func (t *tInt) ChecksDescription() map[string]string {
	return map[string]string{
		"eq":       "int.eq(arg int) : Checks that the value is equal to the argument",
		"gt":       "int.gt(arg int) : Checks that the value is greater than the argument",
		"gte":      "int.gte(arg int) : Checks that the value is greater than or equal to the argument",
		"lt":       "int.lt(arg int) : Checks that the value is less than the argument",
		"lte":      "int.lte(arg int) : Checks that the value is less than or equal to the argument",
		"range":    "int.range(min int, max int) : Checks that the value is in the range [min, max]",
		"toFloat":  "int.toFloat() : Converts the value to a float",
		"toString": "int.toString() : Converts the value to a string",
	}
}
