package analyzer

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ConfigMate/configmate/analyzer/check"
	"github.com/ConfigMate/configmate/analyzer/spec"
	"github.com/ConfigMate/configmate/analyzer/types"
	"github.com/ConfigMate/configmate/files"
	"github.com/ConfigMate/configmate/parsers"
)

type Analyzer interface {
	AnalyzeSpecification(specFilePath string, specFileContent []byte) (*spec.Specification, []CheckResult, *SpecError)
	AllFilesContent(specFilePath string) map[string][]byte
}

type SpecError struct {
	AnalyzerMsg string                  `json:"analyzer_msg"`
	ErrorMsgs   []string                `json:"error_msgs"`
	TokenList   []TokenLocationWithFile `json:"token_list"` // list of tokens involved in the error
}

type CheckStatus int

const (
	CheckPassed CheckStatus = iota
	CheckFailed
	CheckSkipped
)

type CheckResult struct {
	Status        CheckStatus             `json:"status"`         // whether the check passed, failed or was skipped
	ResultComment string                  `json:"result_comment"` // an error msg or comment about the result
	Field         spec.FieldSpec          `json:"field"`          // the rule that was checked
	CheckNum      int                     `json:"check_num"`      // the number of the check that was evaluated
	TokenList     []TokenLocationWithFile `json:"token_list"`     // list of tokens involved in the check
}

// TokenLocationWithFile is a TokenLocation enhanced with a file path;
// representing the file the token is in.
type TokenLocationWithFile struct {
	File     string                `json:"file"`
	Location parsers.TokenLocation `json:"location"`
}

// mainFileAlias is the alias used to reference the main config file.
const mainFileAlias = "main"

type analyzerImpl struct {
	specParser     spec.SpecParser
	checkEvaluator check.CheckEvaluator

	fileFetcher    files.FileFetcher
	parserProvider parsers.ParserProvider
}

func NewAnalyzer(
	specParser spec.SpecParser,
	checkEvaluator check.CheckEvaluator,
	fileFetcher files.FileFetcher,
	parserProvider parsers.ParserProvider) Analyzer {
	return &analyzerImpl{
		specParser:     specParser,
		checkEvaluator: checkEvaluator,
		fileFetcher:    fileFetcher,
		parserProvider: parserProvider,
	}
}

