package types

// The return values of a check must be as follows:
// Check can be correctly executed, and the condition checked is true:
//   - return &tBool{value: true}, nil
//
// Check can be correctly executed, but the condition checked is false:
//   - return &tBool{value: false}, error("message indicating why the check failed")
//
// Check cannot be correctly executed (the arguments are invalid):
//   - return nil, error("error message")
type Check func(args []IType) (IType, error)

type IType interface {
	TypeName() string
	Value() interface{}
	Checks() map[string]Check
	ChecksDescription() map[string]string
}
