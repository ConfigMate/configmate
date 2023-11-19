package spec

import "github.com/ConfigMate/configmate/parsers"

type Specification struct {
	File       string            // File this specification is for
	FileFormat string            // Format of the file
	Imports    map[string]string // Imported rulebooks with their aliases
	Fields     []FieldSpec       // Node that holds the specification of the file

	FileLocation         parsers.TokenLocation            // Location of the file specification
	FileFormatLocation   parsers.TokenLocation            // Location of the file format
	ImportsAliasLocation map[string]parsers.TokenLocation // Location of the imports alias
	ImportsLocation      map[string]parsers.TokenLocation // Location of the imports field
}

type FieldSpec struct {
	Field    string              // Field to check
	Type     string              // Type of the field
	Optional bool                // Whether the field is optional
	Default  string              // Default value of the field
	Notes    string              // Notes about the rule
	Checks   []CheckWithLocation // List of checks to perform

	FieldLocation    parsers.TokenLocation // Location of the field
	TypeLocation     parsers.TokenLocation // Location of the type
	OptionalLocation parsers.TokenLocation // Location of the optional field
	DefaultLocation  parsers.TokenLocation // Location of the default field
	NotesLocation    parsers.TokenLocation // Location of the notes field
}

type CheckWithLocation struct {
	Check    string                // Name of the check
	Location parsers.TokenLocation // Location of the check
}
