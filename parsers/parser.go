package parsers

type Parser interface {
	Parse(filename string) (*Node, error)
}
