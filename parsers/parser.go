package parsers

type Parser interface {
	Parse(filename string) (ConfigFile, error)
}