func (a *analyzerImpl) AnalyzeSpecification(specFilePath string, specFileContent []byte) (*spec.Specification, []CheckResult, *SpecError) {
	// Check if contents were not provided, and get them from the file path then
	if specFileContent == nil {
		var err error
		specFileContent, err = a.fileFetcher.FetchFile(specFilePath)
		if err != nil {
			return nil, nil, &SpecError{
				AnalyzerMsg: "Failed to get specification file",
				ErrorMsgs:   []string{err.Error()},
			}
		}
	}

	// Parse specification from contents
	mainSpec, parserErrors := a.specParser.Parse(specFileContent)
	if len(parserErrors) > 0 { // Check for parser errors
		specError := &SpecError{
			AnalyzerMsg: "Failed to parse specification file",
			ErrorMsgs:   []string{},
			TokenList:   []TokenLocationWithFile{},
		}
		for _, parserError := range parserErrors {
			specError.ErrorMsgs = append(specError.ErrorMsgs, parserError.ErrorMessage)
			specError.TokenList = append(specError.TokenList, TokenLocationWithFile{
				File:     specFilePath,
				Location: parserError.Location,
			})
		}
		return nil, nil, specError
	}

	// Create fields map
	fields := make(map[string][]spec.FieldSpec)
	fields[mainFileAlias] = mainSpec.Fields

	// Get main config file
	mainConfigContent, err := a.fileFetcher.FetchFile(mainSpec.File)
	if err != nil {
		specError := &SpecError{
			AnalyzerMsg: "Failed to get main config file",
			ErrorMsgs:   []string{err.Error()},
			TokenList: []TokenLocationWithFile{
				{
					File:     specFilePath,
					Location: mainSpec.FileLocation,
				},
			},
		}
		return mainSpec, nil, specError
	}

	// Get parser for main config file
	mainConfigParser, err := a.parserProvider.GetParser(mainSpec.FileFormat)
	if err != nil {
		specError := &SpecError{
			AnalyzerMsg: "Failed to get parser for main config file",
			ErrorMsgs:   []string{err.Error()},
			TokenList: []TokenLocationWithFile{
				{
					File:     specFilePath,
					Location: mainSpec.FileFormatLocation,
				},
			},
		}
		return mainSpec, nil, specError
	}

	// Parse main config file
	mainConfig, parserErrs := mainConfigParser.Parse(mainConfigContent)
	if len(parserErrs) > 0 {
		specError := &SpecError{
			AnalyzerMsg: "Failed to parse main config file",
			ErrorMsgs:   []string{},
			TokenList: []TokenLocationWithFile{
				{
					File:     specFilePath,
					Location: mainSpec.FileLocation,
				},
			},
		}
		for _, parserError := range parserErrs {
			specError.ErrorMsgs = append(specError.ErrorMsgs, parserError.Message)
			specError.TokenList = append(specError.TokenList, TokenLocationWithFile{
				File:     mainSpec.File,
				Location: parserError.Location,
			})
		}
		return mainSpec, nil, specError
	}

	// Create file paths maps for token locations in errors
	specFilePaths := make(map[string]string)
	configFilePaths := make(map[string]string)

	// Add main spec and config file path map
	specFilePaths[mainFileAlias] = specFilePath
	configFilePaths[mainFileAlias] = mainSpec.File

	// Create files map
	files := make(map[string]*parsers.Node)
	files[mainFileAlias] = mainConfig

	// Fetch imported spec files
	for alias, importedSpecFilePath := range mainSpec.Imports {
		// Check that alias doesn't conflict with main spec
		if alias == mainFileAlias {
			specError := &SpecError{
				AnalyzerMsg: "Alias conflicts: '" + mainFileAlias + "' is reserved and used internally",
				ErrorMsgs:   []string{},
				TokenList: []TokenLocationWithFile{
					{
						File:     specFilePath,
						Location: mainSpec.ImportsAliasLocation[alias],
					},
				},
			}
			return mainSpec, nil, specError
		}

		// Check that alias doesn't conflict with other imported specs
		if _, ok := files[alias]; ok {
			specError := &SpecError{
				AnalyzerMsg: fmt.Sprintf("Alias conflicts: '%s' is already used for another imported spec", alias),
				ErrorMsgs:   []string{},
				TokenList: []TokenLocationWithFile{
					{
						File:     specFilePath,
						Location: mainSpec.ImportsAliasLocation[alias],
					},
				},
			}
			return mainSpec, nil, specError
		}

		// Get imported spec file
		importedSpecBytes, err := a.fileFetcher.FetchFile(importedSpecFilePath)
		if err != nil {
			specError := &SpecError{
				AnalyzerMsg: "Failed to get imported specification file",
				ErrorMsgs:   []string{err.Error()},
				TokenList: []TokenLocationWithFile{
					{
						File:     specFilePath,
						Location: mainSpec.ImportsLocation[alias],
					},
				},
			}
			return mainSpec, nil, specError
		}

		// Parse imported spec file
		importedSpec, parserErrors := a.specParser.Parse(importedSpecBytes)
		if len(parserErrors) > 0 {
			specError := &SpecError{
				AnalyzerMsg: fmt.Sprintf("Failed to parse imported spec file %s", importedSpecFilePath),
				ErrorMsgs:   []string{},
				TokenList: []TokenLocationWithFile{
					{
						File:     specFilePath,
						Location: mainSpec.ImportsLocation[alias],
					},
				},
			}
			for _, parserError := range parserErrors {
				specError.ErrorMsgs = append(specError.ErrorMsgs, parserError.ErrorMessage)
				specError.TokenList = append(specError.TokenList, TokenLocationWithFile{
					File:     importedSpecFilePath,
					Location: parserError.Location,
				})
			}
			return mainSpec, nil, specError
		}

		// Add imported spec file to fields map
		fields[alias] = importedSpec.Fields

		// Get imported config file
		importedConfigContent, err := a.fileFetcher.FetchFile(importedSpec.File)
		if err != nil {
			specError := &SpecError{
				AnalyzerMsg: fmt.Sprintf("Failed to get imported config file from spec %s", alias),
				ErrorMsgs:   []string{err.Error()},
				TokenList: []TokenLocationWithFile{
					{
						File:     specFilePath,
						Location: mainSpec.ImportsLocation[alias],
					},
					{
						File:     importedSpecFilePath,
						Location: importedSpec.FileLocation,
					},
				},
			}
			return mainSpec, nil, specError
		}

		// Get parser for imported config file
		importedConfigParser, err := a.parserProvider.GetParser(importedSpec.FileFormat)
		if err != nil {
			specError := &SpecError{
				AnalyzerMsg: fmt.Sprintf("Failed to get parser for imported config file from spec %s", alias),
				ErrorMsgs:   []string{err.Error()},
				TokenList: []TokenLocationWithFile{
					{
						File:     specFilePath,
						Location: mainSpec.ImportsLocation[alias],
					},
					{
						File:     importedSpecFilePath,
						Location: importedSpec.FileFormatLocation,
					},
				},
			}
			return mainSpec, nil, specError
		}

		// Parse imported config file
		importedConfig, parserErrs := importedConfigParser.Parse(importedConfigContent)
		if len(parserErrs) > 0 {
			specError := &SpecError{
				AnalyzerMsg: fmt.Sprintf("Failed to parse imported config file from spec %s", alias),
				ErrorMsgs:   []string{err.Error()},
				TokenList: []TokenLocationWithFile{
					{
						File:     specFilePath,
						Location: mainSpec.ImportsLocation[alias],
					},
					{
						File:     importedSpecFilePath,
						Location: importedSpec.FileLocation,
					},
				},
			}
			for _, parserError := range parserErrs {
				specError.ErrorMsgs = append(specError.ErrorMsgs, parserError.Message)
				specError.TokenList = append(specError.TokenList, TokenLocationWithFile{
					File:     importedSpec.File,
					Location: parserError.Location,
				})
			}
			return mainSpec, nil, specError
		}

		// Add imported config file to files map
		files[alias] = importedConfig

		// Add imported spec and config file to path map
		specFilePaths[alias] = importedSpecFilePath
		configFilePaths[alias] = importedSpec.File
	}

	// Find all fields and parse them
	// optMissingFields is a map of optional fields that are missing
	fieldValues, fieldLocations, optMissingFields, specError := a.findAndParseAllFields(
		files,
		fields,
		specFilePaths,
		configFilePaths,
	)
	if specError != nil {
		return mainSpec, nil, specError
	}

	// Run checks
	res, specError := a.runChecks(
		mainSpec.Fields,
		fieldValues,
		fieldLocations,
		optMissingFields,
		specFilePaths,
	)
	if specError != nil {
		return mainSpec, nil, specError
	}

	// Analyze config files
	return mainSpec, res, nil
}

