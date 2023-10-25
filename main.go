package main

import (
	"fmt"
	"os"

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
				UsageText: "configm check --rulebook <rulebook>",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "rulebook",
						Aliases:  []string{"r"},
						Usage:    "Rulebook to use for checking.",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					// Get the rulebook path from the command line argument
					rulebookPath := c.String("rulebook")

					// Read the rulebook file
					ruleBookData, err := os.ReadFile(rulebookPath)
					if err != nil {
						return err
					}

					// Decode TOML into a Rulebook object
					ruleBook, err := utils.DecodeRulebook(ruleBookData)
					if err != nil {
						return err
					}

					// Parse rulebooks
					files := make(map[string]*parsers.Node)
					for _, file := range ruleBook.Files {
						// Read the file
						data, err := os.ReadFile(file.Path)
						if err != nil {
							return err
						}

						// Parse the file
						parser, err := parsers.Parse(data, file.Format)
						if err != nil {
							return err
						}

						// Append the parse result to the files map
						files[file.Path] = parser
					}

					// Get rules
					rules := ruleBook.Rules

					// Get analyzer
					analyzer := &analyzer.AnalyzerImpl{}
					res, err := analyzer.AnalyzeConfigFiles(files, rules)
					if err != nil {
						return err
					}

					// Print results
					for _, result := range res {
						fmt.Println("Result:", result)
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
					// Get analyzer
					analyzer := &analyzer.AnalyzerImpl{}

					// Create server
					srv := server.CreateServer(c.Int("port"), analyzer)

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
