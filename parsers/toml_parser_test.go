package parsers

import (
	"os"
	"reflect"
	"testing"
)

type tomlParserTestCase struct {
	input        []byte
	expected     *Node
	expectedErrs []CMParserError
}

func TestParseSimpleConfig_tomlParser(t *testing.T) {
	// Input
	testConfig, err := os.ReadFile("./test_configs/simple.toml")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	// Test cases
	testCases := []tomlParserTestCase{
		{
			input: testConfig,
			expected: &Node{
				Type: Object,
				Value: map[string]*Node{
					"title": {
						Type:          String,
						Value:         "TOML Sample Configuration",
						NameLocation:  TokenLocation{Start: CharLocation{Line: 0, Column: 0}, End: CharLocation{Line: 0, Column: 5}},
						ValueLocation: TokenLocation{Start: CharLocation{Line: 0, Column: 8}, End: CharLocation{Line: 0, Column: 35}},
					},
					"project_lead": {
						Type: Object,
						Value: map[string]*Node{
							"name": {
								Type:          String,
								Value:         "Charlie Example",
								NameLocation:  TokenLocation{Start: CharLocation{Line: 3, Column: 0}, End: CharLocation{Line: 3, Column: 4}},
								ValueLocation: TokenLocation{Start: CharLocation{Line: 3, Column: 7}, End: CharLocation{Line: 3, Column: 24}},
							},
							"birthdate": {
								Type:          String,
								Value:         "1980-11-15T09:45:00-05:00",
								NameLocation:  TokenLocation{Start: CharLocation{Line: 4, Column: 0}, End: CharLocation{Line: 4, Column: 9}},
								ValueLocation: TokenLocation{Start: CharLocation{Line: 4, Column: 12}, End: CharLocation{Line: 4, Column: 37}},
							},
						},
						NameLocation:  TokenLocation{Start: CharLocation{Line: 2, Column: 1}, End: CharLocation{Line: 2, Column: 13}},
						ValueLocation: TokenLocation{Start: CharLocation{Line: 2, Column: 1}, End: CharLocation{Line: 2, Column: 13}},
					},
					"database_settings": {
						Type: Object,
						Value: map[string]*Node{
							"host": {
								Type:          String,
								Value:         "192.168.1.100",
								NameLocation:  TokenLocation{Start: CharLocation{Line: 7, Column: 0}, End: CharLocation{Line: 7, Column: 4}},
								ValueLocation: TokenLocation{Start: CharLocation{Line: 7, Column: 7}, End: CharLocation{Line: 7, Column: 22}},
							},
							"port_numbers": {
								Type: Array,
								Value: []*Node{
									{
										Type:          Int,
										Value:         8005,
										ValueLocation: TokenLocation{Start: CharLocation{Line: 8, Column: 17}, End: CharLocation{Line: 8, Column: 21}},
									},
									{
										Type:          Int,
										Value:         8006,
										ValueLocation: TokenLocation{Start: CharLocation{Line: 8, Column: 23}, End: CharLocation{Line: 8, Column: 27}},
									},
									{
										Type:          Int,
										Value:         8007,
										ValueLocation: TokenLocation{Start: CharLocation{Line: 8, Column: 29}, End: CharLocation{Line: 8, Column: 33}},
									},
								},
								NameLocation:  TokenLocation{Start: CharLocation{Line: 8, Column: 0}, End: CharLocation{Line: 8, Column: 12}},
								ValueLocation: TokenLocation{Start: CharLocation{Line: 8, Column: 15}, End: CharLocation{Line: 8, Column: 35}},
							},
							"max_connections": {
								Type:          Int,
								Value:         6000,
								NameLocation:  TokenLocation{Start: CharLocation{Line: 9, Column: 0}, End: CharLocation{Line: 9, Column: 15}},
								ValueLocation: TokenLocation{Start: CharLocation{Line: 9, Column: 18}, End: CharLocation{Line: 9, Column: 22}},
							},
							"is_active": {
								Type:          Bool,
								Value:         false,
								NameLocation:  TokenLocation{Start: CharLocation{Line: 10, Column: 0}, End: CharLocation{Line: 10, Column: 9}},
								ValueLocation: TokenLocation{Start: CharLocation{Line: 10, Column: 12}, End: CharLocation{Line: 10, Column: 17}},
							},
						},
						NameLocation:  TokenLocation{Start: CharLocation{Line: 6, Column: 1}, End: CharLocation{Line: 6, Column: 18}},
						ValueLocation: TokenLocation{Start: CharLocation{Line: 6, Column: 1}, End: CharLocation{Line: 6, Column: 18}},
					},
					"web_servers": {
						Type: Object,
						Value: map[string]*Node{
							"primary": {
								Type: Object,
								Value: map[string]*Node{
									"address": {
										Type:          String,
										Value:         "10.0.0.3",
										NameLocation:  TokenLocation{Start: CharLocation{Line: 16, Column: 2}, End: CharLocation{Line: 16, Column: 9}},
										ValueLocation: TokenLocation{Start: CharLocation{Line: 16, Column: 12}, End: CharLocation{Line: 16, Column: 22}},
									},
									"datacenter": {
										Type:          String,
										Value:         "dc01",
										NameLocation:  TokenLocation{Start: CharLocation{Line: 17, Column: 2}, End: CharLocation{Line: 17, Column: 12}},
										ValueLocation: TokenLocation{Start: CharLocation{Line: 17, Column: 15}, End: CharLocation{Line: 17, Column: 21}},
									},
								},
								NameLocation:  TokenLocation{Start: CharLocation{Line: 15, Column: 3}, End: CharLocation{Line: 15, Column: 24}},
								ValueLocation: TokenLocation{Start: CharLocation{Line: 15, Column: 3}, End: CharLocation{Line: 15, Column: 24}},
							},
							"secondary": {
								Type: Object,
								Value: map[string]*Node{
									"address": {
										Type:          String,
										Value:         "10.0.0.4",
										NameLocation:  TokenLocation{Start: CharLocation{Line: 20, Column: 2}, End: CharLocation{Line: 20, Column: 9}},
										ValueLocation: TokenLocation{Start: CharLocation{Line: 20, Column: 12}, End: CharLocation{Line: 20, Column: 22}},
									},
									"datacenter": {
										Type:          String,
										Value:         "dc02",
										NameLocation:  TokenLocation{Start: CharLocation{Line: 21, Column: 2}, End: CharLocation{Line: 21, Column: 12}},
										ValueLocation: TokenLocation{Start: CharLocation{Line: 21, Column: 15}, End: CharLocation{Line: 21, Column: 21}},
									},
								},
								NameLocation:  TokenLocation{Start: CharLocation{Line: 19, Column: 3}, End: CharLocation{Line: 19, Column: 24}},
								ValueLocation: TokenLocation{Start: CharLocation{Line: 19, Column: 3}, End: CharLocation{Line: 19, Column: 24}},
							},
						},
						NameLocation:  TokenLocation{Start: CharLocation{Line: 12, Column: 1}, End: CharLocation{Line: 12, Column: 12}},
						ValueLocation: TokenLocation{Start: CharLocation{Line: 12, Column: 1}, End: CharLocation{Line: 12, Column: 12}},
					},
					"client_info": {
						Type: Object,
						Value: map[string]*Node{
							"metadata": {
								Type: Array,
								Value: []*Node{
									{
										Type: Array,
										Value: []*Node{
											{
												Type:          String,
												Value:         "epsilon",
												ValueLocation: TokenLocation{Start: CharLocation{Line: 24, Column: 14}, End: CharLocation{Line: 24, Column: 23}},
											},
											{
												Type:          String,
												Value:         "zeta",
												ValueLocation: TokenLocation{Start: CharLocation{Line: 24, Column: 25}, End: CharLocation{Line: 24, Column: 31}},
											},
										},
										ValueLocation: TokenLocation{Start: CharLocation{Line: 24, Column: 13}, End: CharLocation{Line: 24, Column: 32}},
									},
									{
										Type: Array,
										Value: []*Node{
											{
												Type:          Int,
												Value:         3,
												ValueLocation: TokenLocation{Start: CharLocation{Line: 24, Column: 35}, End: CharLocation{Line: 24, Column: 36}},
											},
											{
												Type:          Int,
												Value:         4,
												ValueLocation: TokenLocation{Start: CharLocation{Line: 24, Column: 38}, End: CharLocation{Line: 24, Column: 39}},
											},
										},
										ValueLocation: TokenLocation{Start: CharLocation{Line: 24, Column: 34}, End: CharLocation{Line: 24, Column: 40}},
									},
								},
								NameLocation:  TokenLocation{Start: CharLocation{Line: 24, Column: 0}, End: CharLocation{Line: 24, Column: 8}},
								ValueLocation: TokenLocation{Start: CharLocation{Line: 24, Column: 11}, End: CharLocation{Line: 24, Column: 42}},
							},
							"access_nodes": {
								Type: Array,
								Value: []*Node{
									{
										Type:          String,
										Value:         "beta",
										ValueLocation: TokenLocation{Start: CharLocation{Line: 28, Column: 2}, End: CharLocation{Line: 28, Column: 8}},
									},
									{
										Type:          String,
										Value:         "theta",
										ValueLocation: TokenLocation{Start: CharLocation{Line: 29, Column: 2}, End: CharLocation{Line: 29, Column: 9}},
									},
								},
								NameLocation:  TokenLocation{Start: CharLocation{Line: 27, Column: 0}, End: CharLocation{Line: 27, Column: 12}},
								ValueLocation: TokenLocation{Start: CharLocation{Line: 27, Column: 15}, End: CharLocation{Line: 30, Column: 1}},
							},
						},
						NameLocation:  TokenLocation{Start: CharLocation{Line: 23, Column: 1}, End: CharLocation{Line: 23, Column: 12}},
						ValueLocation: TokenLocation{Start: CharLocation{Line: 23, Column: 1}, End: CharLocation{Line: 23, Column: 12}},
					},
				},
				NameLocation:  TokenLocation{Start: CharLocation{Line: 0, Column: 0}, End: CharLocation{Line: 0, Column: 0}},
				ValueLocation: TokenLocation{Start: CharLocation{Line: 0, Column: 0}, End: CharLocation{Line: 0, Column: 0}},
			},
			expectedErrs: []CMParserError{},
		},
	}

	// Run tests
	for _, test := range testCases {
		parser := &tomlParser{}
		result, errs := parser.Parse(test.input)

		if len(errs) > 0 {
			t.Errorf("Unexpected errors: %#v", errs)
		} else if !reflect.DeepEqual(test.expected, result) {
			t.Errorf("Expected %#+v, got %#+v", test.expected, result)
		}
	}
}

