package analyzer

import (
	"fmt"
	reflect "reflect"
	"testing"

	"github.com/ConfigMate/configmate/analyzer/types"
)

func TestEvaluate(t *testing.T) {
	type testStructure struct {
		// Input values
		primaryField     types.IType          // Primary field
		fields           map[string]FieldInfo // Fields
		optMissingFields map[string]bool      // Optional missing fields
		checks           []string             // Checks

		// Expected values
		expectedRes     types.IType // Expected result
		expectedSkipped bool        // Expected skipped value
		expectedErr     error       // Expected error
	}

	// Test cases
	tests := []testStructure{
		// Test 1:
		// This is a simple case where we are checking
		// a field of type bool with a check verifing
		// is it true. No other fields involved.
		// The check should pass.
		func() testStructure {
			primaryField, _ := types.MakeType("bool", true) // Create primary field
			fields := make(map[string]FieldInfo)            // Create fields
			optMissingFields := make(map[string]bool)       // Create optional missing fields
			checks := []string{"eq(true)"}                  // Create checks

			expectedRes, _ := types.MakeType("bool", true) // Create expected result
			expectedSkipped := false                       // Create expected skipped value
			expectedErr := error(nil)                      // Create expected error

			return testStructure{
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
		func() testStructure {
			primaryField, _ := types.MakeType("bool", false) // Create primary field
			fields := make(map[string]FieldInfo)             // Create fields
			optMissingFields := make(map[string]bool)        // Create optional missing fields
			checks := []string{"eq(true)"}                   // Create checks

			expectedRes, _ := types.MakeType("bool", false)            // Create expected result
			expectedSkipped := false                                   // Create expected skipped value
			expectedErr := fmt.Errorf("bool.eq failed: true != false") // Create expected error

			return testStructure{
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
		func() testStructure {
			primaryField, _ := types.MakeType("int", 5) // Create primary field
			fields := make(map[string]FieldInfo)        // Create fields
			optMissingFields := make(map[string]bool)   // Create optional missing fields
			checks := []string{"range(0, 10).eq(true)"} // Create checks

			expectedRes, _ := types.MakeType("bool", true) // Create expected result
			expectedSkipped := false                       // Create expected skipped value
			expectedErr := error(nil)                      // Create expected error

			return testStructure{
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
		func() testStructure {
			primaryField, _ := types.MakeType("int", 5) // Create primary field
			fields := make(map[string]FieldInfo)        // Create fields
			optMissingFields := make(map[string]bool)   // Create optional missing fields
			checks := []string{"range(0, 3).eq(false)"} // Create checks

			expectedRes, _ := types.MakeType("bool", true) // Create expected result
			expectedSkipped := false                       // Create expected skipped value
			expectedErr := error(nil)                      // Create expected error

			return testStructure{
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
		func() testStructure {
			primaryField, _ := types.MakeType("float", 3.14) // Create primary field
			fields := make(map[string]FieldInfo)             // Create fields
			optMissingFields := make(map[string]bool)        // Create optional missing fields
			checks := []string{"toInt().range(0, 2)"}        // Create checks

			expectedRes, _ := types.MakeType("bool", false)                      // Create expected result
			expectedSkipped := false                                             // Create expected skipped value
			expectedErr := fmt.Errorf("int.range failed: 3 not in range [0, 2]") // Create expected error

			return testStructure{
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
		func() testStructure {
			primaryField, _ := types.MakeType("int", 5) // Create primary field

			fields := make(map[string]FieldInfo)       // Create fields
			otherField, _ := types.MakeType("int", 10) // Create other field
			fields["config.other"] = FieldInfo{Value: otherField}

			optMissingFields := make(map[string]bool) // Create optional missing fields
			checks := []string{"lt(config.other)"}    // Create checks

			expectedRes, _ := types.MakeType("bool", true) // Create expected result
			expectedSkipped := false                       // Create expected skipped value
			expectedErr := error(nil)                      // Create expected error

			return testStructure{
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
		func() testStructure {
			primaryField, _ := types.MakeType("int", 5) // Create primary field

			fields := make(map[string]FieldInfo)       // Create fields
			otherField1, _ := types.MakeType("int", 0) // Create other field 1
			otherField2, _ := types.MakeType("int", 10)
			fields["config.other1"] = FieldInfo{Value: otherField1}
			fields["config.other2"] = FieldInfo{Value: otherField2}

			optMissingFields := make(map[string]bool)                 // Create optional missing fields
			checks := []string{"range(config.other1, config.other2)"} // Create checks

			expectedRes, _ := types.MakeType("bool", true) // Create expected result
			expectedSkipped := false                       // Create expected skipped value
			expectedErr := error(nil)                      // Create expected error

			return testStructure{
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
		func() testStructure {
			primaryField, _ := types.MakeType("int", 5) // Create primary field

			fields := make(map[string]FieldInfo) // Create fields
			otherField, _ := types.MakeType("float", 10.0)
			fields["config.other"] = FieldInfo{Value: otherField}

			optMissingFields := make(map[string]bool)      // Create optional missing fields
			checks := []string{"lt(config.other.toInt())"} // Create checks

			expectedRes, _ := types.MakeType("bool", true) // Create expected result
			expectedSkipped := false                       // Create expected skipped value
			expectedErr := error(nil)                      // Create expected error

			return testStructure{
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
		func() testStructure {
			primaryField, _ := types.MakeType("int", 5) // Create primary field

			fields := make(map[string]FieldInfo) // Create fields
			otherField1, _ := types.MakeType("float", 0.5)
			otherField2, _ := types.MakeType("float", 3.14)
			fields["config.other1"] = FieldInfo{Value: otherField1}
			fields["config.something.other2"] = FieldInfo{Value: otherField2}

			optMissingFields := make(map[string]bool)                                           // Create optional missing fields
			checks := []string{"range(config.other1.toInt(), config.something.other2.toInt())"} // Create checks

			expectedRes, _ := types.MakeType("bool", false)                      // Create expected result
			expectedSkipped := false                                             // Create expected skipped value
			expectedErr := fmt.Errorf("int.range failed: 5 not in range [0, 3]") // Create expected error

			return testStructure{
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
		evaluator := NewCheckEvaluator(test.primaryField, test.fields, test.optMissingFields)

		// Evaluate checks
		for _, check := range test.checks {
			res, skipped, err := evaluator.Evaluate(check)
			if !reflect.DeepEqual(res, test.expectedRes) || !reflect.DeepEqual(skipped, test.expectedSkipped) || !reflect.DeepEqual(err, test.expectedErr) {
				t.Errorf("Evaluate(%v, %v, %v, %v) = %v, %v, %v, want %v, %v, %v", test.primaryField, test.fields, test.optMissingFields, check, res, skipped, err, test.expectedRes, test.expectedSkipped, test.expectedErr)
			}
		}
	}
}
