package types

import "fmt"

type tFactoryMethod func(value interface{}) (IType, error)

func MakeType(typename string, value interface{}) (IType, error) {
	factories := map[string]tFactoryMethod{
		"bool":   boolFactory,
		"int":    intFactory,
		"float":  floatFactory,
		"string": stringFactory,
	}

	if factory, ok := factories[typename]; ok {
		return factory(value)
	}

	return nil, fmt.Errorf("type %s does not exist", typename)
}
