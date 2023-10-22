package types

type Check func(args ...interface{}) (IType, error)

type IType interface {
	Value() interface{}
	Checks() map[string]Check
	ChecksDescription() map[string]string
}
