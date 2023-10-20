package analyzer

import (
	"fmt"

	"github.com/ConfigMate/configmate/parsers"
)

const rulebookFileAlias = "rulebook"

type Analyzer interface {
	AnalyzeRuleBook(rulebookFileTree *parsers.Node) (res []Result)
	AnalyzeConfigFiles(Files map[string]*parsers.Node, Rules []Rule) (res []Result, err error)
}

type Result struct {
	Passed        bool            `json:"passed"`         // true if the check passed, false if it failed
	ResultComment string          `json:"result_comment"` // an error msg or comment about the result
	TokenList     []TokenLocation `json:"token_list"`     // a list of tokens that were involved in the rule
}

type TokenLocation struct {
	File  string `json:"file"`
	Start parsers.Location
	End   parsers.Location
}

type AnalyzerImpl struct {
	checks map[string]Check // map of available checks
}

func (a *AnalyzerImpl) AnalyzeConfigFiles(files map[string]*parsers.Node, rules []Rule) (res []Result, err error) {
	// Check rules
	for _, rule := range rules {
		errors := false                               // true if there were errors in rule arguments
		args := make([]interface{}, 0)                // list of arguments for the check
		allTokenLocations := make([]TokenLocation, 0) // list of token locations for the check

		// Get values for arguments
		for _, strArg := range rule.Args {
			// Parse argument
			arg, err := ParseCheckArg(strArg)
			if err != nil {
				return nil, fmt.Errorf("failed to parse argument %s: %s", strArg, err.Error())
			}

			switch arg.s { // Switch based on the argument source
			case File:
				fArg := arg.v.(FileValue)                                       // Cast value as *FileValue (unsafe)
				if value, err := files[fArg.alias].Get(fArg.path); err != nil { // Get value from file
					res = append(res, Result{
						Passed:        false,
						ResultComment: fmt.Sprintf("Value at %s in file %s could not be found: %s", fArg.path, fArg.alias, err.Error()),
						TokenList:     []TokenLocation{makeTOFTokenLocation(fArg.alias)},
					})
					errors = true
				} else if !equalType(value.Type, arg.t) { // Ensure value type is correct
					res = append(res, Result{
						Passed:        false,
						ResultComment: fmt.Sprintf("Value at %s in file %s must be a %s, got %s", fArg.path, fArg.alias, arg.t.String(), value.Type.String()),
						TokenList:     []TokenLocation{makeValueTokenLocation(fArg.alias, value)},
					})
					errors = true
				} else { // Add value to args and token location to allTokenLocations
					args = append(args, value.Value)
					allTokenLocations = append(allTokenLocations, makeValueTokenLocation(fArg.alias, value))
				}

			case Literal:
				args = append(args, arg.v)
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
				ResultComment: fmt.Sprintf("%s: %s", rule.CheckName, comment),
				TokenList:     allTokenLocations,
			})
		}

	}

	return res, nil
}

