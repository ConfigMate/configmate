package types

import "fmt"

type tObject struct{}

func objectFactory(value interface{}) (IType, error) {
	return &tObject{}, nil
}

func (t tObject) TypeName() string {
	return "object"
}

func (t tObject) Value() interface{} {
	return nil
}

func (t tObject) Methods() []string {
	return []string{}
}

func (t tObject) MethodDescription(method string) string {
	return ""
}

func (t tObject) GetMethod(method string) Method {
	return func(args []IType) (IType, error) {
		return nil, fmt.Errorf("object does not have a method %s", method)
	}
}
