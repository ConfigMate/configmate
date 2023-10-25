package analyzer

// // Valid Setup
// func TestAnalyzeConfigFiles_ValidRuleArgument(t *testing.T) {
// 	// Create check mock
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	mockCheck := NewMockCheck(ctrl)

// 	// Define behavior of mock
// 	mockCheck.EXPECT().Check(gomock.Any()).Return(true, "Mock check passed", nil).AnyTimes()
// 	mockCheck.EXPECT().GetArgsSourceAndTypes().Return(
// 		[]CheckArgSource{CheckArgSource(File)},
// 		[]CheckArgType{CheckArgType(Int)},
// 	).AnyTimes()

// 	// Test configFile
// 	configFile := &parsers.Node{
// 		Type: parsers.Object,
// 		Value: map[string]*parsers.Node{
// 			"server": {
// 				Type: parsers.Object,
// 				Value: map[string]*parsers.Node{
// 					"port": {
// 						Type:  parsers.Int,
// 						Value: 8080,
// 					},
// 				},
// 			},
// 		},
// 	}

// 	// Test configFilesMap
// 	files := map[string]*parsers.Node{
// 		"test": configFile,
// 	}

// 	// Test rules
// 	rules := []Rule{
// 		{
// 			Description: "Sample rule",
// 			CheckName:   "sampleCheck",
// 			Args:        []string{"f:i:test.server.port"},
// 		},
// 	}

// 	// Setup
// 	analyzer := &analyzerImpl{
// 		checks: map[string]Check{
// 			"sampleCheck": mockCheck,
// 		},
// 	}

// 	// Execute
// 	res, err := analyzer.AnalyzeConfigFiles(files, rules)

// 	// Assert
// 	assert.NoError(t, err)
// 	assert.Len(t, res, 1)
// 	assert.True(t, res[0].Passed)
// 	assert.Contains(t, res[0].ResultComment, "sampleCheck: Mock check passed")
// }

// // Invalid: Field does not exist
// func TestAnalyzeConfigFiles_InvalidRuleArgument(t *testing.T) {
// 	// Create check mock
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	mockCheck := NewMockCheck(ctrl)

// 	// Define behavior of mock
// 	mockCheck.EXPECT().Check(gomock.Any()).Return(true, "Mock check passed", nil).AnyTimes()
// 	mockCheck.EXPECT().GetArgsSourceAndTypes().Return(
// 		[]CheckArgSource{CheckArgSource(File)},
// 		[]CheckArgType{CheckArgType(Int)},
// 	).AnyTimes()

// 	// Test configFile
// 	configFile := &parsers.Node{
// 		Type: parsers.Object,
// 		Value: map[string]*parsers.Node{
// 			"server": {
// 				Type: parsers.Object,
// 				Value: map[string]*parsers.Node{
// 					"port": {
// 						Type:  parsers.Int,
// 						Value: 8080,
// 					},
// 				},
// 			},
// 		},
// 	}

// 	// Test configFilesMap
// 	files := map[string]*parsers.Node{
// 		"test": configFile,
// 	}

// 	// Test rules
// 	rules := []Rule{
// 		{
// 			Description: "Sample rule",
// 			CheckName:   "sampleCheck",
// 			Args:        []string{"f:i:test.server.port[0]"},
// 		},
// 	}

// 	// Setup
// 	analyzer := &analyzerImpl{
// 		checks: map[string]Check{
// 			"sampleCheck": mockCheck,
// 		},
// 	}

// 	// Execute
// 	res, err := analyzer.AnalyzeConfigFiles(files, rules)

// 	// Assert
// 	assert.NoError(t, err)
// 	assert.Len(t, res, 1)
// 	assert.False(t, res[0].Passed)
// 	assert.Contains(t, res[0].ResultComment, "Value at server.port[0] in file test could not be found")
// }

// // Invalid: Field is not of the correct type
// func TestAnalyzeConfigFiles_InvalidRuleArgumentType(t *testing.T) {
// 	// Create check mock
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	mockCheck := NewMockCheck(ctrl)

// 	// Define behavior of mock
// 	mockCheck.EXPECT().Check(gomock.Any()).Return(true, "Mock check passed", nil).AnyTimes()
// 	mockCheck.EXPECT().GetArgsSourceAndTypes().Return(
// 		[]CheckArgSource{CheckArgSource(File)},
// 		[]CheckArgType{CheckArgType(Int)},
// 	).AnyTimes()

// 	// Test configFile
// 	configFile := &parsers.Node{
// 		Type: parsers.Object,
// 		Value: map[string]*parsers.Node{
// 			"server": {
// 				Type: parsers.Object,
// 				Value: map[string]*parsers.Node{
// 					"port": {
// 						Type:  parsers.String,
// 						Value: "8080",
// 					},
// 				},
// 			},
// 		},
// 	}

// 	// Test configFilesMap
// 	files := map[string]*parsers.Node{
// 		"test": configFile,
// 	}

// 	// Test rules
// 	rules := []Rule{
// 		{
// 			Description: "Sample rule",
// 			CheckName:   "sampleCheck",
// 			Args:        []string{"f:i:test.server.port"},
// 		},
// 	}

// 	// Setup
// 	analyzer := &analyzerImpl{
// 		checks: map[string]Check{
// 			"sampleCheck": mockCheck,
// 		},
// 	}

// 	// Execute
// 	res, err := analyzer.AnalyzeConfigFiles(files, rules)

// 	// Assert
// 	assert.NoError(t, err)
// 	assert.Len(t, res, 1)
// 	assert.False(t, res[0].Passed)
// 	assert.Contains(t, res[0].ResultComment, "Value at server.port in file test must be a int, got string")
// }
