package types

// The return values of a check must be as follows:
// Method can be correctly executed, and the condition checked is true:
//   - return &tBool{value: true}, nil
//
// Method can be correctly executed, but the condition checked is false:
//   - return &tBool{value: false}, error("message indicating why the check failed")
//
// Method cannot be correctly executed (the arguments are invalid):
//   - return nil, error("error message")
type Method func(args []IType) (IType, error)

type IType interface {
	TypeName() string
	Value() interface{}
	GetMethod(string) Method
}
