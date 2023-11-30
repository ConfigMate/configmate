package types

import (
	"fmt"

	"github.com/ConfigMate/configmate/analyzer/spec"
	"github.com/ConfigMate/configmate/parsers"
)

var tCustomObjectMethodsDescriptions map[string]string = map[string]string{
	"get": "custom_object.get(field string) fieldtype : Gets the specified field",
}

// Special error to indicate that an optional field is missing
var OptMissFieldError error

type tCustomObject struct {
	ObjectName       string
	Fields           map[string]IType
	OptMissingFields map[string]bool
}

func customObjectFactory(objValue map[string]*parsers.Node, definition spec.ObjectDef) (IType, error) {
	// Create a new customObj
	customObj := &tCustomObject{
		ObjectName:       definition.Name,
		Fields:           make(map[string]IType),
		OptMissingFields: make(map[string]bool),
	}

	for _, prop := range definition.Properties {
		if propNode, ok := objValue[prop.Name]; ok {
			t, err := MakeType(prop.Type, propNode.Value)
			if err != nil {
				return nil, err
			}
			customObj.Fields[prop.Name] = t
		} else if prop.Optional {
			customObj.OptMissingFields[prop.Name] = true
		} else {
			return nil, fmt.Errorf("missing required property %s", prop.Name)
		}
	}

	return customObj, nil
}

func (t tCustomObject) TypeName() string {
	return t.ObjectName
}

func (t tCustomObject) Value() interface{} {
	return t.Fields
}

func (t tCustomObject) GetMethod(method string) Method {
	tCustomObjectMethods := map[string]Method{
		"get": func(args []IType) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 1 {
				return nil, fmt.Errorf("%s.get expects 1 argument", t.ObjectName)
			}

			// Cast argument to string type
			field, ok := args[0].(*tString)
			if !ok {
				return nil, fmt.Errorf("argument to %s.get must be a string", t.ObjectName)
			}

			// Check if the field is optional and missing
			if _, ok := t.OptMissingFields[field.value]; ok {
				return nil, OptMissFieldError
			}

			// Check that the field exists
			if _, ok := t.Fields[field.value]; !ok {
				return nil, fmt.Errorf("%s does not have field %s", t.ObjectName, field.Value().(string))
			}

			return t.Fields[field.value], nil
		},
	}

	// Check if method doesn't exist
	if _, ok := tCustomObjectMethods[method]; !ok {
		return func(args []IType) (IType, error) {
			return nil, fmt.Errorf("%s does not have method %s", t.ObjectName, method)
		}
	}

	return tCustomObjectMethods[method]
}
