package types

import (
	"fmt"
	"strings"

	"github.com/ConfigMate/configmate/parsers"
)

type tFactoryMethod func(value interface{}) (IType, error)

func MakeType(typename string, value interface{}) (IType, error) {
	if strings.HasPrefix(typename, "list:") {
		return listFactory(typename[5:], value.([]*parsers.Node))
	}

	factories := map[string]tFactoryMethod{
		"bool":   boolFactory,
		"int":    intFactory,
		"float":  floatFactory,
		"string": stringFactory,
		"object": objectFactory,
	}

	if factory, ok := factories[typename]; ok {
		return factory(value)
	}

	return nil, fmt.Errorf("type %s does not exist", typename)
}
