package types

import (
	"fmt"

	"github.com/ConfigMate/configmate/parsers"
)

type tList struct {
	listType string
	values   []IType
}

func listFactory(typename string, values []*parsers.Node) (IType, error) {
	var err error
	list := &tList{
		listType: typename,
		values:   make([]IType, len(values)),
	}
	for i, value := range values {
		if list.values[i], err = MakeType(typename, value.Value); err != nil {
			return nil, err
		}
	}

	return list, nil
}

func (t tList) TypeName() string {
	if len(t.values) == 0 {
		return "list:" + t.listType
	}
	return "list:" + t.values[0].TypeName()
}

func (t tList) Value() interface{} {
	return t.values
}

func (t tList) Checks() map[string]Check {
	return map[string]Check{
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
}

func (t tList) ChecksDescription() map[string]string {
	return map[string]string{
		"at":  "list.at(index int) - returns the element at the given index",
		"len": "list.len() - returns the length of the list",
	}
}