func TestShortSamples_tomlParser(t *testing.T) {
	// Input
	var shortTomlConfig0 = []byte(`
	"" = "blank"     
	'' = 'blank'
	`)

	var shortTomlConfig1 = []byte(`
	bare_key = "value"
	bare-key = "value"
	1234 = "value"
	"127.0.0.1" = "value"
	"character encoding" = "value"
	"ʎǝʞ" = "value"
	'key2' = "value"
	'quoted "value"' = "value"
	`)

	var shortTomlConfig2 = []byte(`
	name = "Orange"
	physical.color = "orange"
	physical.shape = "round"
	site."google.com" = true
	`)

	var shortTomlConfig3 = []byte(`
	fruit.name = "banana"
	fruit. color = "yellow"
	fruit . flavor = "banana"
	`)

	var shortTomlConfig4 = []byte(`
	apple.type = "fruit"
	orange.type = "fruit"

	apple.skin = "thin"
	orange.skin = "thick"
	`)

	var shortTomlConfig5 = []byte(`
	3.14159 = "pi"
	`)

	var shortTomlConfig6 = []byte(`
	str = "I'm a string. \"You can quote me\". Name\tJos\u00E9\nLocation\tSF."
	str2 = """
	Roses are red
	Violets are blue"""
	str3 = """The quick brown \
	fox jumps over \
	the lazy dog."""
	str4 = '\\ServerX\admin$\system32\'
	str5 = '''
	The first newline is
	trimmed in raw strings.
	   All other whitespace
	   is preserved.
	'''
	`)

	var shortTomlConfig7 = []byte(`
	flt2 = 3.1415
	flt3 = -0.01
	flt4 = 5e+22
	`)

	var shortTomlConfig8 = []byte(`
	contributors = [
		"Foo Bar <foo@example.com>",
		{ name = "Baz Qux", email = "bazqux@example.com", url = "https://example.com/bazqux" }
	]
	`)

	var shortTomlConfig9 = []byte(`
	[table-1]
	key1 = "some string"
	key2 = 123

	[table-2]
	key1 = "another string"
	key2 = 456
	`)

	var shortTomlConfig10 = []byte(`
	[dog."tater.man"]
	type.name = "pug"
	`)

	var shortTomlConfig11 = []byte(`
	[a.b.c]            
	[ d.e.f ]          
	[ g .  h  . i ]    
	[ j . "ʞ" . 'l' ]  
	`)

	var shortTomlConfig12 = []byte(`
	[x.y.z.w]
	some = "thing"
	[x]
	other = "thing"
	`)

	var shortTomlConfig13 = []byte(`
	[fruit.apple]
	[animal]
	[fruit.orange]
	`)

	var shortTomlConfig14 = []byte(`
	name = { first = "Tom", last = "Preston-Werner" }
	point = { x = 1, y = 2 }
	animal = { type.name = "pug" }
	`)

	var shortTomlConfig15 = []byte(`
	[[products]]
	name = "Hammer"
	sku = 738594937

	[[products]]  # empty table within the array

	[[products]]
	name = "Nail"
	sku = 284758393

	color = "gray"
	`)

	var shortTomlConfig16 = []byte(`
	[[fruits]]
	name = "apple"

	[fruits.physical]  # subtable
	color = "red"
	shape = "round"

	[[fruits.varieties]]  # nested array of tables
	name = "red delicious"

	[[fruits.varieties]]
	name = "granny smith"


	[[fruits]]
	name = "banana"

	[[fruits.varieties]]
	name = "plantain"
	`)

	var shortTomlConfig17 = []byte(`
	points = [ { x = 1, y = 2, z = 3 },
           { x = 7, y = 8, z = 9 },
           { x = 2, y = 4, z = 8 } ]
	`)

	testCases := []jsonParserTestCase{
		{
			input:        shortTomlConfig0,
			expected:     nil,
			expectedErrs: nil,
		},
		{
			input:        shortTomlConfig1,
			expected:     nil,
			expectedErrs: nil,
		},
		{
			input:        shortTomlConfig2,
			expected:     nil,
			expectedErrs: nil,
		},
		{
			input:        shortTomlConfig3,
			expected:     nil,
			expectedErrs: nil,
		},
	}

	// Run tests
	for _, test := range testCases {
		parser := &tomlParser{}
		result, errs := parser.Parse(test.input)

		if len(errs) > 0 {
			t.Errorf("Unexpected errors: %#v", errs)
		} else if !reflect.DeepEqual(test.expected, result) {
			t.Errorf("Expected %#+v, got %#+v", test.expected, result)
		}
	}
}

