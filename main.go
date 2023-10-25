package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/ConfigMate/configmate/analyzer"
	"github.com/ConfigMate/configmate/parsers"
	"github.com/ConfigMate/configmate/server"
	"github.com/ConfigMate/configmate/utils"
	"github.com/urfave/cli/v2"
)

// Version contains the semver version number for the controller. (set via ldflags during build)
var Version = "0.0.0-none"

// BuildDate contains the build date and time in RFC 3339 format. (set via ldflags during build)
var BuildDate = "Not Provided"

// GitHash contains the Git commit hash from which the controller was built. (set via ldflags during build)
var GitHash = "Not Provided"

func main() {
	app := &cli.App{
		Name:                 "configm",
		Usage:                "Tool to check configuration files for errors in content.",
		UsageText:            "configm [global options] command [command options] [arguments...]",
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "Enable verbose output.",
			},
		},
		Commands: []*cli.Command{
			{
				Name:      "version",
				Usage:     "Print version information.",
				UsageText: "configm version",
				Action: func(c *cli.Context) error {
					// Print version information
					fmt.Printf("Version: %s\n", Version)

					// Print build information if verbose flag is set
					if c.Bool("verbose") {
						fmt.Printf("Build Date: %s\n", BuildDate)
						fmt.Printf("Git Hash: %s\n", GitHash)
					}

					return nil
				},
			},
			{
				Name:      "check",
				Usage:     "Check configuration files for errors in content.",
				UsageText: "configm check <path-to-rulebook>",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "all",
						Aliases: []string{"a"},
						Usage:   "Outputs the result for successful checks also.",
					},
				},
				Action: func(c *cli.Context) error {
					// Check number of arguments
					if c.NArg() != 1 {
						return fmt.Errorf("invalid number of arguments")
					}

					// Get the rulebook path from the arguments
					rulebookPath := c.Args().Get(0)

					// Read the rulebook file
					ruleBookData, err := os.ReadFile(rulebookPath)
					if err != nil {
						return err
					}

					// Decode the TOML data into the Rulebook struct
					var rulebook analyzer.Rulebook
					if _, err := toml.Decode(string(ruleBookData), &rulebook); err != nil {
						return fmt.Errorf("error decoding file into a rulebook object: %v", err)
					}

					// Parse rulebooks
					files := make(map[string]*parsers.Node)
					for alias, file := range rulebook.Files {
						// Read the file
						data, err := os.ReadFile(file.Path)
						if err != nil {
							return err
						}

						// Parse the file
						parsedFile, err := parsers.Parse(data, file.Format)
						if err != nil {
							return err
						}

						// Append the parse result to the files map
						files[alias] = parsedFile
					}

					// Get rules
					rules := rulebook.Rules

					// Map the file content to the corresponding line numbers
					filesLines := make(map[string]map[int]string)
					for alias, details := range rulebook.Files {
						// Read file
						fileData, err := os.ReadFile(details.Path)
						if err != nil {
							return err
						}

						// Create line map
						lineMap, err := utils.CreateLineMap(fileData)
						if err != nil {
							return nil
						}

						// Append line map to filesLines
						filesLines[alias] = lineMap
					}

					// Get analyzer
					a := analyzer.NewAnalyzer()
					res, err := a.AnalyzeConfigFiles(files, rules)
					if err != nil {
						return err
					}

					successfulChecks := make([]analyzer.Result, 0)
					// Print results for failed checks
					for _, result := range res {
						if !result.Passed {
							formattedResult := utils.FormatResult(result, rulebook.Files, filesLines)
							fmt.Print(formattedResult)
						} else {
							successfulChecks = append(successfulChecks, result)
						}
					}

					if len(successfulChecks) == len(res) {
						fmt.Println("All checks passed!")
					}

					// Print results for successful checks if --all flag is set
					if c.Bool("all") {
						for _, result := range successfulChecks {
							formattedResult := utils.FormatResult(result, rulebook.Files, filesLines)
							fmt.Print(formattedResult)
						}
					}

					return nil
				},
			},
			{
				Name:      "serve",
				Usage:     "Start a web server to check configuration files for errors in content.",
				UsageText: "configm serve --port <port>",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "port",
						Aliases: []string{"p"},
						Usage:   "Port to listen on.",
						Value:   10007,
					},
				},
				Action: func(c *cli.Context) error {
					// Create server
					srv := server.CreateServer(c.Int("port"), analyzer.NewAnalyzer())

					// Start server
					if err := srv.Serve(); err != nil {
						return err
					}

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		// Print error message
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	// Exit with success
	os.Exit(0)
}
