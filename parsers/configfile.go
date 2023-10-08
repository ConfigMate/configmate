package parsers

// FieldType is the type of a field in a configuration file.
type FieldType int

const (
	Null FieldType = iota
	Bool
	Int
	Float
	String
	Array
	Object
)

func (ft FieldType) String() string {
	switch ft {
	case Unknown:
		return "unknown"
	case Bool:
		return "bool"
	case Int:
		return "int"
	case Float:
		return "float"
	case String:
		return "string"
	case Array:
		return "array"
	case Object:
		return "object"
	default:
		return "unknown"
	}
}

// Node is a node in a configuration file. Fields of type Object will be encoded as
// a map[string]*Node and fields of type Array will be encoded as a []*Node.
type Node struct {
	Type      FieldType   // Type of field
	ArrayType FieldType   // Type of elements in array (if Type == Array)
	Value     interface{} // Value of field

	NameLocation struct { // Location of field name in configuration file
		Line   int
		Column int
		Length int
	}
	ValueLocation struct { // Location of field value in configuration file
		Line   int
		Column int
		Length int
	}
}
