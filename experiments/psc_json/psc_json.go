package main

import (
	"fmt"
	"io"
	"os"

	"github.com/ConfigMate/configmate/parsers"
)

func main() {
	// Load file: ../examples/sample_config.json
	file, err := os.Open("../../examples/sample_config.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	parser := &parsers.JsonParser{}

	b, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	res, err := parser.Parse(b)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("%+v\n", res)

	os.Exit(0)
}
