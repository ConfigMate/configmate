package types

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

func (t tObject) Checks() map[string]Check {
	return map[string]Check{}
}

func (t tObject) ChecksDescription() map[string]string {
	return map[string]string{}
}
