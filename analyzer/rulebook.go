package analyzer

type Rulebook struct {
	Name        string            `json:"name" toml:"name"`
	Description string            `json:"description" toml:"description"`
	Files       map[string]string `json:"files" toml:"files"`
	Rules       []Rule            `json:"rules" toml:"rules"`
}

// Rule declares a check and it's arguments in a rulebook.
// If argument is a fields/value in a config file, it should be in the format "f:file_alias.key",
// where file_alias is the alias of the file in the rulebook and key is the path separeted by dots of the field/value.
// If argument is a literal, it should be in the format "l:value".
type Rule struct {
	Description string   `json:"description" toml:"description"`
	CheckName   string   `json:"check" toml:"check"`
	Args        []string `json:"args" toml:"args"`
}
