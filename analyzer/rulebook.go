package analyzer

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
	Description string   `json:"description" toml:"description"` // Description of the rule
	CheckName   string   `json:"check" toml:"check"`             // Name of the check
	Args        []string `json:"args" toml:"args"`               // List of arguments to the check
}
