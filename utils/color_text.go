package utils

type StdOutColor string

const (
	Red    StdOutColor = "\033[31m"
	Green  StdOutColor = "\033[32m"
	Yellow StdOutColor = "\033[33m"
	Blue   StdOutColor = "\033[34m"
	Purple StdOutColor = "\033[35m"
	Cyan   StdOutColor = "\033[36m"
	Gray   StdOutColor = "\033[37m"
	White  StdOutColor = "\033[97m"
)

const Reset string = "\033[0m"

func ColorText(text string, color StdOutColor) string {
	value := string(color) + text + string(Reset)
	return value
}
