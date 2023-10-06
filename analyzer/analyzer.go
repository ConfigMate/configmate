package analyzer

import (
	"fmt"

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

type Check interface {
	Check(values []interface{}) (passed bool, comment string, err error)
	// TODO: move to different file and define proper types
}

type AnalyzerImpl struct {
	checks map[string]Check // map of available checks
}

func (a *AnalyzerImpl) AnalyzeConfigFiles(files map[string]parsers.ConfigFile, rules []Rule) (res []Result, err error) {
	// Check rules
	for _, rule := range rules {
		args := make([]interface{}, 0)
		allValueTokenLocations := make([]TokenLocation, 0)
		// Get values for arguments
		for _, arg := range rule.Args {
			if isFile, alias, key := decodeFileValue(arg); isFile {
				// Get value from config file
				value, err := getNodeFromConfigFileNode(files[alias], key)
				if err != nil {
					return nil, fmt.Errorf("failed to get value from config file %s: %s", alias, err.Error())
				}
				allValueTokenLocations = append(allValueTokenLocations, TokenLocation{
					File:   alias,
					Line:   value.ValueLocation.Line,
					Column: value.ValueLocation.Column,
					Length: value.ValueLocation.Length,
				})
				args = append(args, value.Value)
			} else if isLiteral, value := decodeLiteralValue(arg); isLiteral {
				args = append(args, value)

			} else {
				return nil, fmt.Errorf("failed to decode argument %s", arg)
			}
		}

		// Apply check
		passed, comment, err := a.checks[rule.CheckName].Check(args)
		if err != nil {
			return nil, fmt.Errorf("failed to apply check %s: %s", rule.CheckName, err.Error())
		}

		// Add result
		res = append(res, Result{
			Passed:        passed,
			ResultComment: fmt.Sprintf("%s:%s", rule.CheckName, comment),
			TokenList:     allValueTokenLocations,
		})
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