func (a *analyzerImpl) AllFilesContent(specFilePath string) map[string][]byte {
	// Create files map
	files := make(map[string][]byte)

	// Get specification file
	specBytes, err := a.fileFetcher.FetchFile(specFilePath)
	if err != nil {
		return nil
	}

	// Add spec file to files map
	files[specFilePath] = specBytes

	// Parse specification
	spec, parserErrors := a.specParser.Parse(specBytes)
	if len(parserErrors) > 0 {
		return files
	}

	// Get config file
	configBytes, err := a.fileFetcher.FetchFile(spec.File)
	if err != nil {
		return files
	}

	// Add config file to files map
	files[spec.File] = configBytes

	// Fetch imported spec files
	for _, importedSpecFilePath := range spec.Imports {
		// Get imported spec file
		importedSpecBytes, err := a.fileFetcher.FetchFile(importedSpecFilePath)
		if err != nil {
			continue
		}

		// Add imported spec file to files map
		files[importedSpecFilePath] = importedSpecBytes

		// Parse imported spec file
		importedSpec, parserErrors := a.specParser.Parse(importedSpecBytes)
		if len(parserErrors) > 0 {
			continue
		}

		// Get imported config file
		importedConfigBytes, err := a.fileFetcher.FetchFile(importedSpec.File)
		if err != nil {
			continue
		}

		// Add imported config file to files map
		files[importedSpec.File] = importedConfigBytes
	}

	return files
}

