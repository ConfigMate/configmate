package analyzer

import (
	reflect "reflect"
	"testing"

	"github.com/ConfigMate/configmate/analyzer/types"
)

func TestEvaluate(t *testing.T) {
	tests := []struct {
		// Primary field
		primaryFieldType  string      // Primary field type
		primaryFieldValue interface{} // Primary field value

		// Fields
		fields []struct {
			fieldName  string      // Field name
			fieldType  string      // Field type
			fieldValue interface{} // Field value
		}

		optMissingFields map[string]bool // Optional missing fields
		check            string          // Check

		// Expected values
		expectedResType  string      // Expected result type
		expectedResValue interface{} // Expected result value
		skipped          bool        // Expected skipped value
		expectedErr      string      // Expected error
	}{
		// Test 1
		{
			// Input values
			primaryFieldType:  "string",
			primaryFieldValue: "hello",
			fields: []struct {
				fieldName  string
				fieldType  string
				fieldValue interface{}
			}{
				{
					fieldName:  "field1.something",
					fieldType:  "string",
					fieldValue: "world",
				},
			},
			optMissingFields: map[string]bool{},
			check:            "eq(\"world\")",

			// Expected values
			expectedResType:  "bool",
			expectedResValue: false,
			skipped:          false,
			expectedErr:      "string.eq failed: hello != world",
		},
		// Test 2
		{
			// Input values
			primaryFieldType:  "int",
			primaryFieldValue: 5,
			fields: []struct {
				fieldName  string
				fieldType  string
				fieldValue interface{}
			}{
				{
					fieldName:  "field1.othernumber",
					fieldType:  "float",
					fieldValue: 3.14,
				},
			},
			optMissingFields: map[string]bool{},
			check:            "gt(field1.othernumber.toInt())",

			// Expected values
			expectedResType:  "bool",
			expectedResValue: true,
			skipped:          false,
			expectedErr:      "",
		},
	}

	for _, test := range tests {
		// Create primary field (input should be valid)
		primaryField, _ := types.MakeType(test.primaryFieldType, test.primaryFieldValue)

		// Create fields
		fields := make(map[string]types.IType)
		for _, field := range test.fields {
			// Create field (input should be valid)
			t, _ := types.MakeType(field.fieldType, field.fieldValue)
			fields[field.fieldName] = t
		}

		// Create expected result (input should be valid)
		expectedRes, _ := types.MakeType(test.expectedResType, test.expectedResValue)

		// Create evaluator
		evaluator := NewCheckEvaluator(primaryField, fields, test.optMissingFields)

		res, skipped, err := evaluator.Evaluate(test.check)
		if err != nil && err.Error() != test.expectedErr {
			t.Errorf("Evaluate(%q) = _, _, %q, want _, _, %q", test.check, err, test.expectedErr)
		} else if err == nil && (!reflect.DeepEqual(res, expectedRes) || skipped != test.skipped) {
			t.Errorf("Evaluate(%q) = %q, %v, _, want %q, %v, _", test.check, res, skipped, expectedRes, test.skipped)
		}
	}
}