func TestHighLevelErrorConditions_tomlParser(t *testing.T) {
	// Input
	var hlErrTomlConfig0 = []byte(`
	name = "Tom"
	name = "Pradyun" 
	`)

	var hlErrTomlConfig1 = []byte(`
	spelling = "favorite"
	"spelling" = "favourite"
	`)

	var hlErrTomlConfig2 = []byte(`
	fruit.apple = 1
	fruit.apple.smooth = true
	`)

	var hlErrTomlConfig3 = []byte(`
	apos15 = "Here are fifteen apostrophes: '''''''''''''''"
	`)

	var hlErrTomlConfig4 = []byte(`
	invalid_float_1 = .7
	invalid_float_2 = 7.
	invalid_float_3 = 3.e+20
	`)

	var hlErrTomlConfig5 = []byte(`
	[fruit]
	apple = "red"

	[fruit]
	orange = "orange"
	`)

	var hlErrTomlConfig6 = []byte(`
	[fruit]
	apple = "red"

	[fruit.apple]
	texture = "smooth"
	`)

	var hlErrTomlConfig7 = []byte(`
	[product]
	type = { name = "Nail" }
	type.edible = false
	`)

	var hlErrTomlConfig8 = []byte(`
	[product]
	type.name = "Nail"
	type = { edible = false }
	`)

	var hlErrTomlConfig9 = []byte(`
	[fruit.physical]
	color = "red"
	shape = "round"

	[[fruit]]
	name = "apple"
	`)

	var hlErrTomlConfig10 = []byte(`
	fruits = []

	[[fruits]]
	`)

	var hlErrTomlConfig11 = []byte(`
	[[fruits]]
	name = "apple"

	[[fruits.varieties]]
	name = "red delicious"

	[fruits.varieties]
	name = "granny smith"

	[fruits.physical]
	color = "red"
	shape = "round"

	[[fruits.physical]]
	color = "green"
	`)

	testCases := []jsonParserTestCase{
		{
			input:    hlErrTomlConfig0,
			expected: nil,
			expectedErrs: []CMParserError{
				{
					Message: "mismatched input '<EOF>' expecting {'[', '{', BOOLEAN, BASIC_STRING, ML_BASIC_STRING, LITERAL_STRING, ML_LITERAL_STRING, FLOAT, INF, NAN, DEC_INT, HEX_INT, OCT_INT, BIN_INT, OFFSET_DATE_TIME, LOCAL_DATE_TIME, LOCAL_DATE, LOCAL_TIME}",
					Location: TokenLocation{
						Start: CharLocation{Line: 0, Column: 6},
						End:   CharLocation{Line: 0, Column: 7},
					},
				},
			},
		},
		{
			input:    hlErrTomlConfig1,
			expected: nil,
			expectedErrs: []CMParserError{
				{
					Message: "mismatched input 'last' expecting {<EOF>, NL}",
					Location: TokenLocation{
						Start: CharLocation{Line: 0, Column: 14},
						End:   CharLocation{Line: 0, Column: 15},
					},
				},
			},
		},
		{
			input:    hlErrTomlConfig2,
			expected: nil,
			expectedErrs: []CMParserError{
				{
					Message: "no viable alternative at input '='",
					Location: TokenLocation{
						Start: CharLocation{Line: 0, Column: 0},
						End:   CharLocation{Line: 0, Column: 1},
					},
				},
			},
		},
	}

	// Run tests
	for _, test := range testCases {
		parser := &tomlParser{}
		res, errs := parser.Parse(test.input)

		if len(errs) == 0 {
			t.Errorf("Expected errors, got none")
		} else if !reflect.DeepEqual(test.expectedErrs, errs) {
			t.Errorf("Expected %v, got %v, %v", test.expectedErrs, errs, res)
		}
	}
}

