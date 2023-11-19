package check

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ConfigMate/configmate/analyzer/types"
	"github.com/ConfigMate/configmate/parsers"
)

type checkEvaluatorTestStructure struct {
	// Input values
	primaryField     string                 // Primary field name
	fields           map[string]types.IType // Fields
	optMissingFields map[string]bool        // Optional missing fields
	checks           []string               // Checks

	// Expected values
	expectedRes     types.IType // Expected result
	expectedSkipped bool        // Expected skipped value
	expectedErr     error       // Expected error
}

// TestEvaluateBasicFunctionality tests the basic functionality of the
// check evaluator. It tests checks like:
//   - eq(5)
//   - range(0, 10)
//   - lt(5)
//   - lt(config.other)
//   - range(config.other1, config.other2)
func TestEvaluateBasicFunctionality(t *testing.T) {
	// Test cases
	tests := []checkEvaluatorTestStructure{
		// Test 1:
		// This is a simple case where we are checking
		// a field of type bool with a check verifing
		// is it true. No other fields involved.
		// The check should pass.
		func() checkEvaluatorTestStructure {
			primaryField := "config.primary" // Create primary field

			fields := make(map[string]types.IType)     // Create fields
			pFValue, _ := types.MakeType("bool", true) // Create primary field
			fields[primaryField] = pFValue             // Add primary field to fields

			optMissingFields := make(map[string]bool) // Create optional missing fields
			checks := []string{"eq(true)"}            // Create checks

			expectedRes, _ := types.MakeType("bool", true) // Create expected result
			expectedSkipped := false                       // Create expected skipped value
			expectedErr := error(nil)                      // Create expected error

			return checkEvaluatorTestStructure{
				primaryField:     primaryField,
				fields:           fields,
				optMissingFields: optMissingFields,
				checks:           checks,
				expectedRes:      expectedRes,
				expectedSkipped:  expectedSkipped,
				expectedErr:      expectedErr,
			}
		}(),
		// Test 2:
		// This test is the same as test 1, except
		// the check should fail.
		func() checkEvaluatorTestStructure {
			primaryField := "primary" // Create primary field

			fields := make(map[string]types.IType)      // Create fields
			pFValue, _ := types.MakeType("bool", false) // Create primary field
			fields[primaryField] = pFValue              // Add primary field to fields

			optMissingFields := make(map[string]bool) // Create optional missing fields
			checks := []string{"eq(true)"}            // Create checks

			expectedRes, _ := types.MakeType("bool", false)            // Create expected result
			expectedSkipped := false                                   // Create expected skipped value
			expectedErr := fmt.Errorf("bool.eq failed: true != false") // Create expected error

			return checkEvaluatorTestStructure{
				primaryField:     primaryField,
				fields:           fields,
				optMissingFields: optMissingFields,
				checks:           checks,
				expectedRes:      expectedRes,
				expectedSkipped:  expectedSkipped,
				expectedErr:      expectedErr,
			}
		}(),
		// Test 3:
		// This test verifies the functionality of invoking subsequet checks.
		// All checks pass individually. We start with a field of type int,
		// check that it is within a range (this check passes), and then check
		// we check that that result is true (redundant)
		func() checkEvaluatorTestStructure {
			primaryField := "primary.field" // Create primary field

			fields := make(map[string]types.IType) // Create fields
			pFValue, _ := types.MakeType("int", 5) // Create primary field
			fields[primaryField] = pFValue         // Add primary field to fields

			optMissingFields := make(map[string]bool)   // Create optional missing fields
			checks := []string{"range(0, 10).eq(true)"} // Create checks

			expectedRes, _ := types.MakeType("bool", true) // Create expected result
			expectedSkipped := false                       // Create expected skipped value
			expectedErr := error(nil)                      // Create expected error

			return checkEvaluatorTestStructure{
				primaryField:     primaryField,
				fields:           fields,
				optMissingFields: optMissingFields,
				checks:           checks,
				expectedRes:      expectedRes,
				expectedSkipped:  expectedSkipped,
				expectedErr:      expectedErr,
			}
		}(),
		// Test 4:
		// This test also verifies the functionality of invoking subsequet checks.
		// The overall check passes but an intermidiate check returns false.
		// The overall results should be nothing but a success.
		func() checkEvaluatorTestStructure {
			primaryField := "primary.field" // Create primary field

			fields := make(map[string]types.IType) // Create fields
			pFValue, _ := types.MakeType("int", 5) // Create primary field
			fields[primaryField] = pFValue         // Add primary field to fields

			optMissingFields := make(map[string]bool) // Create optional missing fields
			checks := []string{"!range(0, 3)"}        // Create checks

			expectedRes, _ := types.MakeType("bool", true) // Create expected result
			expectedSkipped := false                       // Create expected skipped value
			expectedErr := error(nil)                      // Create expected error

			return checkEvaluatorTestStructure{
				primaryField:     primaryField,
				fields:           fields,
				optMissingFields: optMissingFields,
				checks:           checks,
				expectedRes:      expectedRes,
				expectedSkipped:  expectedSkipped,
				expectedErr:      expectedErr,
			}
		}(),
		// Test 5:
		// This test also verifies the functionality of invoking subsequet checks.
		// In this case one of them returns an intermidiate type (not a bool), the
		// overall check fails.
		func() checkEvaluatorTestStructure {
			primaryField := "primary.field" // Create primary field

			fields := make(map[string]types.IType)      // Create fields
			pFValue, _ := types.MakeType("float", 3.14) // Create primary field
			fields[primaryField] = pFValue              // Add primary field to fields

			optMissingFields := make(map[string]bool) // Create optional missing fields
			checks := []string{"toInt().range(0, 2)"} // Create checks

			expectedRes, _ := types.MakeType("bool", false)                      // Create expected result
			expectedSkipped := false                                             // Create expected skipped value
			expectedErr := fmt.Errorf("int.range failed: 3 not in range [0, 2]") // Create expected error

			return checkEvaluatorTestStructure{
				primaryField:     primaryField,
				fields:           fields,
				optMissingFields: optMissingFields,
				checks:           checks,
				expectedRes:      expectedRes,
				expectedSkipped:  expectedSkipped,
				expectedErr:      expectedErr,
			}
		}(),
		// Test 6:
		// This test will verify the functionality of using other fields
		// as parameters for checks. In this case we have a primary field
		// of type int, and we check that is it less than another field.
		// The check should pass.
		func() checkEvaluatorTestStructure {
			primaryField := "primary.field" // Create primary field

			fields := make(map[string]types.IType) // Create fields
			pFValue, _ := types.MakeType("int", 5) // Create primary field
			fields[primaryField] = pFValue         // Add primary field to fields

			otherField, _ := types.MakeType("int", 10) // Create other field
			fields["config.other"] = otherField

			optMissingFields := make(map[string]bool) // Create optional missing fields
			checks := []string{"lt(config.other)"}    // Create checks

			expectedRes, _ := types.MakeType("bool", true) // Create expected result
			expectedSkipped := false                       // Create expected skipped value
			expectedErr := error(nil)                      // Create expected error

			return checkEvaluatorTestStructure{
				primaryField:     primaryField,
				fields:           fields,
				optMissingFields: optMissingFields,
				checks:           checks,
				expectedRes:      expectedRes,
				expectedSkipped:  expectedSkipped,
				expectedErr:      expectedErr,
			}
		}(),
		// Test 7:
		// This test will verify the functionality of using multiple other fields
		// as parameters for checks. In this case we have a primary field of
		// type int, and we check that it is between two other fields. The check
		// should pass.
		func() checkEvaluatorTestStructure {
			primaryField := "field" // Create primary field

			fields := make(map[string]types.IType) // Create fields
			pFValue, _ := types.MakeType("int", 5) // Create primary field
			fields[primaryField] = pFValue

			otherField1, _ := types.MakeType("int", 0) // Create other field 1
			otherField2, _ := types.MakeType("int", 10)
			fields["config.other1"] = otherField1
			fields["config.other2"] = otherField2

			optMissingFields := make(map[string]bool)                 // Create optional missing fields
			checks := []string{"range(config.other1, config.other2)"} // Create checks

			expectedRes, _ := types.MakeType("bool", true) // Create expected result
			expectedSkipped := false                       // Create expected skipped value
			expectedErr := error(nil)                      // Create expected error

			return checkEvaluatorTestStructure{
				primaryField:     primaryField,
				fields:           fields,
				optMissingFields: optMissingFields,
				checks:           checks,
				expectedRes:      expectedRes,
				expectedSkipped:  expectedSkipped,
				expectedErr:      expectedErr,
			}
		}(),
		// Test 8:
		// This test will verify the functionality of using other fields with their own
		// checks as parameters for checks. In this case we have a primary field of type
		// int, and we check that is it less than another field of type float. To do this
		// we have to convert the float to int. The check should pass.
		func() checkEvaluatorTestStructure {
			primaryField := "primary.field" // Create primary field

			fields := make(map[string]types.IType) // Create fields
			pFValue, _ := types.MakeType("int", 5) // Create primary field
			fields[primaryField] = pFValue         // Add primary field to fields

			otherField, _ := types.MakeType("float", 10.0)
			fields["config.other"] = otherField

			optMissingFields := make(map[string]bool)      // Create optional missing fields
			checks := []string{"lt(config.other.toInt())"} // Create checks

			expectedRes, _ := types.MakeType("bool", true) // Create expected result
			expectedSkipped := false                       // Create expected skipped value
			expectedErr := error(nil)                      // Create expected error

			return checkEvaluatorTestStructure{
				primaryField:     primaryField,
				fields:           fields,
				optMissingFields: optMissingFields,
				checks:           checks,
				expectedRes:      expectedRes,
				expectedSkipped:  expectedSkipped,
				expectedErr:      expectedErr,
			}
		}(),
		// Test 9:
		// Similarly to 6 and 7, this test will verify the functionality of using multiple
		// other fields with their own checks as parameters for checks. In this case we have
		// a primary field of type int, and we check that it is between two other fields of
		// type float. To do this we have to convert the floats to ints. The check should fail.
		func() checkEvaluatorTestStructure {
			primaryField := "primary.field" // Create primary field

			fields := make(map[string]types.IType) // Create fields
			pFValue, _ := types.MakeType("int", 5) // Create primary field
			fields[primaryField] = pFValue         // Add primary field to fields

			otherField1, _ := types.MakeType("float", 0.5)
			otherField2, _ := types.MakeType("float", 3.14)
			fields["config.other1"] = otherField1
			fields["config.something.other2"] = otherField2

			optMissingFields := make(map[string]bool)                                           // Create optional missing fields
			checks := []string{"range(config.other1.toInt(), config.something.other2.toInt())"} // Create checks

			expectedRes, _ := types.MakeType("bool", false)                      // Create expected result
			expectedSkipped := false                                             // Create expected skipped value
			expectedErr := fmt.Errorf("int.range failed: 5 not in range [0, 3]") // Create expected error

			return checkEvaluatorTestStructure{
				primaryField:     primaryField,
				fields:           fields,
				optMissingFields: optMissingFields,
				checks:           checks,
				expectedRes:      expectedRes,
				expectedSkipped:  expectedSkipped,
				expectedErr:      expectedErr,
			}
		}(),
	}

	for _, test := range tests {
		// Create evaluator
		evaluator := NewCheckEvaluator()

		// Evaluate checks
		for _, check := range test.checks {
			res, skipped, err := evaluator.Evaluate(check, test.primaryField, test.fields, test.optMissingFields)
			errMessage := ""
			if err != nil {
				errMessage = err.Error()
			}
			expectedErrMessage := ""
			if test.expectedErr != nil {
				expectedErrMessage = test.expectedErr.Error()
			}
			if !reflect.DeepEqual(res, test.expectedRes) || !reflect.DeepEqual(skipped, test.expectedSkipped) || errMessage != expectedErrMessage {
				t.Errorf("Evaluate(%v, %v, %v, %v) = %v, %v, %v, want %v, %v, %v", test.primaryField, test.fields, test.optMissingFields, check, res, skipped, errMessage, test.expectedRes, test.expectedSkipped, expectedErrMessage)
			}
		}
	}
}

