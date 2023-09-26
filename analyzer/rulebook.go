package analyzer

type Rulebook struct {
	Name        string            `json:"name" toml:"name"`
	Description string            `json:"description" toml:"description"`
	Files       map[string]string `json:"files" toml:"files"`
	Rules       []Rule            `json:"rules" toml:"rules"`
}

type Rule struct {
	Description string   `json:"description" toml:"description"`
	CheckName   string   `json:"check" toml:"check"`
	Args        []string `json:"args" toml:"args"`
}
