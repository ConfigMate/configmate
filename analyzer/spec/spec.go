package spec

import "github.com/ConfigMate/configmate/parsers"

type Specification struct {
	File       string            `json:"file"`        // File this specification is for
	FileFormat string            `json:"file_format"` // Format of the file
	Imports    map[string]string `json:"imports"`     // Imported rulebooks with their aliases
	Fields     []FieldSpec       `json:"fields"`      // Node that holds the specification of the file

	FileLocation         parsers.TokenLocation            `json:"file_location"`          // Location of the file specification
	FileFormatLocation   parsers.TokenLocation            `json:"file_format_location"`   // Location of the file format
	ImportsAliasLocation map[string]parsers.TokenLocation `json:"imports_alias_location"` // Location of the imports alias
	ImportsLocation      map[string]parsers.TokenLocation `json:"imports_location"`       // Location of the imports field
}

type FieldSpec struct {
	Field    *parsers.NodeKey    `json:"field"`    // Field to check
	Type     string              `json:"type"`     // Type of the field
	Optional bool                `json:"optional"` // Whether the field is optional
	Default  string              `json:"default"`  // Default value of the field
	Notes    string              `json:"notes"`    // Notes about the rule
	Checks   []CheckWithLocation `json:"checks"`   // List of checks to perform

	FieldLocation    parsers.TokenLocation `json:"field_location"`    // Location of the field
	TypeLocation     parsers.TokenLocation `json:"type_location"`     // Location of the type
	OptionalLocation parsers.TokenLocation `json:"optional_location"` // Location of the optional field
	DefaultLocation  parsers.TokenLocation `json:"default_location"`  // Location of the default field
	NotesLocation    parsers.TokenLocation `json:"notes_location"`    // Location of the notes field
}

type CheckWithLocation struct {
	Check    string                `json:"check"`    // Name of the check
	Location parsers.TokenLocation `json:"location"` // Location of the check
}