func (a *analyzerImpl) findAndParseAllFields(
	files map[string]*parsers.Node,
	fields map[string][]spec.FieldSpec,
	specFilePaths map[string]string,
	configFilePaths map[string]string) (map[string]types.IType, map[string]TokenLocationWithFile, map[string]bool, *SpecError) {
	// Create maps
	fieldValues := make(map[string]types.IType)
	fieldLocations := make(map[string]TokenLocationWithFile)
	optMissingFields := make(map[string]bool)

	for fileAlias, fileFields := range fields {
		// Sort file specs by field name lenght (shortest first)
		// This guarantees parent fields are checked before child fields
		sort.Slice(fileFields, func(i, j int) bool {
			return len(fileFields[i].Field.String()) < len(fileFields[j].Field.String())
		})

		for _, fspec := range fileFields {
			// Check if a parent field is an optional
			// field that is missing, which makes the
			// current field optional as well
			for optMissingField := range optMissingFields {
				if strings.HasPrefix(fileAlias+"."+fspec.Field.String(), optMissingField) {
					optMissingFields[fileAlias+"."+fspec.Field.String()] = true
					break
				}
			}
			if optMissingFields[fileAlias+"."+fspec.Field.String()] {
				continue
			}

			// Get field from file tree
			fnode, err := files[fileAlias].Get(fspec.Field)
			if err != nil {
				return nil, nil, nil, &SpecError{
					AnalyzerMsg: fmt.Sprintf("Failed to get field %s from file %s", fspec.Field, configFilePaths[fileAlias]),
					ErrorMsgs:   []string{err.Error()},
					TokenList: []TokenLocationWithFile{
						{
							File:     specFilePaths[fileAlias],
							Location: fspec.FieldLocation,
						},
					},
				}
			} else if fnode == nil && !fspec.Optional { // Field not found and not optial
				return nil, nil, nil, &SpecError{
					AnalyzerMsg: fmt.Sprintf("Field %s not found in file %s", fspec.Field, configFilePaths[fileAlias]),
					ErrorMsgs:   []string{},
					TokenList: []TokenLocationWithFile{
						{
							File:     specFilePaths[fileAlias],
							Location: fspec.FieldLocation,
						},
					},
				}
			} else if fnode == nil && fspec.Optional { // Field not found and optional
				optMissingFields[fileAlias+"."+fspec.Field.String()] = true
			} else { // Field found
				t, err := types.MakeType(fspec.Type, fnode.Value)
				if err != nil {
					return nil, nil, nil, &SpecError{
						AnalyzerMsg: fmt.Sprintf("failed to parse field %s from file %s as type %s",
							fspec.Field, configFilePaths[fileAlias], fspec.Type,
						),
						ErrorMsgs: []string{err.Error()},
						TokenList: []TokenLocationWithFile{
							{
								File:     specFilePaths[fileAlias],
								Location: fspec.TypeLocation,
							},
							{
								File:     configFilePaths[fileAlias],
								Location: fnode.ValueLocation,
							},
						},
					}
				} else {
					fieldValues[fileAlias+"."+fspec.Field.String()] = t
					fieldLocations[fileAlias+"."+fspec.Field.String()] = TokenLocationWithFile{
						File:     configFilePaths[fileAlias],
						Location: fnode.ValueLocation,
					}
				}
			}
		}
	}

	return fieldValues, fieldLocations, optMissingFields, nil
}

func (a *analyzerImpl) runChecks(
	mainFieldSpecs []spec.FieldSpec,
	fieldValues map[string]types.IType,
	fieldLocations map[string]TokenLocationWithFile,
	optMissingFields map[string]bool,
	specFilePaths map[string]string) (res []CheckResult, err *SpecError) {

	// Create results list
	res = []CheckResult{}

	for index, fspec := range mainFieldSpecs {
		// Evaluate checks
		for checkNum, checkInfo := range fspec.Checks {
			// Evaluate check
			result, skipping, err := a.checkEvaluator.Evaluate(
				checkInfo.Check,
				mainFileAlias+"."+fspec.Field.String(),
				fieldValues,
				optMissingFields,
			)
			if result == nil {
				return nil, &SpecError{
					AnalyzerMsg: fmt.Sprintf("failed to evaluate check %s for field %s", checkInfo.Check, fspec.Field),
					ErrorMsgs:   []string{err.Error()},
					TokenList: []TokenLocationWithFile{
						{
							File:     specFilePaths[mainFileAlias],
							Location: fspec.Checks[checkNum].Location,
						},
					},
				}
			} else if skipping {
				res = append(res, CheckResult{
					Status:        CheckSkipped,
					ResultComment: err.Error(),
					Field:         mainFieldSpecs[index],
					CheckNum:      checkNum,
					TokenList:     []TokenLocationWithFile{},
				})
			} else {
				resComment := ""
				if err != nil {
					resComment = err.Error()
				}

				checkStatus := CheckPassed
				if !result.Value().(bool) {
					checkStatus = CheckFailed
				}

				res = append(res, CheckResult{
					Status:        checkStatus,
					ResultComment: resComment,
					Field:         mainFieldSpecs[index],
					CheckNum:      checkNum,
					TokenList: []TokenLocationWithFile{
						fieldLocations[mainFileAlias+"."+fspec.Field.String()],
					},
				})
			}
		}
	}

	return res, nil
}
