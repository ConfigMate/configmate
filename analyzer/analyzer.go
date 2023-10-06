package analyzer

import (
	"fmt"
	"strconv"

	"github.com/ConfigMate/configmate/parsers"
)

const rulebookFileAlias = "rulebook"

type Analyzer interface {
	AnalyzeRuleBook(rulebookFileTree parsers.ConfigFile) (res []Result)
	AnalyzeConfigFiles(Files map[string]parsers.ConfigFile, Rules []Rule) (res []Result, err error)
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
	checks map[string]Check // map of available checks
}

func (a *AnalyzerImpl) AnalyzeConfigFiles(files map[string]parsers.ConfigFile, rules []Rule) (res []Result, err error) {
	// Check rules
	for _, rule := range rules {
		errors := false                               // true if there were errors in rule arguments
		args := make([]interface{}, 0)                // list of arguments for the check
		allTokenLocations := make([]TokenLocation, 0) // list of token locations for the check

		// Get values for arguments
		for _, arg := range rule.Args {
			// Validate argument
			if err := CheckArg(arg).Valid(); err != nil {
				return nil, fmt.Errorf("failed to validate argument %s: %s", arg, err.Error())
			}

			// Get argument source, type and value
			argSource := CheckArg(arg).Source()
			argType := CheckArg(arg).Type()
			argValue := CheckArg(arg).Value()

			switch argSource {
			case File:
				// Get file alias and path
				fileAlias, path := decodeFileValue(argValue)

				// Get value
				if value, err := getNodeFromConfigFileNode(files[fileAlias], path); err != nil {
					res = append(res, Result{
						Passed:        false,
						ResultComment: fmt.Sprintf("Value at %s in file %s could not be found: %s", path, fileAlias, err.Error()),
						TokenList:     []TokenLocation{{File: fileAlias}},
					})
					errors = true
				} else {
					switch argType {
					case Int:
						// Verify value is an integer
						if value.Type != parsers.Int {
							res = append(res, Result{
								Passed:        false,
								ResultComment: fmt.Sprintf("Value at %s in file %s must be an int, got %s", path, fileAlias, value.Type.String()),
								TokenList: []TokenLocation{
									{
										File:   fileAlias,
										Line:   value.ValueLocation.Line,
										Column: value.ValueLocation.Column,
										Length: value.ValueLocation.Length,
									},
								},
							})
							errors = true
						} else if _, ok := value.Value.(int); !ok {
							res = append(res, Result{
								Passed:        false,
								ResultComment: fmt.Sprintf("Value at %s in file %s must be an int, got %s", path, fileAlias, value.Type.String()),
								TokenList: []TokenLocation{
									{
										File:   fileAlias,
										Line:   value.ValueLocation.Line,
										Column: value.ValueLocation.Column,
										Length: value.ValueLocation.Length,
									},
								},
							})
							errors = true
						}
					case Float:
						// Verify value is a float
						if value.Type != parsers.Float {
							res = append(res, Result{
								Passed:        false,
								ResultComment: fmt.Sprintf("Value at %s in file %s must be a float, got %s", path, fileAlias, value.Type.String()),
								TokenList: []TokenLocation{
									{
										File:   fileAlias,
										Line:   value.ValueLocation.Line,
										Column: value.ValueLocation.Column,
										Length: value.ValueLocation.Length,
									},
								},
							})
							errors = true
						} else if _, ok := value.Value.(float64); !ok {
							res = append(res, Result{
								Passed:        false,
								ResultComment: fmt.Sprintf("Value at %s in file %s must be a float, got %s", path, fileAlias, value.Type.String()),
								TokenList: []TokenLocation{
									{
										File:   fileAlias,
										Line:   value.ValueLocation.Line,
										Column: value.ValueLocation.Column,
										Length: value.ValueLocation.Length,
									},
								},
							})
							errors = true
						}
					case Bool:
						// Verify value is a bool
						if value.Type != parsers.Bool {
							res = append(res, Result{
								Passed:        false,
								ResultComment: fmt.Sprintf("Value at %s in file %s must be a bool, got %s", path, fileAlias, value.Type.String()),
								TokenList: []TokenLocation{
									{
										File:   fileAlias,
										Line:   value.ValueLocation.Line,
										Column: value.ValueLocation.Column,
										Length: value.ValueLocation.Length,
									},
								},
							})
							errors = true
						} else if _, ok := value.Value.(bool); !ok {
							res = append(res, Result{
								Passed:        false,
								ResultComment: fmt.Sprintf("Value at %s in file %s must be a bool, got %s", path, fileAlias, value.Type.String()),
								TokenList: []TokenLocation{
									{
										File:   fileAlias,
										Line:   value.ValueLocation.Line,
										Column: value.ValueLocation.Column,
										Length: value.ValueLocation.Length,
									},
								},
							})
							errors = true
						}
					case String:
						// Verify value is a string
						if value.Type != parsers.String {
							res = append(res, Result{
								Passed:        false,
								ResultComment: fmt.Sprintf("Value at %s in file %s must be a string, got %s", path, fileAlias, value.Type.String()),
								TokenList: []TokenLocation{
									{
										File:   fileAlias,
										Line:   value.ValueLocation.Line,
										Column: value.ValueLocation.Column,
										Length: value.ValueLocation.Length,
									},
								},
							})
							errors = true
						} else if _, ok := value.Value.(string); !ok {
							res = append(res, Result{
								Passed:        false,
								ResultComment: fmt.Sprintf("Value at %s in file %s must be a string, got %s", path, fileAlias, value.Type.String()),
								TokenList: []TokenLocation{
									{
										File:   fileAlias,
										Line:   value.ValueLocation.Line,
										Column: value.ValueLocation.Column,
										Length: value.ValueLocation.Length,
									},
								},
							})
							errors = true
						}
					}

					// Add value
					args = append(args, value.Value)

					// Add token location
					allTokenLocations = append(allTokenLocations, TokenLocation{
						File:   fileAlias,
						Line:   value.ValueLocation.Line,
						Column: value.ValueLocation.Column,
						Length: value.ValueLocation.Length,
					})
				}

			case Literal:
				// Get value
				switch argType {
				case Int:
					// Verify value is an integer
					if value, err := strconv.Atoi(argValue); err != nil {
						res = append(res, Result{
							Passed:        false,
							ResultComment: fmt.Sprintf("Value %s must be an int, got %s", argValue, err.Error()),
							TokenList:     []TokenLocation{{File: rulebookFileAlias}},
						})
						errors = true
					} else {
						// Add value
						args = append(args, value)
					}
				case Float:
					// Verify value is a float
					if value, err := strconv.ParseFloat(argValue, 64); err != nil {
						res = append(res, Result{
							Passed:        false,
							ResultComment: fmt.Sprintf("Value %s must be a float, got %s", argValue, err.Error()),
							TokenList:     []TokenLocation{{File: rulebookFileAlias}},
						})
						errors = true
					} else {
						// Add value
						args = append(args, value)
					}
				case Bool:
					// Verify value is a bool
					if value, err := strconv.ParseBool(argValue); err != nil {
						res = append(res, Result{
							Passed:        false,
							ResultComment: fmt.Sprintf("Value %s must be a bool, got %s", argValue, err.Error()),
							TokenList:     []TokenLocation{{File: rulebookFileAlias}},
						})
						errors = true
					} else {
						// Add value
						args = append(args, value)
					}
				case String:
					// Add value
					args = append(args, argValue)
				}
			}
		}

		// Apply check if there were no errors getting arguments
		if !errors {
			// Apply check
			passed, comment, err := a.checks[rule.CheckName].Check(args)
			if err != nil {
				return nil, fmt.Errorf("failed to apply check %s: %s", rule.CheckName, err.Error())
			}

			// Add result
			res = append(res, Result{
				Passed:        passed,
				ResultComment: fmt.Sprintf("%s:%s", rule.CheckName, comment),
				TokenList:     allTokenLocations,
			})
		}

	}

	return res, nil
}

