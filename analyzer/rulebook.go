package analyzer

type Rulebook struct {
	Name        string            `json:"name" toml:"name"`               // Name of the rulebook
	Description string            `json:"description" toml:"description"` // Description of the rulebook
	Files       map[string]string `json:"files" toml:"files"`             // Map of file aliases to file paths
	Rules       []Rule            `json:"rules" toml:"rules"`             // List of rules
}

type Rule struct {
	Field       string   `json:"field" toml:"field"`             // Field to check
	Description string   `json:"description" toml:"description"` // Description of the rule
	Type        string   `json:"type" toml:"type"`               // Type of the field
	Optional    bool     `json:"optional" toml:"optional"`       // Whether the field is optional
	Default     string   `json:"default" toml:"default"`         // Default value of the field
	Checks      []string `json:"checks" toml:"checks"`           // List of checks to perform
	Notes       string   `json:"notes" toml:"notes"`             // Notes about the rule
}