// TestEvaluateOptionalMissingFields tests the functionality of the check
// evaluator when optional missing fields are involved. It tests checks like:
//   - eq(config.missingOptField)
//   - range(config.missingOptField1, 128)
func TestEvaluateOptionalMissingFields(t *testing.T) {
	// Test cases
	tests := []checkEvaluatorTestStructure{
		// Test 1:
		// This test verifies the functionality when optional missing field is involved.
		// In this case, the check should skip.
		func() checkEvaluatorTestStructure {
			primaryField := "primary.field" // Create primary field

			fields := make(map[string]types.IType) // Create fields
			pFValue, _ := types.MakeType("int", 5) // Create primary field value
			fields[primaryField] = pFValue         // Add primary field to fields

			optMissingFields := make(map[string]bool) // Create optional missing fields
			optMissingFields["config.missingOptField"] = true
			checks := []string{"eq(config.missingOptField)"} // Create checks

			expectedRes, _ := types.MakeType("bool", false)                                                                 // Create expected result
			expectedSkipped := true                                                                                         // Create expected skipped value
			expectedErr := fmt.Errorf("skipping check because referenced optional field config.missingOptField is missing") // Create expected error

			return checkEvaluatorTestStructure{
				primaryField:     primaryField,
				fields:           fields,
				optMissingFields: optMissingFields,
				checks:           checks,
				expectedRes:      expectedRes,
				expectedSkipped:  expectedSkipped,
				expectedErr:      expectedErr,
			}
		}(),
		// Test 2:
		// This test verifies the functionality when multiple optional missing fields are involved.
		// In this case, the check should skip as well.
		func() checkEvaluatorTestStructure {
			primaryField := "primary.field" // Create primary field

			fields := make(map[string]types.IType) // Create fields
			pFValue, _ := types.MakeType("int", 5) // Create primary field value
			fields[primaryField] = pFValue         // Add primary field to fields

			optMissingFields := make(map[string]bool) // Create optional missing fields
			optMissingFields["config.missingOptField1"] = true
			checks := []string{"range(config.missingOptField1, 128)"} // Create checks

			expectedRes, _ := types.MakeType("bool", false)                                                                  // Create expected result
			expectedSkipped := true                                                                                          // Create expected skipped value
			expectedErr := fmt.Errorf("skipping check because referenced optional field config.missingOptField1 is missing") // Create expected error

			return checkEvaluatorTestStructure{
				primaryField:     primaryField,
				fields:           fields,
				optMissingFields: optMissingFields,
				checks:           checks,
				expectedRes:      expectedRes,
				expectedSkipped:  expectedSkipped,
				expectedErr:      expectedErr,
			}
		}(),
	}

	for _, test := range tests {
		// Create evaluator
		evaluator := NewCheckEvaluator()

		// Evaluate checks
		for _, check := range test.checks {
			res, skipped, err := evaluator.Evaluate(check, test.primaryField, test.fields, test.optMissingFields)
			errMessage := ""
			if err != nil {
				errMessage = err.Error()
			}
			expectedErrMessage := ""
			if test.expectedErr != nil {
				expectedErrMessage = test.expectedErr.Error()
			}
			if !reflect.DeepEqual(res, test.expectedRes) || !reflect.DeepEqual(skipped, test.expectedSkipped) || errMessage != expectedErrMessage {
				t.Errorf("Evaluate(%v, %v, %v, %v) = %v, %v, %v, want %v, %v, %v", test.primaryField, test.fields, test.optMissingFields, check, res, skipped, errMessage, test.expectedRes, test.expectedSkipped, expectedErrMessage)
			}
		}
	}
}

