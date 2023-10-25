package analyzer

type Rulebook struct {
	Name  string                 `json:"name" toml:"name"`   // Name of the rulebook
	Files map[string]FileDetails `json:"files" toml:"files"` // Map of file aliases to file paths and file formats
	Rules []Rule                 `json:"rules" toml:"rules"` // List of rules
}

type FileDetails struct {
	Path   string `json:"path" toml:"path"`
	Format string `json:"format" toml:"format"`
}

type Rule struct {
	Field    string      `json:"field" toml:"field"`       // Field to check
	Type     string      `json:"type" toml:"type"`         // Type of the field
	Optional bool        `json:"optional" toml:"optional"` // Whether the field is optional
	Checks   []string    `json:"checks" toml:"checks"`     // List of checks to perform
	Default  interface{} `json:"default" toml:"default"`   // Default value of the field
	Notes    string      `json:"notes" toml:"notes"`       // Notes about the rule
}
