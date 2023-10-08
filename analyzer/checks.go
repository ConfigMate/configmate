package analyzer

import (
	"fmt"
	"strconv"
	"strings"
)

// Check is a binary function that varifies a specific condition.
type Check interface {
	Check(values []interface{}) (passed bool, comment string, err error)
	GetArgsSourceAndTypes() ([]CheckArgSource, []CheckArgType)
}

// CheckArg is an argument to a rule. It is encoded as a string, and it
// has the format "source:type:value", where source is either "f:" or "l:",
// type is either "i:", "f:", "b:", "s:" and value is the value
// of the argument.
// If the argument is a literal, the value is parsed from the provided string value.
// If the argument is a field/value in a config file, the value has the format "file_alias.path.to.field".
type CheckArg struct {
	v interface{}    // value
	s CheckArgSource // source
	t CheckArgType   // type
}

// CheckArgSource is the source of a rule argument.
// It can either come from a file or be a literal.
type CheckArgSource int

const (
	File    CheckArgSource = iota // represented with "f:"
	Literal                       // represented with "l:"
)

func (s CheckArgSource) String() string {
	switch s {
	case File:
		return "file"
	case Literal:
		return "literal"
	default:
		return ""
	}
}

// CheckArgType is the type of a rule argument.
// It can either be an int, float, bool or string.
type CheckArgType int

const (
	Int    CheckArgType = iota // represented with "i:"
	Float                      // represented with "f:"
	Bool                       // represented with "b:"
	String                     // represented with "s:"
)

func (t CheckArgType) String() string {
	switch t {
	case Int:
		return "int"
	case Float:
		return "float"
	case Bool:
		return "bool"
	case String:
		return "string"
	default:
		return ""
	}
}

func ParseCheckArg(arg string) (*CheckArg, error) {
	// Split argument
	segments := strings.Split(string(arg), ":")

	// Check number of segments
	if len(segments) != 3 {
		return nil, fmt.Errorf("invalid number of segments in rule argument: %s", arg)
	}

	// Check s
	var s CheckArgSource
	switch segments[0] {
	case "f":
		s = File
	case "l":
		s = Literal
	default:
		return nil, fmt.Errorf("invalid source in rule argument: %s", arg)
	}

	// Check type
	var t CheckArgType
	switch segments[1] {
	case "i":
		t = Int
	case "f":
		t = Float
	case "b":
		t = Bool
	case "s":
		t = String
	default:
		return nil, fmt.Errorf("invalid type in rule argument: %s", arg)
	}

	// Check value is not empty except for literal string
	if (segments[0] == "f" || segments[1] != "s") && segments[2] == "" {
		return nil, fmt.Errorf("value cannot be empty in rule argument: %s", arg)
	}

	// Decode value
	var v interface{}
	switch s {
	case File:
		// Decode file value
		if fileValue, err := decodeFileValue(segments[2]); err != nil {
			return nil, fmt.Errorf("failed to decode file value: %s", err.Error())
		} else {
			v = *fileValue
		}
	case Literal:
		// Decode literal value
		if literalValue, err := decodeLiteral(segments[2], t); err != nil {
			return nil, fmt.Errorf("failed to decode literal value: %s", err.Error())
		} else {
			v = literalValue
		}
	}

	return &CheckArg{
		s: s,
		t: t,
		v: v,
	}, nil
}

type FileValue struct {
	alias string
	path  string
}

// decodeFileValue returns the alias and path of the given file value.
// File values look like these: "file_alias.server.port", "file_alias.settings.users[0].name".
func decodeFileValue(value string) (*FileValue, error) {
	// Split the value based on the dot
	segments := strings.SplitN(value, ".", 2)

	// Check number of segments
	if len(segments) != 2 {
		return nil, fmt.Errorf("invalid number of segments in file value: %s", value)
	}

	// Check neither segment is empty
	if segments[0] == "" || segments[1] == "" {
		return nil, fmt.Errorf("invalid file value: %s", value)
	}

	// Return file value
	return &FileValue{
		alias: segments[0],
		path:  segments[1],
	}, nil
}

// decodeLiteral attempts to convert the given literal argument to the given type.
func decodeLiteral(value string, t CheckArgType) (interface{}, error) {
	switch t {
	case Int:
		// Verify value is an integer
		if value, err := strconv.Atoi(value); err != nil {
			return nil, fmt.Errorf("failed to decode %v as int: %s", value, err.Error())
		} else {
			return value, nil
		}
	case Float:
		// Verify value is a float
		if value, err := strconv.ParseFloat(value, 64); err != nil {
			return nil, fmt.Errorf("failed to decode %v as float: %s", value, err.Error())
		} else {
			return value, nil
		}
	case Bool:
		// Verify value is a bool
		if value, err := strconv.ParseBool(value); err != nil {
			return nil, fmt.Errorf("failed to interpret %v as bool: %s", value, err.Error())
		} else {
			return value, nil
		}
	case String:
		// Add value
		return value, nil
	default:
		return nil, fmt.Errorf("unknown argument type: %s", t.String())
	}
}