// TestEvaluateLogicalExpressions tests the functionality of the check
// evaluator when logical expressions are involved. It tests checks like:
//   - eq(5) && eq(10)
//   - eq(5) || eq(10)
//   - (range(0, 5) || range(10, 15)) && range(0, 25)
func TestEvaluateLogicalExpressions(t *testing.T) {
	// Test cases
	tests := []checkEvaluatorTestStructure{
		// Test 1: OR expression, left side true
		func() checkEvaluatorTestStructure {
			primaryField := "primary.field"

			pFValue, _ := types.MakeType("int", 10)
			fields := map[string]types.IType{primaryField: pFValue}

			checks := []string{"range(0, 15) || range(20, 25)"}

			expectedRes, _ := types.MakeType("bool", true)

			return checkEvaluatorTestStructure{
				primaryField:    primaryField,
				fields:          fields,
				checks:          checks,
				expectedRes:     expectedRes,
				expectedSkipped: false,
				expectedErr:     nil,
			}
		}(),
		// Test 2: OR expression, both sides true
		func() checkEvaluatorTestStructure {
			primaryField := "primary.field"

			pFValue, _ := types.MakeType("int", 10)
			fields := map[string]types.IType{primaryField: pFValue}

			checks := []string{"range(0, 15) || range(5, 25)"}

			expectedRes, _ := types.MakeType("bool", true)

			return checkEvaluatorTestStructure{
				primaryField:    primaryField,
				fields:          fields,
				checks:          checks,
				expectedRes:     expectedRes,
				expectedSkipped: false,
				expectedErr:     nil,
			}
		}(),
		// Test 3: OR expression, both sides false
		func() checkEvaluatorTestStructure {
			primaryField := "primary.field"

			pFValue, _ := types.MakeType("int", 30)
			fields := map[string]types.IType{primaryField: pFValue}

			checks := []string{"range(0, 15) || range(20, 25)"}

			expectedRes, _ := types.MakeType("bool", false)

			return checkEvaluatorTestStructure{
				primaryField:    primaryField,
				fields:          fields,
				checks:          checks,
				expectedRes:     expectedRes,
				expectedSkipped: false,
				expectedErr:     fmt.Errorf("int.range failed: 30 not in range [0, 15]; int.range failed: 30 not in range [20, 25]"),
			}
		}(),
		// Test 4: AND expression, both sides true
		func() checkEvaluatorTestStructure {
			primaryField := "primary.field"

			pFValue, _ := types.MakeType("int", 10)
			fields := map[string]types.IType{primaryField: pFValue}

			checks := []string{"range(0, 15) && range(5, 25)"}

			expectedRes, _ := types.MakeType("bool", true)

			return checkEvaluatorTestStructure{
				primaryField:    primaryField,
				fields:          fields,
				checks:          checks,
				expectedRes:     expectedRes,
				expectedSkipped: false,
				expectedErr:     nil,
			}
		}(),
		// Test 5: AND expression, left side false
		func() checkEvaluatorTestStructure {
			primaryField := "primary.field"

			pFValue, _ := types.MakeType("int", 30)
			fields := map[string]types.IType{primaryField: pFValue}

			checks := []string{"range(0, 15) && range(20, 40)"}

			expectedRes, _ := types.MakeType("bool", false)

			return checkEvaluatorTestStructure{
				primaryField:    primaryField,
				fields:          fields,
				checks:          checks,
				expectedRes:     expectedRes,
				expectedSkipped: false,
				expectedErr:     fmt.Errorf("int.range failed: 30 not in range [0, 15]"),
			}
		}(),
		// Test 6: Nested expression, mixed operators, all true
		func() checkEvaluatorTestStructure {
			primaryField := "primary.field"

			pFValue, _ := types.MakeType("int", 10)
			fields := map[string]types.IType{primaryField: pFValue}

			checks := []string{"range(0, 15) && (range(5, 25) || range(30, 40))"}

			expectedRes, _ := types.MakeType("bool", true)

			return checkEvaluatorTestStructure{
				primaryField:    primaryField,
				fields:          fields,
				checks:          checks,
				expectedRes:     expectedRes,
				expectedSkipped: false,
				expectedErr:     nil,
			}
		}(),
		// Test 7: Nested expression, mixed operators, one false
		func() checkEvaluatorTestStructure {
			primaryField := "primary.field"

			pFValue, _ := types.MakeType("int", 30)
			fields := map[string]types.IType{primaryField: pFValue}

			checks := []string{"range(0, 15) && (range(20, 40) || range(45, 50))"}

			expectedRes, _ := types.MakeType("bool", false)

			return checkEvaluatorTestStructure{
				primaryField:    primaryField,
				fields:          fields,
				checks:          checks,
				expectedRes:     expectedRes,
				expectedSkipped: false,
				expectedErr:     fmt.Errorf("int.range failed: 30 not in range [0, 15]"),
			}
		}(),
		// Test 8: AND expression with eq function, both sides true
		func() checkEvaluatorTestStructure {
			primaryField := "primary.field"

			pFValue, _ := types.MakeType("int", 10)
			oFValue, _ := types.MakeType("int", 10)
			fields := map[string]types.IType{primaryField: pFValue, "other.field": oFValue}

			checks := []string{"eq(10) && other.field.eq(10)"}

			expectedRes, _ := types.MakeType("bool", true)

			return checkEvaluatorTestStructure{
				primaryField:    primaryField,
				fields:          fields,
				checks:          checks,
				expectedRes:     expectedRes,
				expectedSkipped: false,
				expectedErr:     nil,
			}
		}(),
		// Test 9: OR expression with gt function, left side true
		func() checkEvaluatorTestStructure {
			primaryField := "primary.field"

			pFValue, _ := types.MakeType("float", 15.5)
			fields := map[string]types.IType{primaryField: pFValue}

			checks := []string{"gt(10.0) || lt(10.0)"}

			expectedRes, _ := types.MakeType("bool", true)

			return checkEvaluatorTestStructure{
				primaryField:    primaryField,
				fields:          fields,
				checks:          checks,
				expectedRes:     expectedRes,
				expectedSkipped: false,
				expectedErr:     nil,
			}
		}(),
		// Test 10: Nested expression, mixed operators, with different functions
		func() checkEvaluatorTestStructure {
			primaryField := "primary.field"

			pFValue, _ := types.MakeType("string", "hello world")
			fields := map[string]types.IType{primaryField: pFValue}

			checks := []string{"regex(\"^hello\") && (eq(\"hello world\") || regex(\"world$\"))"}

			expectedRes, _ := types.MakeType("bool", true)

			return checkEvaluatorTestStructure{
				primaryField:    primaryField,
				fields:          fields,
				checks:          checks,
				expectedRes:     expectedRes,
				expectedSkipped: false,
				expectedErr:     nil,
			}
		}(),
		// Test 11: AND expression with toInt function, both sides true
		func() checkEvaluatorTestStructure {
			primaryField := "primary.field"

			pFValue, _ := types.MakeType("float", 10.0)
			oFValue, _ := types.MakeType("int", 10)
			fields := map[string]types.IType{primaryField: pFValue, "other.field": oFValue}

			checks := []string{"toInt().eq(10) && other.field.eq(10)"}

			expectedRes, _ := types.MakeType("bool", true)

			return checkEvaluatorTestStructure{
				primaryField:    primaryField,
				fields:          fields,
				checks:          checks,
				expectedRes:     expectedRes,
				expectedSkipped: false,
				expectedErr:     nil,
			}
		}(),
	}

	for _, test := range tests {
		// Create evaluator
		evaluator := NewCheckEvaluator()

		// Evaluate check
		for _, check := range test.checks {
			res, skipped, err := evaluator.Evaluate(check, test.primaryField, test.fields, test.optMissingFields)
			errMessage := ""
			if err != nil {
				errMessage = err.Error()
			}
			expectedErrMessage := ""
			if test.expectedErr != nil {
				expectedErrMessage = test.expectedErr.Error()
			}
			if !reflect.DeepEqual(res, test.expectedRes) || !reflect.DeepEqual(skipped, test.expectedSkipped) || errMessage != expectedErrMessage {
				t.Errorf("Evaluate(%v, %v, %v, %v) = %v, %v, %v, want %v, %v, %v", test.primaryField, test.fields, test.optMissingFields, check, res, skipped, errMessage, test.expectedRes, test.expectedSkipped, expectedErrMessage)
			}
		}
	}
}

