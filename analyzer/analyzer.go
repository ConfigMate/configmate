package analyzer

type Analyzer interface {
	Analyze(rb Rulebook) (res []Result, err error)
}

type Result struct {
	Passed        bool            `json:"passed"`         // true if the check passed, false if it failed
	ResultComment string          `json:"result_comment"` // an error msg or comment about the result
	TokenList     []TokenLocation `json:"token_list"`     // a list of tokens that were involved in the rule
}

type TokenLocation struct {
	File   string `json:"file"`
	Line   int    `json:"line"`
	Column int    `json:"column"`
	Length int    `json:"length"`
}

type AnalyzerImpl struct {
}

func (a *AnalyzerImpl) Analyze(rb Rulebook) (res []Result, err error) {
	return []Result{
		{
			Passed:        true,
			ResultComment: "This is a test",
			TokenList: []TokenLocation{
				{
					File:   "test.toml",
					Line:   1,
					Column: 1,
					Length: 1,
				},
			},
		},
	}, nil
}
