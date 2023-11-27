package types

import (
	"fmt"
	"strings"

	"github.com/ConfigMate/configmate/analyzer/spec"
	"github.com/ConfigMate/configmate/parsers"
)

func MakeType(typename string, value interface{}) (IType, error) {
	return tf.makeType(typename, value)
}

func AddCustomObjTypes(customDefs []spec.ObjectDef) {
	tf.addCustomObjTypes(customDefs)
}

var tf typeFactory

type typeFactory struct {
	factories      map[string]tFactoryMethod
	customObjTypes map[string]spec.ObjectDef
}

type tFactoryMethod func(value interface{}) (IType, error)

func init() {
	tf = typeFactory{
		factories: map[string]tFactoryMethod{
			"bool":      boolFactory,
			"int":       intFactory,
			"float":     floatFactory,
			"string":    stringFactory,
			"object":    objectFactory,
			"host":      hostFactory,
			"port":      portFactory,
			"host_port": hostPortFactory,
			"file":      fileFactory,
		},
		customObjTypes: make(map[string]spec.ObjectDef),
	}
}

func (tf *typeFactory) makeType(typename string, value interface{}) (IType, error) {
	if strings.HasPrefix(typename, "list<") && strings.HasSuffix(typename, ">") {
		return listFactory(typename[5:len(typename)-1], value)
	}

	if customDef, ok := tf.customObjTypes[typename]; ok {
		return customObjectFactory(value.(map[string]*parsers.Node), customDef)
	}

	if factory, ok := tf.factories[typename]; ok {
		return factory(value)
	}

	return nil, fmt.Errorf("type %s does not exist", typename)
}

func (tf *typeFactory) addCustomObjTypes(customDefs []spec.ObjectDef) {
	for _, def := range customDefs {
		tf.customObjTypes[def.Name] = def
	}
}