// TestEvaluateControlExpressions tests the functionality of the check
// evaluator when control expressions are involved. It tests checks like:
//   - if(eq(true)){ eq(false) }
//   - if(eq(false)){ eq(true) } elseif(eq(true)){ eq(true) } else{ eq(false) }
//   - foreach(li : this){ li.gt(0) }
func TestEvaluateControlExpressions(t *testing.T) {
	// Test cases
	tests := []checkEvaluatorTestStructure{
		// Test 1: IF statement, false condition
		func() checkEvaluatorTestStructure {
			primaryField := "primary.field"

			pFValue, _ := types.MakeType("bool", false)
			fields := map[string]types.IType{primaryField: pFValue}

			checks := []string{"if(eq(true)){ eq(false) }"}

			expectedRes, _ := types.MakeType("bool", true)

			return checkEvaluatorTestStructure{
				primaryField:    primaryField,
				fields:          fields,
				checks:          checks,
				expectedRes:     expectedRes,
				expectedSkipped: false,
				expectedErr:     nil,
			}
		}(),
		// Test 2: IF-ELSEIF-ELSE statement, elseif condition true
		func() checkEvaluatorTestStructure {
			primaryField := "primary.field"

			pFValue, _ := types.MakeType("int", 10)
			fields := map[string]types.IType{primaryField: pFValue}

			checks := []string{"if(gt(20)){ eq(5) } elseif(lt(5)){ eq(2) } else{ eq(10) }"}

			expectedRes, _ := types.MakeType("bool", true)

			return checkEvaluatorTestStructure{
				primaryField:    primaryField,
				fields:          fields,
				checks:          checks,
				expectedRes:     expectedRes,
				expectedSkipped: false,
				expectedErr:     nil,
			}
		}(),
		// Test 3: FOREACH statement, simple list iteration
		func() checkEvaluatorTestStructure {
			primaryField := "primary.field"

			listItems := []*parsers.Node{
				{Value: 5},
				{Value: 10},
				{Value: 15},
			}
			pFValue, _ := types.MakeType("list:int", listItems)
			fields := map[string]types.IType{primaryField: pFValue}

			checks := []string{"foreach(li : this){ li.gt(0) }"}

			expectedRes, _ := types.MakeType("bool", true)

			return checkEvaluatorTestStructure{
				primaryField:    primaryField,
				fields:          fields,
				checks:          checks,
				expectedRes:     expectedRes,
				expectedSkipped: false,
				expectedErr:     nil,
			}
		}(),
		// Test 4: FOREACH statement, nested list iteration
		func() checkEvaluatorTestStructure {
			primaryField := "primary.field"

			innerListItems1 := []*parsers.Node{
				{Value: 1},
				{Value: 2},
			}
			innerListItems2 := []*parsers.Node{
				{Value: 3},
				{Value: -5},
			}
			listItems := []*parsers.Node{
				{Value: innerListItems1},
				{Value: innerListItems2},
			}
			pFValue, _ := types.MakeType("list:list:int", listItems)
			fields := map[string]types.IType{primaryField: pFValue}

			checks := []string{"foreach(li : this){ foreach(innerLi : li){ innerLi.gt(0) } }"}

			expectedRes, _ := types.MakeType("bool", false)

			return checkEvaluatorTestStructure{
				primaryField:    primaryField,
				fields:          fields,
				checks:          checks,
				expectedRes:     expectedRes,
				expectedSkipped: false,
				expectedErr:     fmt.Errorf("foreach body failed: item 1: foreach body failed: item 1: int.gt failed: -5 <= 0"),
			}
		}(),
	}

	for _, test := range tests {
		// Create evaluator
		evaluator := NewCheckEvaluator()

		// Evaluate check
		for _, check := range test.checks {
			res, skipped, err := evaluator.Evaluate(check, test.primaryField, test.fields, test.optMissingFields)
			errMessage := ""
			if err != nil {
				errMessage = err.Error()
			}
			expectedErrMessage := ""
			if test.expectedErr != nil {
				expectedErrMessage = test.expectedErr.Error()
			}
			if !reflect.DeepEqual(res, test.expectedRes) || !reflect.DeepEqual(skipped, test.expectedSkipped) || errMessage != expectedErrMessage {
				t.Errorf("Evaluate(%v, %v, %v, %v) = %v, %v, %v, want %v, %v, %v", test.primaryField, test.fields, test.optMissingFields, check, res, skipped, errMessage, test.expectedRes, test.expectedSkipped, expectedErrMessage)
			}
		}
	}
}