// AnalyzeRuleBook analyzes a rulebook file and find errors in it.
func (a *AnalyzerImpl) AnalyzeRuleBook(rulebookFileTree parsers.ConfigFile) (res []Result) {
	// Check name
	if name, err := getNodeFromConfigFileNode(rulebookFileTree, "name"); err != nil {
		res = append(res, Result{
			Passed:        false,
			ResultComment: fmt.Sprintf("Rulebook must have a name: %s", err.Error()),
			TokenList:     []TokenLocation{{File: rulebookFileAlias}},
		})
	} else if name.Type != parsers.String {
		res = append(res, Result{
			Passed:        false,
			ResultComment: fmt.Sprintf("Rulebook name must be a string, got %s", name.Type.String()),
			TokenList: []TokenLocation{
				{
					File:   rulebookFileAlias,
					Line:   name.ValueLocation.Line,
					Column: name.ValueLocation.Column,
					Length: name.ValueLocation.Length,
				},
			},
		})
	} else if name.Value.(string) == "" {
		res = append(res, Result{
			Passed:        false,
			ResultComment: "Rulebook name must not be empty",
			TokenList: []TokenLocation{
				{
					File:   rulebookFileAlias,
					Line:   name.ValueLocation.Line,
					Column: name.ValueLocation.Column,
					Length: name.ValueLocation.Length,
				},
			},
		})
	}

	// Check description
	if description, err := getNodeFromConfigFileNode(rulebookFileTree, "description"); err != nil {
		res = append(res, Result{
			Passed:        false,
			ResultComment: fmt.Sprintf("Rulebook must have a description: %s", err.Error()),
			TokenList:     []TokenLocation{{File: rulebookFileAlias}},
		})
	} else if description.Type != parsers.String {
		res = append(res, Result{
			Passed:        false,
			ResultComment: fmt.Sprintf("Rulebook description must be a string, got %s", description.Type.String()),
			TokenList: []TokenLocation{
				{
					File:   rulebookFileAlias,
					Line:   description.ValueLocation.Line,
					Column: description.ValueLocation.Column,
					Length: description.ValueLocation.Length,
				},
			},
		})
	} else if description.Value.(string) == "" {
		res = append(res, Result{
			Passed:        false,
			ResultComment: "Rulebook description must not be empty",
			TokenList: []TokenLocation{
				{
					File:   rulebookFileAlias,
					Line:   description.ValueLocation.Line,
					Column: description.ValueLocation.Column,
					Length: description.ValueLocation.Length,
				},
			},
		})
	}

	// Check files
	if files, err := getNodeFromConfigFileNode(rulebookFileTree, "files"); err != nil {
		res = append(res, Result{
			Passed:        false,
			ResultComment: fmt.Sprintf("Rulebook must have a files section: %s", err.Error()),
			TokenList:     []TokenLocation{{File: rulebookFileAlias}},
		})
	} else if files.Type != parsers.Object {
		res = append(res, Result{
			Passed:        false,
			ResultComment: fmt.Sprintf("Rulebook files section must be an object(map), got %s", files.Type.String()),
			TokenList: []TokenLocation{
				{
					File:   rulebookFileAlias,
					Line:   files.ValueLocation.Line,
					Column: files.ValueLocation.Column,
					Length: files.ValueLocation.Length,
				},
			},
		})
	} else if len(files.Value.(map[string]*parsers.Node)) == 0 {
		res = append(res, Result{
			Passed:        false,
			ResultComment: "Rulebook files section must not be empty",
			TokenList: []TokenLocation{
				{
					File:   rulebookFileAlias,
					Line:   files.ValueLocation.Line,
					Column: files.ValueLocation.Column,
					Length: files.ValueLocation.Length,
				},
			},
		})
	} else {
		// Check files
		for key, value := range files.Value.(map[string]*parsers.Node) {
			// Check key
			if key == "" {
				res = append(res, Result{
					Passed:        false,
					ResultComment: "File key must not be empty",
					TokenList: []TokenLocation{
						{
							File:   rulebookFileAlias,
							Line:   value.NameLocation.Line,
							Column: value.NameLocation.Column,
							Length: value.NameLocation.Length,
						},
					},
				})
			}

			// Check value
			if value.Type != parsers.String {
				res = append(res, Result{
					Passed:        false,
					ResultComment: fmt.Sprintf("File value must be a string, got %s", value.Type.String()),
					TokenList: []TokenLocation{
						{
							File:   rulebookFileAlias,
							Line:   value.ValueLocation.Line,
							Column: value.ValueLocation.Column,
							Length: value.ValueLocation.Length,
						},
					},
				})
			} else if value.Value.(string) == "" {
				res = append(res, Result{
					Passed:        false,
					ResultComment: "File value must not be empty",
					TokenList: []TokenLocation{
						{
							File:   rulebookFileAlias,
							Line:   value.ValueLocation.Line,
							Column: value.ValueLocation.Column,
							Length: value.ValueLocation.Length,
						},
					},
				})
			}
		}
	}

	// Check rules
	if rules, err := getNodeFromConfigFileNode(rulebookFileTree, "rules"); err != nil {
		res = append(res, Result{
			Passed:        false,
			ResultComment: fmt.Sprintf("Rulebook must have a rules section: %s", err.Error()),
			TokenList:     []TokenLocation{{File: rulebookFileAlias}},
		})
	} else if rules.Type != parsers.Array {
		res = append(res, Result{
			Passed:        false,
			ResultComment: fmt.Sprintf("Rules section must be an array, got %s", rules.Type.String()),
			TokenList: []TokenLocation{
				{
					File:   rulebookFileAlias,
					Line:   rules.ValueLocation.Line,
					Column: rules.ValueLocation.Column,
					Length: rules.ValueLocation.Length,
				},
			},
		})
	} else if len(rules.Value.([]*parsers.Node)) == 0 {
		res = append(res, Result{
			Passed:        false,
			ResultComment: "Rules section must not be empty",
			TokenList: []TokenLocation{
				{
					File:   rulebookFileAlias,
					Line:   rules.ValueLocation.Line,
					Column: rules.ValueLocation.Column,
					Length: rules.ValueLocation.Length,
				},
			},
		})
	} else {
		// Check rule by rule
		for _, rule := range rules.Value.([]*parsers.Node) {
			// Check rule
			if rule.Type != parsers.Object {
				res = append(res, Result{
					Passed:        false,
					ResultComment: fmt.Sprintf("Rule must be an object(map), got %s", rule.Type.String()),
					TokenList: []TokenLocation{
						{
							File:   rulebookFileAlias,
							Line:   rule.ValueLocation.Line,
							Column: rule.ValueLocation.Column,
							Length: rule.ValueLocation.Length,
						},
					},
				})
			} else {
				// Check rules description
				if description, err := getNodeFromConfigFileNode(rule, "description"); err != nil {
					res = append(res, Result{
						Passed:        false,
						ResultComment: fmt.Sprintf("Rule must have a description: %s", err.Error()),
						TokenList:     []TokenLocation{{File: rulebookFileAlias}},
					})
				} else if description.Type != parsers.String {
					res = append(res, Result{
						Passed:        false,
						ResultComment: fmt.Sprintf("Rule description must be a string, got %s", description.Type.String()),
						TokenList: []TokenLocation{
							{
								File:   rulebookFileAlias,
								Line:   description.ValueLocation.Line,
								Column: description.ValueLocation.Column,
								Length: description.ValueLocation.Length,
							},
						},
					})
				} else if description.Value.(string) == "" {
					res = append(res, Result{
						Passed:        false,
						ResultComment: "Rule description must not be empty",
						TokenList: []TokenLocation{
							{
								File:   rulebookFileAlias,
								Line:   description.ValueLocation.Line,
								Column: description.ValueLocation.Column,
								Length: description.ValueLocation.Length,
							},
						},
					})
				}

				// Check rules check name
				if checkName, err := getNodeFromConfigFileNode(rule, "check"); err != nil {
					res = append(res, Result{
						Passed:        false,
						ResultComment: fmt.Sprintf("Rule must have a check: %s", err.Error()),
						TokenList: []TokenLocation{
							{
								File:   rulebookFileAlias,
								Line:   rule.ValueLocation.Line,
								Column: rule.ValueLocation.Column,
								Length: rule.ValueLocation.Length,
							},
						},
					})
				} else if checkName.Type != parsers.String {
					res = append(res, Result{
						Passed:        false,
						ResultComment: fmt.Sprintf("Rule check must be a string, got %s", checkName.Type.String()),
						TokenList: []TokenLocation{
							{
								File:   rulebookFileAlias,
								Line:   checkName.ValueLocation.Line,
								Column: checkName.ValueLocation.Column,
								Length: checkName.ValueLocation.Length,
							},
						},
					})
				} else if checkName.Value.(string) == "" {
					res = append(res, Result{
						Passed:        false,
						ResultComment: "Rule check must not be empty",
						TokenList: []TokenLocation{
							{
								File:   rulebookFileAlias,
								Line:   checkName.ValueLocation.Line,
								Column: checkName.ValueLocation.Column,
								Length: checkName.ValueLocation.Length,
							},
						},
					})
				} else if _, ok := a.checks[checkName.Value.(string)]; !ok {
					res = append(res, Result{
						Passed:        false,
						ResultComment: fmt.Sprintf("Rule check %s does not exist", checkName.Value.(string)),
						TokenList: []TokenLocation{
							{
								File:   rulebookFileAlias,
								Line:   checkName.ValueLocation.Line,
								Column: checkName.ValueLocation.Column,
								Length: checkName.ValueLocation.Length,
							},
						},
					})
				} else {
					// // Get desired args
					// desiredArgs := a.checks[checkName.Value.(string)].GetArgTypes()

					// Check rules args
					if args, err := getNodeFromConfigFileNode(rule, "args"); err != nil {
						res = append(res, Result{
							Passed:        false,
							ResultComment: fmt.Sprintf("Rulebook rule must have args: %s", err.Error()),
							TokenList: []TokenLocation{
								{
									File:   rulebookFileAlias,
									Line:   rule.ValueLocation.Line,
									Column: rule.ValueLocation.Column,
									Length: rule.ValueLocation.Length,
								},
							},
						})
					} else if args.Type != parsers.Array {
						res = append(res, Result{
							Passed:        false,
							ResultComment: fmt.Sprintf("Rulebook rule args must be an array, got %s", args.Type.String()),
							TokenList: []TokenLocation{
								{
									File:   rulebookFileAlias,
									Line:   args.ValueLocation.Line,
									Column: args.ValueLocation.Column,
									Length: args.ValueLocation.Length,
								},
							},
						})
						// } else if len(args.Value.([]*parsers.Node)) != len(desiredArgs) {
						// 	res = append(res, Result{
						// 		Passed:        false,
						// 		ResultComment: fmt.Sprintf("Rulebook rule args must have %d args, got %d", len(desiredArgs), len(args.Value.([]*parsers.Node))),
						// 		TokenList: []TokenLocation{
						// 			{
						// 				File:   rulebookFileAlias,
						// 				Line:   args.ValueLocation.Line,
						// 				Column: args.ValueLocation.Column,
						// 				Length: args.ValueLocation.Length,
						// 			},
						// 		},
						// 	})
						// } else {
						// 	// Check args
						// 	for i, arg := range args.Value.([]*parsers.Node) {
						// 		// Check arg
						// 		if arg.Type != desiredArgs[i] {
						// 			res = append(res, Result{
						// 				Passed:        false,
						// 				ResultComment: fmt.Sprintf("Rulebook rule arg %d must be of type %s, got %s", i, desiredArgs[i], arg.Type.String()),
						// 				TokenList: []TokenLocation{
						// 					{
						// 						File:   rulebookFileAlias,
						// 						Line:   arg.ValueLocation.Line,
						// 						Column: arg.ValueLocation.Column,
						// 						Length: arg.ValueLocation.Length,
						// 					},
						// 				},
						// 			})
						// 		}
						// 	}
					}
				}
			}
		}
	}

	return res
}