func TestErrorConditions_tomlParser(t *testing.T) {
	// Input
	var errTomlConfig0 = []byte(`key = `)
	var errTomlConfig1 = []byte(`first = "Tom" last = "Preston-Werner"`)
	var errTomlConfig2 = []byte(`= "no key name"`)

	testCases := []jsonParserTestCase{
		{
			input:    errTomlConfig0,
			expected: nil,
			expectedErrs: []CMParserError{
				{
					Message: "mismatched input '<EOF>' expecting {'[', '{', BOOLEAN, BASIC_STRING, ML_BASIC_STRING, LITERAL_STRING, ML_LITERAL_STRING, FLOAT, INF, NAN, DEC_INT, HEX_INT, OCT_INT, BIN_INT, OFFSET_DATE_TIME, LOCAL_DATE_TIME, LOCAL_DATE, LOCAL_TIME}",
					Location: TokenLocation{
						Start: CharLocation{Line: 0, Column: 6},
						End:   CharLocation{Line: 0, Column: 7},
					},
				},
			},
		},
		{
			input:    errTomlConfig1,
			expected: nil,
			expectedErrs: []CMParserError{
				{
					Message: "mismatched input 'last' expecting {<EOF>, NL}",
					Location: TokenLocation{
						Start: CharLocation{Line: 0, Column: 14},
						End:   CharLocation{Line: 0, Column: 15},
					},
				},
			},
		},
		{
			input:    errTomlConfig2,
			expected: nil,
			expectedErrs: []CMParserError{
				{
					Message: "no viable alternative at input '='",
					Location: TokenLocation{
						Start: CharLocation{Line: 0, Column: 0},
						End:   CharLocation{Line: 0, Column: 1},
					},
				},
			},
		},
	}

	// Run tests
	for _, test := range testCases {
		parser := &tomlParser{}
		res, errs := parser.Parse(test.input)

		if len(errs) == 0 {
			t.Errorf("Expected errors, got none")
		} else if !reflect.DeepEqual(test.expectedErrs, errs) {
			t.Errorf("Expected %v, got %v, %v", test.expectedErrs, errs, res)
		}
	}
}
