package types

import (
	"fmt"
	"strconv"
)

var tIntMethodsDescriptions map[string]string = map[string]string{
	"eq":       "int.eq(arg int) bool : Checks that the value is equal to the argument",
	"gt":       "int.gt(arg int) bool : Checks that the value is greater than the argument",
	"gte":      "int.gte(arg int) bool : Checks that the value is greater than or equal to the argument",
	"lt":       "int.lt(arg int) bool : Checks that the value is less than the argument",
	"lte":      "int.lte(arg int) bool : Checks that the value is less than or equal to the argument",
	"range":    "int.range(min int, max int) bool : Checks that the value is in the range [min, max]",
	"toFloat":  "int.toFloat() float : Converts the value to a float",
	"toString": "int.toString() string : Converts the value to a string",
}

type tInt struct {
	value int
}

func intFactory(value interface{}) (IType, error) {
	if value, ok := value.(int); ok {
		return &tInt{value: value}, nil
	}

	return nil, fmt.Errorf("value is not an int")
}

func (t tInt) TypeName() string {
	return "int"
}

func (t tInt) Value() interface{} {
	return t.value
}

func (t tInt) GetMethod(method string) Method {
	tIntMethods := map[string]Method{
		"eq": func(args []IType) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 1 {
				return nil, fmt.Errorf("int.eq expects 1 argument")
			}

			// Cast argument to int type
			i, ok := args[0].(*tInt)
			if !ok {
				return nil, fmt.Errorf("int.eq expects an int argument")
			}

			// Check that the argument is equal to the value
			if i.value != t.value {
				return &tBool{value: false}, fmt.Errorf("int.eq failed: %v != %v", i.value, t.value)
			}
			return &tBool{value: true}, nil
		},
		"gt": func(args []IType) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 1 {
				return nil, fmt.Errorf("int.gt expects 1 argument")
			}

			// Cast argument to int type
			i, ok := args[0].(*tInt)
			if !ok {
				return nil, fmt.Errorf("int.gt expects an int argument")
			}

			// Check that the value is greater than the argument
			if t.value <= i.value {
				return &tBool{value: false}, fmt.Errorf("int.gt failed: %v <= %v", t.value, i.value)
			}
			return &tBool{value: true}, nil
		},
		"gte": func(args []IType) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 1 {
				return nil, fmt.Errorf("int.gte expects 1 argument")
			}

			// Cast argument to int type
			i, ok := args[0].(*tInt)
			if !ok {
				return nil, fmt.Errorf("int.gte expects an int argument")
			}

			// Check that the value is greater than or equal to the argument
			if t.value < i.value {
				return &tBool{value: false}, fmt.Errorf("int.gte failed: %v < %v", t.value, i.value)
			}
			return &tBool{value: true}, nil
		},
		"lt": func(args []IType) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 1 {
				return nil, fmt.Errorf("int.lt expects 1 argument")
			}

			// Cast argument to int type
			i, ok := args[0].(*tInt)
			if !ok {
				return nil, fmt.Errorf("int.lt expects an int argument")
			}

			// Check that the value is less than the argument
			if t.value >= i.value {
				return &tBool{value: false}, fmt.Errorf("int.lt failed: %v >= %v", t.value, i.value)
			}
			return &tBool{value: true}, nil
		},
		"lte": func(args []IType) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 1 {
				return nil, fmt.Errorf("int.lte expects 1 argument")
			}

			// Cast argument to int type
			i, ok := args[0].(*tInt)
			if !ok {
				return nil, fmt.Errorf("int.lte expects an int argument")
			}

			// Check that the value is less than or equal to the argument
			if t.value > i.value {
				return &tBool{value: false}, fmt.Errorf("int.lte failed: %v > %v", t.value, i.value)
			}
			return &tBool{value: true}, nil
		},
		"range": func(args []IType) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 2 {
				return nil, fmt.Errorf("int.range expects 2 arguments")
			}

			// Cast arguments to int type
			i0, ok := args[0].(*tInt)
			if !ok {
				return nil, fmt.Errorf("int.range expects an int for first argument")
			}

			i1, ok := args[1].(*tInt)
			if !ok {
				return nil, fmt.Errorf("int.range expects an int for second argument")
			}

			// Check that the value is in the range
			if t.value < i0.value || t.value > i1.value {
				return &tBool{value: false}, fmt.Errorf("int.range failed: %v not in range [%v, %v]", t.value, i0.value, i1.value)
			}
			return &tBool{value: true}, nil
		},
		"toFloat": func(args []IType) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 0 {
				return nil, fmt.Errorf("int.toFloat expects 0 arguments")
			}

			// Convert to float
			return &tFloat{value: float64(t.value)}, nil
		},
		"toString": func(args []IType) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 0 {
				return nil, fmt.Errorf("int.toString expects 0 arguments")
			}

			// Convert to string
			return &tString{value: strconv.Itoa(t.value)}, nil
		},
	}

	// Check if method doesn't exist
	if _, ok := tIntMethods[method]; !ok {
		return func(args []IType) (IType, error) {
			return nil, fmt.Errorf("int does not have method %s", method)
		}
	}

	return tIntMethods[method]
}
