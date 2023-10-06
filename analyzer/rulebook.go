package analyzer

import (
	"fmt"
	"strings"
)

type Rulebook struct {
	Name        string            `json:"name" toml:"name"`               // Name of the rulebook
	Description string            `json:"description" toml:"description"` // Description of the rulebook
	Files       map[string]string `json:"files" toml:"files"`             // Map of file aliases to file paths
	Rules       []Rule            `json:"rules" toml:"rules"`             // List of rules
}

// Rule declares a check and it's arguments in a rulebook.
// If argument is a fields/value in a config file, it should be in the format "f:<type>:file_alias.key",
// where file_alias is the alias of the file in the rulebook and key is the path separeted by dots of the field/value.
// If argument is a literal, it should be in the format "l:<type>:value".
type Rule struct {
	Description string     `json:"description" toml:"description"` // Description of the rule
	CheckName   string     `json:"check" toml:"check"`             // Name of the check
	Args        []CheckArg `json:"args" toml:"args"`               // List of arguments to the check
}

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
type CheckArg string

// CheckArgSource is the source of a rule argument.
// It can either come from a file or be a literal.
type CheckArgSource int

const (
	File    CheckArgSource = iota // represented with "f:"
	Literal                       // represented with "l:"
)

// CheckArgType is the type of a rule argument.
// It can either be an int, float, bool, string, array or object.
type CheckArgType int

const (
	Int    CheckArgType = iota // represented with "i:"
	Float                      // represented with "f:"
	Bool                       // represented with "b:"
	String                     // represented with "s:"
)

func (r CheckArg) Valid() error {
	// Split argument
	segments := strings.Split(string(r), ":")

	// Check number of segments
	if len(segments) != 3 {
		return fmt.Errorf("invalid number of segments in rule argument: %s", r)
	}

	// Check source
	if segments[0] != "f" && segments[0] != "l" {
		return fmt.Errorf("invalid source in rule argument: %s", r)
	}

	// Check type
	if segments[1] != "i" && segments[1] != "f" && segments[1] != "b" && segments[1] != "s" {
		return fmt.Errorf("invalid type in rule argument: %s", r)
	}

	// Check value is not empty except for literal string
	if (segments[0] == "f" || segments[1] != "s") && segments[2] == "" {
		return fmt.Errorf("value cannot be empty in rule argument: %s", r)
	}

	return nil
}

func (r CheckArg) Source() CheckArgSource {
	// Validate
	if err := r.Valid(); err != nil {
		return -1
	}

	// Split argument
	segments := strings.Split(string(r), ":")

	if segments[0] == "f" {
		return File
	}

	return Literal
}

func (r CheckArg) Type() CheckArgType {
	// Validate
	if err := r.Valid(); err != nil {
		return -1
	}

	// Split argument
	segments := strings.Split(string(r), ":")

	switch segments[1] {
	case "i":
		return Int
	case "f":
		return Float
	case "b":
		return Bool
	case "s":
		return String
	default:
		return -1
	}
}

func (r CheckArg) Value() string {
	// Validate
	if err := r.Valid(); err != nil {
		return ""
	}

	// Split argument
	segments := strings.Split(string(r), ":")

	return segments[2]
}