// AnalyzeRuleBook analyzes a rulebook file and find errors in it.
func (a *AnalyzerImpl) AnalyzeRuleBook(rulebookFileTree *parsers.Node) (res []Result) {
	// Check nameNode
	if nameNode, err := rulebookFileTree.Get("name"); err != nil {
		res = append(res, Result{
			Passed:        false,
			ResultComment: fmt.Sprintf("Rulebook must have a name: %s", err.Error()),
			TokenList:     []TokenLocation{makeTOFTokenLocation(rulebookFileAlias)},
		})
	} else {
		if nameNode.Type != parsers.String {
			res = append(res, Result{
				Passed:        false,
				ResultComment: fmt.Sprintf("Rulebook name must be a string, got %s", nameNode.Type.String()),
				TokenList:     []TokenLocation{makeValueTokenLocation(rulebookFileAlias, nameNode)},
			})
		} else if nameNode.Value.(string) == "" {
			res = append(res, Result{
				Passed:        false,
				ResultComment: "Rulebook name must not be empty",
				TokenList:     []TokenLocation{makeValueTokenLocation(rulebookFileAlias, nameNode)},
			})
		}
	}

	// Check description
	if description, err := rulebookFileTree.Get("description"); err != nil {
		res = append(res, Result{
			Passed:        false,
			ResultComment: fmt.Sprintf("Rulebook must have a description: %s", err.Error()),
			TokenList:     []TokenLocation{makeTOFTokenLocation(rulebookFileAlias)},
		})
	} else {
		if description.Type != parsers.String {
			res = append(res, Result{
				Passed:        false,
				ResultComment: fmt.Sprintf("Rulebook description must be a string, got %s", description.Type.String()),
				TokenList:     []TokenLocation{makeValueTokenLocation(rulebookFileAlias, description)},
			})
		} else if description.Value.(string) == "" {
			res = append(res, Result{
				Passed:        false,
				ResultComment: "Rulebook description must not be empty",
				TokenList:     []TokenLocation{makeValueTokenLocation(rulebookFileAlias, description)},
			})
		}
	}

	// Check files
	if files, err := rulebookFileTree.Get("files"); err != nil {
		res = append(res, Result{
			Passed:        false,
			ResultComment: fmt.Sprintf("Rulebook must have a files section: %s", err.Error()),
			TokenList:     []TokenLocation{makeTOFTokenLocation(rulebookFileAlias)},
		})
	} else {
		if files.Type != parsers.Object {
			res = append(res, Result{
				Passed:        false,
				ResultComment: fmt.Sprintf("Rulebook files section must be an object(map), got %s", files.Type.String()),
				TokenList:     []TokenLocation{makeValueTokenLocation(rulebookFileAlias, files)},
			})
		} else if len(files.Value.(map[string]*parsers.Node)) == 0 {
			res = append(res, Result{
				Passed:        false,
				ResultComment: "Rulebook files section must not be empty",
				TokenList:     []TokenLocation{makeValueTokenLocation(rulebookFileAlias, files)},
			})
		} else {
			// Check files
			for key, value := range files.Value.(map[string]*parsers.Node) {
				// Check key
				if key == "" {
					res = append(res, Result{
						Passed:        false,
						ResultComment: "File key must not be empty",
						TokenList:     []TokenLocation{makeNameTokenLocation(rulebookFileAlias, value)},
					})
				}

				// Check value
				if value.Type != parsers.String {
					res = append(res, Result{
						Passed:        false,
						ResultComment: fmt.Sprintf("File value must be a string, got %s", value.Type.String()),
						TokenList:     []TokenLocation{makeValueTokenLocation(rulebookFileAlias, value)},
					})
				} else if value.Value.(string) == "" {
					res = append(res, Result{
						Passed:        false,
						ResultComment: "File value must not be empty",
						TokenList:     []TokenLocation{makeValueTokenLocation(rulebookFileAlias, value)},
					})
				}
			}
		}
	}

	// Check rules
	if rules, err := rulebookFileTree.Get("rules"); err != nil {
		res = append(res, Result{
			Passed:        false,
			ResultComment: fmt.Sprintf("Rulebook must have a rules section: %s", err.Error()),
			TokenList:     []TokenLocation{makeTOFTokenLocation(rulebookFileAlias)},
		})
	} else {
		if rules.Type != parsers.Array {
			res = append(res, Result{
				Passed:        false,
				ResultComment: fmt.Sprintf("Rules section must be an array, got %s", rules.Type.String()),
				TokenList:     []TokenLocation{makeValueTokenLocation(rulebookFileAlias, rules)},
			})
		} else if len(rules.Value.([]*parsers.Node)) == 0 {
			res = append(res, Result{
				Passed:        false,
				ResultComment: "Rules section must not be empty",
				TokenList:     []TokenLocation{makeValueTokenLocation(rulebookFileAlias, rules)},
			})
		} else {
			// Check rule by rule
			for _, rule := range rules.Value.([]*parsers.Node) {
				// Check rule
				if rule.Type != parsers.Object {
					res = append(res, Result{
						Passed:        false,
						ResultComment: fmt.Sprintf("Rule must be an object(map), got %s", rule.Type.String()),
						TokenList:     []TokenLocation{makeValueTokenLocation(rulebookFileAlias, rule)},
					})
				} else {
					// Check rules description
					if description, err := rule.Get("description"); err != nil {
						res = append(res, Result{
							Passed:        false,
							ResultComment: fmt.Sprintf("Rule must have a description: %s", err.Error()),
							TokenList:     []TokenLocation{makeTOFTokenLocation(rulebookFileAlias)},
						})
					} else {
						if description.Type != parsers.String {
							res = append(res, Result{
								Passed:        false,
								ResultComment: fmt.Sprintf("Rule description must be a string, got %s", description.Type.String()),
								TokenList:     []TokenLocation{makeValueTokenLocation(rulebookFileAlias, description)},
							})
						} else if description.Value.(string) == "" {
							res = append(res, Result{
								Passed:        false,
								ResultComment: "Rule description must not be empty",
								TokenList:     []TokenLocation{makeValueTokenLocation(rulebookFileAlias, description)},
							})
						}
					}

					// Check rules check name
					if checkName, err := rule.Get("check"); err != nil {
						res = append(res, Result{
							Passed:        false,
							ResultComment: fmt.Sprintf("Rule must have a check: %s", err.Error()),
							TokenList:     []TokenLocation{makeValueTokenLocation(rulebookFileAlias, rule)},
						})
					} else {
						if checkName.Type != parsers.String {
							res = append(res, Result{
								Passed:        false,
								ResultComment: fmt.Sprintf("Rule check must be a string, got %s", checkName.Type.String()),
								TokenList:     []TokenLocation{makeValueTokenLocation(rulebookFileAlias, checkName)},
							})
						} else if checkName.Value.(string) == "" {
							res = append(res, Result{
								Passed:        false,
								ResultComment: "Rule check must not be empty",
								TokenList:     []TokenLocation{makeValueTokenLocation(rulebookFileAlias, checkName)},
							})
						} else if _, ok := a.checks[checkName.Value.(string)]; !ok {
							res = append(res, Result{
								Passed:        false,
								ResultComment: fmt.Sprintf("Rule check %s does not exist", checkName.Value.(string)),
								TokenList:     []TokenLocation{makeValueTokenLocation(rulebookFileAlias, checkName)},
							})
						} else {
							// Get desired args sources and types
							desiredSources, desiredTypes := a.checks[checkName.Value.(string)].GetArgsSourceAndTypes()

							// Check rules args
							if args, err := rule.Get("args"); err != nil {
								res = append(res, Result{
									Passed:        false,
									ResultComment: fmt.Sprintf("Rulebook rule must have args: %s", err.Error()),
									TokenList:     []TokenLocation{makeValueTokenLocation(rulebookFileAlias, rule)},
								})
							} else {
								if args.Type != parsers.Array {
									res = append(res, Result{
										Passed:        false,
										ResultComment: fmt.Sprintf("Rulebook rule args must be an array, got %s", args.Type.String()),
										TokenList:     []TokenLocation{makeValueTokenLocation(rulebookFileAlias, args)},
									})
								} else if len(args.Value.([]*parsers.Node)) != len(desiredSources) {
									res = append(res, Result{
										Passed:        false,
										ResultComment: fmt.Sprintf("Rulebook rule args must have %d args, got %d", len(desiredSources), len(args.Value.([]*parsers.Node))),
										TokenList:     []TokenLocation{makeValueTokenLocation(rulebookFileAlias, args)},
									})
								} else {
									// Check args
									for i, argNode := range args.Value.([]*parsers.Node) {
										if _, ok := argNode.Value.(string); !ok { // Ensure value is represented as a string
											res = append(res, Result{
												Passed:        false,
												ResultComment: fmt.Sprintf("All rule arguments must be represented with strings: %v must be a string, got %s", argNode.Value, argNode.Type.String()),
												TokenList:     []TokenLocation{makeValueTokenLocation(rulebookFileAlias, argNode)},
											})
										} else if arg, err := ParseCheckArg(argNode.Value.(string)); err != nil {
											res = append(res, Result{
												Passed:        false,
												ResultComment: fmt.Sprintf("Failed to parse argument %s: %s", argNode.Value.(string), err.Error()),
												TokenList:     []TokenLocation{makeValueTokenLocation(rulebookFileAlias, argNode)},
											})
										} else if arg.s != desiredSources[i] { // Ensure value source is correct
											res = append(res, Result{
												Passed:        false,
												ResultComment: fmt.Sprintf("Rulebook rule arg %d must be from %s, got %s", i, desiredSources[i].String(), arg.s.String()),
												TokenList:     []TokenLocation{makeValueTokenLocation(rulebookFileAlias, argNode)},
											})
										} else if arg.t != desiredTypes[i] { // Ensure value type is correct
											res = append(res, Result{
												Passed:        false,
												ResultComment: fmt.Sprintf("Rulebook rule arg %d must be of type %s, got %s", i, desiredTypes[i].String(), arg.t.String()),
												TokenList:     []TokenLocation{makeValueTokenLocation(rulebookFileAlias, argNode)},
											})
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return res
}
