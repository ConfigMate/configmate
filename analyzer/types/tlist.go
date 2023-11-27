package types

import (
	"fmt"

	"github.com/ConfigMate/configmate/parsers"
)

type tList struct {
	listType string
	values   []IType
}

func listFactory(typename string, value interface{}) (IType, error) {
	// Check that the value is a list
	listValues, ok := value.([]*parsers.Node)
	if !ok {
		return nil, fmt.Errorf("value is not a list")
	}

	// Create a new list
	list := &tList{
		listType: typename,
		values:   make([]IType, len(listValues)),
	}

	for i, value := range listValues {
		var err error
		if list.values[i], err = MakeType(typename, value.Value); err != nil {
			return nil, err
		}
	}

	return list, nil
}

func (t tList) TypeName() string {
	return "list<" + t.listType + ">"
}

func (t tList) Value() interface{} {
	return t.values
}

func (t tList) Methods() []string {
	return []string{
		"at",
		"len",
	}
}

func (t tList) MethodDescription(method string) string {
	tListMethodsDescriptions := map[string]string{
		"at":  "list.at(index int) elementtype - returns the element at the given index",
		"len": "list.len() int - returns the length of the list",
	}

	return tListMethodsDescriptions[method]
}

func (t tList) GetMethod(method string) Method {
	tListMethods := map[string]Method{
		"at": func(args []IType) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 1 {
				return nil, fmt.Errorf("list.at expects 1 argument")
			}

			// Cast argument to int type
			i, ok := args[0].(*tInt)
			if !ok {
				return nil, fmt.Errorf("list.at expects an int argument")
			}

			// Check that the index is in range
			if i.value < 0 || i.value >= len(t.values) {
				return nil, fmt.Errorf("list.at failed: index %v out of range", i.value)
			}
			return t.values[i.value], nil
		},
		"len": func(args []IType) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 0 {
				return nil, fmt.Errorf("list.len expects 0 arguments")
			}

			// Return the length of the list
			return &tInt{value: len(t.values)}, nil
		},
	}

	// Check if method doesn't exist
	if _, ok := tListMethods[method]; !ok {
		return func(args []IType) (IType, error) {
			return nil, fmt.Errorf("list does not have method %s", method)
		}
	}

	return tListMethods[method]
}
