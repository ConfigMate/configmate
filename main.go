package main

import (
	"fmt"
	"os"

	"github.com/ConfigMate/configmate/analyzer"
	"github.com/ConfigMate/configmate/analyzer/check"
	"github.com/ConfigMate/configmate/analyzer/spec"
	"github.com/ConfigMate/configmate/analyzer/types"
	"github.com/ConfigMate/configmate/files"
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
				Name:      "types",
				Usage:     "Print supported types.",
				UsageText: "configm types",
				Action: func(c *cli.Context) error {
					fmt.Println("Supported Types:")
					for _, t := range types.GetTypes() {
						fmt.Printf("\t%s\n", t)
					}

					return nil
				},
			},
			{
				Name:      "methods",
				Usage:     "Print supported methods.",
				UsageText: "configm methods <type>",
				Action: func(c *cli.Context) error {
					// Check number of arguments
					if c.NArg() != 1 {
						return fmt.Errorf("invalid number of arguments")
					}

					// Get the type from the arguments
					t := c.Args().Get(0)

					// Get methods
					methods := types.GetTypeInfo(t)
					if methods == nil {
						return fmt.Errorf("invalid type")
					}

					fmt.Printf("Supported Methods for %s:\n", t)
					for m, desc := range methods {
						fmt.Printf("\t%s: %s\n", m, desc)
					}

					return nil
				},
			},
			{
				Name:      "check",
				Usage:     "Check a configuration file specification.",
				UsageText: "configm check <path-to-specification>",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "skipped",
						Aliases: []string{"s"},
						Usage:   "Outputs the result for skipped checks.",
					},
					&cli.BoolFlag{
						Name:    "all",
						Aliases: []string{"a"},
						Usage:   "Outputs the result for successful and skipped checks also.",
					},
				},
				Action: func(c *cli.Context) error {
					// Check number of arguments
					if c.NArg() != 1 {
						return fmt.Errorf("invalid number of arguments")
					}

					// Get the rulebook path from the arguments
					specFilePath := c.Args().Get(0)

					// Get analyzer
					a := analyzer.NewAnalyzer(
						spec.NewSpecParser(),
						check.NewCheckEvaluator(),
						files.NewFileFetcher(),
						parsers.NewParserProvider(),
					)

					// Get all files
					files := a.AllFilesContent(specFilePath)

					// Map the files contents to the corresponding line numbers
					filesLines := utils.CreateLinesMapForFiles(files)

					_, res, specError := a.AnalyzeSpecification(specFilePath, nil)
					if specError != nil {
						formattedResult := utils.FormatSpecError(*specError, filesLines)
						fmt.Print(formattedResult)
						return nil
					}

					passedChecks := make([]analyzer.CheckResult, 0)
					skippedChecks := make([]analyzer.CheckResult, 0)
					// Print results for failed checks
					for _, result := range res {
						if result.Status == analyzer.CheckFailed {
							formattedResult := utils.FormatCheckResult(result, filesLines)
							fmt.Print(formattedResult)
						} else if result.Status == analyzer.CheckSkipped {
							skippedChecks = append(skippedChecks, result)
						} else if result.Status == analyzer.CheckPassed {
							passedChecks = append(passedChecks, result)
						}
					}

					if len(passedChecks)+len(skippedChecks) == len(res) {
						fmt.Println("All checks passed!")
					}

					// Print results for successful checks if --all flag is set
					if c.Bool("all") {
						for _, result := range passedChecks {
							formattedResult := utils.FormatCheckResult(result, filesLines)
							fmt.Print(formattedResult)
						}
					}

					// Print results for skipped checks if --all or --skipped flag is set
					if c.Bool("all") || c.Bool("skipped") {
						for _, result := range skippedChecks {
							formattedResult := utils.FormatCheckResult(result, filesLines)
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
				Hidden:    true,
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
					srv := server.CreateServer(c.Int("port"))

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