// TestEvaluateErroneosExpressions tests the functionality of the check
// evaluator when erroneous expressions are involved. It tests checks like:
//   - eq(5, 10) -  This function only receives one parameter
//   - toInt()   -	This does not result in a boolean
//   - eq(non_existent_field)  -  References a field that does not exist
func TestEvaluateErroneousExpressions(t *testing.T) {
	// Test cases
	tests := []checkEvaluatorTestStructure{
		// Test 1: Function with too many parameters
		func() checkEvaluatorTestStructure {
			primaryField := "primary.field"

			pFValue, _ := types.MakeType("int", 10)
			fields := map[string]types.IType{primaryField: pFValue}

			checks := []string{"eq(5, 10)"}

			return checkEvaluatorTestStructure{
				primaryField:    primaryField,
				fields:          fields,
				checks:          checks,
				expectedRes:     nil,
				expectedSkipped: false,
				expectedErr:     fmt.Errorf("int.eq expects 1 argument"),
			}
		}(),
		// Test 2: Function that does not result in a boolean
		func() checkEvaluatorTestStructure {
			primaryField := "primary.field"

			pFValue, _ := types.MakeType("float", 10.0)
			fields := map[string]types.IType{primaryField: pFValue}

			checks := []string{"toInt()"}

			return checkEvaluatorTestStructure{
				primaryField:    primaryField,
				fields:          fields,
				checks:          checks,
				expectedRes:     nil,
				expectedSkipped: false,
				expectedErr:     fmt.Errorf("check must evaluate to a bool"),
			}
		}(),
		// Test 3: Function that references a field that does not exist
		func() checkEvaluatorTestStructure {
			primaryField := "primary.field"

			pFValue, _ := types.MakeType("float", 10.0)
			fields := map[string]types.IType{primaryField: pFValue}

			checks := []string{"eq(non_existent_field)"}

			return checkEvaluatorTestStructure{
				primaryField:    primaryField,
				fields:          fields,
				checks:          checks,
				expectedRes:     nil,
				expectedSkipped: false,
				expectedErr:     fmt.Errorf("field non_existent_field does not exist"),
			}
		}(),
		// Test 4: And expression where right side is not a boolean
		func() checkEvaluatorTestStructure {
			primaryField := "primary.field"

			pFValue, _ := types.MakeType("float", 10.0)
			fields := map[string]types.IType{primaryField: pFValue}

			checks := []string{"eq(10.0) && toInt()"}

			return checkEvaluatorTestStructure{
				primaryField:    primaryField,
				fields:          fields,
				checks:          checks,
				expectedRes:     nil,
				expectedSkipped: false,
				expectedErr:     fmt.Errorf("and expression right expression must be a bool"),
			}
		}(),
		// Test 5: Or expression where left side is not a boolean
		func() checkEvaluatorTestStructure {
			primaryField := "primary.field"

			pFValue, _ := types.MakeType("float", 10.0)
			fields := map[string]types.IType{primaryField: pFValue}

			checks := []string{"toInt() || eq(10.0)"}

			return checkEvaluatorTestStructure{
				primaryField:    primaryField,
				fields:          fields,
				checks:          checks,
				expectedRes:     nil,
				expectedSkipped: false,
				expectedErr:     fmt.Errorf("or expression left expression must be a bool"),
			}
		}(),
		// Test 6: Foreach where the list item alias conflicts with another field
		func() checkEvaluatorTestStructure {
			primaryField := "primary.field"

			pFValue, _ := types.MakeType("list:int", []*parsers.Node{{Value: 10}})
			fields := map[string]types.IType{primaryField: pFValue}

			// Make another field
			otherField, _ := types.MakeType("int", 5)
			fields["otherField"] = otherField

			checks := []string{"foreach(otherField : this){ eq(10) }"}

			return checkEvaluatorTestStructure{
				primaryField:    primaryField,
				fields:          fields,
				checks:          checks,
				expectedRes:     nil,
				expectedSkipped: false,
				expectedErr:     fmt.Errorf("list item alias otherField in foreach conflicts with existing field"),
			}
		}(),
		// Test 7: Calling a function on a field that does not support it
		func() checkEvaluatorTestStructure {
			primaryField := "primary.field"

			pFValue, _ := types.MakeType("bool", true)
			fields := map[string]types.IType{primaryField: pFValue}

			checks := []string{"primary.field.gt(10)"}

			return checkEvaluatorTestStructure{
				primaryField:    primaryField,
				fields:          fields,
				checks:          checks,
				expectedRes:     nil,
				expectedSkipped: false,
				expectedErr:     fmt.Errorf("bool does not have method gt"),
			}
		}(),
		// Test 8: Calling function with wrong type
		func() checkEvaluatorTestStructure {
			primaryField := "primary.field"

			pFValue, _ := types.MakeType("bool", true)
			fields := map[string]types.IType{primaryField: pFValue}

			checks := []string{"this.eq(10)"}

			return checkEvaluatorTestStructure{
				primaryField:    primaryField,
				fields:          fields,
				checks:          checks,
				expectedRes:     nil,
				expectedSkipped: false,
				expectedErr:     fmt.Errorf("bool.eq expects a bool argument"),
			}
		}(),
	}

	for _, test := range tests {
		// Create evaluator
		evaluator := NewCheckEvaluator()

		// Evaluate check
		for _, check := range test.checks {
			res, skipped, err := evaluator.Evaluate(check, test.primaryField, test.fields, test.optMissingFields)
			errMessage := ""
			if err != nil {
				errMessage = err.Error()
			}
			expectedErrMessage := ""
			if test.expectedErr != nil {
				expectedErrMessage = test.expectedErr.Error()
			}
			if !reflect.DeepEqual(res, test.expectedRes) || !reflect.DeepEqual(skipped, test.expectedSkipped) || errMessage != expectedErrMessage {
				t.Errorf("Evaluate(%v, %v, %v, %v) = %v, %v, %v, want %v, %v, %v", test.primaryField, test.fields, test.optMissingFields, check, res, skipped, errMessage, test.expectedRes, test.expectedSkipped, expectedErrMessage)
			}
		}
	}
}

