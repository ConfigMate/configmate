package analyzer

type Analyzer interface {
	Analyze(rb Rulebook) (res Result, err error)
}

type Result struct{}