// TestEvaluateSyntaxErrors tests the parser's ability
// to detect syntax errors.
func TestEvaluateSyntaxErrors(t *testing.T) {
	tests := []checkEvaluatorTestStructure{
		// Test 1: Missing parenthesis
		func() checkEvaluatorTestStructure {
			checks := []string{"if(eq(true){ eq(false) }"}

			return checkEvaluatorTestStructure{
				primaryField:     "",
				fields:           map[string]types.IType{},
				optMissingFields: map[string]bool{},
				checks:           checks,
				expectedErr:      fmt.Errorf("syntax errors: line 1:11 missing ')' at '{'"),
			}
		}(),
		// Test 2: Extraneous token
		func() checkEvaluatorTestStructure {
			checks := []string{"if(eq(true)) unexpected_token { eq(false) }"}

			return checkEvaluatorTestStructure{
				primaryField:     "",
				fields:           map[string]types.IType{},
				optMissingFields: map[string]bool{},
				checks:           checks,
				expectedErr:      fmt.Errorf("syntax errors: line 1:13 extraneous input 'unexpected_token' expecting '{'"),
			}
		}(),
		// Test 3: Extraneous input
		func() checkEvaluatorTestStructure {
			checks := []string{"if(eq(true)) { eq(false) } else else { eq(true) }"}

			return checkEvaluatorTestStructure{
				primaryField:     "",
				fields:           map[string]types.IType{},
				optMissingFields: map[string]bool{},
				checks:           checks,
				expectedErr:      fmt.Errorf("syntax errors: line 1:32 extraneous input 'else' expecting '{'"),
			}
		}(),
		// Test 4: Invalid field expression
		func() checkEvaluatorTestStructure {
			checks := []string{"if(primary..field){ eq(false) }"}

			return checkEvaluatorTestStructure{
				primaryField:     "",
				fields:           map[string]types.IType{},
				optMissingFields: map[string]bool{},
				checks:           checks,
				expectedErr:      fmt.Errorf("syntax errors: line 1:11 extraneous input '.' expecting IDENTIFIER; line 1:17 no viable alternative at input 'field)'"),
			}
		}(),
		// Test 5: Missing curly brace in IF statement
		func() checkEvaluatorTestStructure {
			checks := []string{"if(eq(true)) eq(false) "}

			return checkEvaluatorTestStructure{
				primaryField:     "",
				fields:           map[string]types.IType{},
				optMissingFields: map[string]bool{},
				checks:           checks,
				expectedErr:      fmt.Errorf("syntax errors: line 1:13 missing '{' at 'eq'; line 1:23 missing '}' at '<EOF>'"),
			}
		}(),
		// Test 6: Using square brackets in field (this is not allowed in CMCL)
		func() checkEvaluatorTestStructure {
			checks := []string{"if(field[0]){ eq(false) }"}

			return checkEvaluatorTestStructure{
				primaryField:     "",
				fields:           map[string]types.IType{},
				optMissingFields: map[string]bool{},
				checks:           checks,
				expectedErr:      fmt.Errorf("syntax errors: line 1:8 token recognition error at: '['; line 1:10 token recognition error at: ']'; line 1:9 extraneous input '0' expecting ')'"),
			}
		}(),
		// Test 7: AND expression missing right side
		func() checkEvaluatorTestStructure {
			checks := []string{"eq(5) &&"}

			return checkEvaluatorTestStructure{
				primaryField:     "",
				fields:           map[string]types.IType{},
				optMissingFields: map[string]bool{},
				checks:           checks,
				expectedErr:      fmt.Errorf("syntax errors: line 1:8 mismatched input '<EOF>' expecting {'(', '!', IDENTIFIER}"),
			}
		}(),
		// Test 8: OR expression missing right side
		func() checkEvaluatorTestStructure {
			checks := []string{"eq(5) ||"}

			return checkEvaluatorTestStructure{
				primaryField:     "",
				fields:           map[string]types.IType{},
				optMissingFields: map[string]bool{},
				checks:           checks,
				expectedErr:      fmt.Errorf("syntax errors: line 1:8 mismatched input '<EOF>' expecting {'(', '!', IDENTIFIER}"),
			}
		}(),
	}

	for _, test := range tests {
		// Create evaluator
		evaluator := NewCheckEvaluator()

		// Evaluate check
		for _, check := range test.checks {
			res, skipped, err := evaluator.Evaluate(check, test.primaryField, test.fields, test.optMissingFields)
			errMessage := ""
			if err != nil {
				errMessage = err.Error()
			}
			expectedErrMessage := ""
			if test.expectedErr != nil {
				expectedErrMessage = test.expectedErr.Error()
			}
			if !reflect.DeepEqual(res, test.expectedRes) || !reflect.DeepEqual(skipped, test.expectedSkipped) || errMessage != expectedErrMessage {
				t.Errorf("Evaluate(%v, %v, %v, %v) = %v, %v, %v, want %v, %v, %v", test.primaryField, test.fields, test.optMissingFields, check, res, skipped, errMessage, test.expectedRes, test.expectedSkipped, expectedErrMessage)
			}
		}
	}
}
