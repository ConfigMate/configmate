package analyzer

import (
	"fmt"
	"strings"

	"github.com/ConfigMate/configmate/analyzer/types"
	"github.com/golang-collections/collections/stack"
)

type cmclNodeType int

const (
	cmclIfCheck cmclNodeType = iota
	cmclForeachCheck
	cmclForeachItemAlias
	cmclFieldExpr
	cmclOrExpr
	cmclAndExpr
	cmclNotExpr
	cmclParenExpr
	cmclFunction
	cmclString
	cmclInt
	cmclFloat
	cmclBool
)

type cmclNode struct {
	nodeType cmclNodeType
	value    string
	children []*cmclNode

	// Used by cmclIfCheck
	elseIfStatements []*cmclNode
	elseStatement    *cmclNode
}

type CheckEvaluator struct {
	primaryField     string
	fields           map[string]types.IType
	optMissingFields map[string]bool

	// The evalFieldStack stores the ITypes of
	// the fields that functions
	// are being evaluated on
	evalFieldStack stack.Stack
}

func NewCheckEvaluator(primaryField string, fields map[string]types.IType, optMissingFields map[string]bool) *CheckEvaluator {
	// Create evaluator
	evaluator := &CheckEvaluator{
		primaryField:     primaryField,
		fields:           fields,
		optMissingFields: optMissingFields,
	}

	return evaluator
}

func (ce *CheckEvaluator) Evaluate(check string) (types.IType, bool, error) {
	// Parse check
	parser := &CheckParser{}
	node, err := parser.parse(check)
	if err != nil {
		return nil, false, err
	}

	// Get primary field value
	if pField, ok := ce.fields[ce.primaryField]; !ok {
		return nil, false, fmt.Errorf("primary field %s does not exist", ce.primaryField)
	} else {
		ce.fields["this"] = pField
		ce.evalFieldStack.Push(pField)
	}

	// Evaluate check
	return ce.visit(node)
}

func (ce *CheckEvaluator) visit(node *cmclNode) (types.IType, bool, error) {
	switch node.nodeType {
	case cmclIfCheck:
		return ce.visitIfCheck(node)
	case cmclForeachCheck:
		return ce.visitForeachCheck(node)
	case cmclFieldExpr:
		return ce.visitFieldExpr(node)
	case cmclOrExpr:
		return ce.visitOrExpr(node)
	case cmclAndExpr:
		return ce.visitAndExpr(node)
	case cmclNotExpr:
		return ce.visitNotExpr(node)
	case cmclParenExpr:
		return ce.visitParenExpr(node)
	case cmclFunction:
		return ce.visitFunction(node)
	case cmclString:
		return ce.visitString(node)
	case cmclInt:
		return ce.visitInt(node)
	case cmclFloat:
		return ce.visitFloat(node)
	case cmclBool:
		return ce.visitBool(node)
	default:
		return nil, false, fmt.Errorf("unknown node type %v", node.nodeType)
	}
}

func (ce *CheckEvaluator) visitIfCheck(node *cmclNode) (types.IType, bool, error) {
	// Evaluate if statement
	condition, skipping, err := ce.visit(node.children[0])
	if condition == nil {
		return nil, false, err
	} else if skipping {
		return nil, true, nil
	}

	// Check if the condition is bool
	if _, ok := condition.Value().(bool); !ok {
		return nil, false, fmt.Errorf("if statement condition must be a bool")
	}

	// Check if the condition is true
	if condition.Value().(bool) {
		// Evaluate if statement
		return ce.visit(node.children[1])
	}

	// Evaluate else if statements
	for _, elseIfStatement := range node.elseIfStatements {
		condition, skipping, err := ce.visit(elseIfStatement.children[0])
		if condition == nil {
			return nil, false, err
		} else if skipping {
			return nil, true, nil
		}

		// Check if the condition is bool
		if _, ok := condition.Value().(bool); !ok {
			return nil, false, fmt.Errorf("else if statement condition must be a bool")
		}

		// Check if the condition is true
		if condition.Value().(bool) {
			// Evaluate else if statement
			return ce.visit(elseIfStatement.children[1])
		}
	}

	// Evaluate else statement
	if node.elseStatement != nil {
		return ce.visit(node.elseStatement.children[0])
	}

	// Make bool true to return
	t, _ := types.MakeType("bool", true)
	return t, false, nil
}

func (ce *CheckEvaluator) visitForeachCheck(node *cmclNode) (types.IType, bool, error) {
	// Get alias for list items during evaluation
	alias := node.children[0].value

	if _, ok := ce.fields[alias]; ok {
		return nil, false, fmt.Errorf("alias %s in foreach conflicts with existing field", alias)
	}

	// Evaluate field expression. Result should be a list
	list, skipping, err := ce.visit(node.children[1])
	if list == nil {
		return nil, false, err
	} else if skipping {
		return nil, true, nil
	}

	// Check if the list is a list
	if !strings.HasSuffix(list.TypeName(), "list") {
		return nil, false, fmt.Errorf("foreach statement must be a list")
	}

	// Evaluate foreach statement. Overall result will be true
	// if all foreach body are true
	resultErrors := make([]error, 0)
	for i, value := range list.Value().([]types.IType) {
		// Add alias to for list item
		ce.fields[alias] = value

		// Evaluate foreach body
		result, skipping, err := ce.visit(node.children[2])
		if err != nil {
			return nil, false, err
		} else if skipping {
			return nil, true, nil
		}

		// Check if the result is a bool
		if _, ok := result.Value().(bool); !ok {
			return nil, false, fmt.Errorf("foreach body must evaluate to a bool")
		}

		// Return if the result is false
		if !result.Value().(bool) {
			resultErrors = append(resultErrors, fmt.Errorf("check failed for item %d", i))
		}

		// Pop value from stack
		ce.evalFieldStack.Pop()
	}

	if len(resultErrors) > 0 {
		return nil, false, fmt.Errorf("foreach body failed: %v", resultErrors)
	}

	// Make bool true to return
	t, _ := types.MakeType("bool", true)
	return t, false, nil
}

func (ce *CheckEvaluator) visitFieldExpr(node *cmclNode) (types.IType, bool, error) {
	// Get field name
	fieldName := node.value
	// Check if the field exists
	if field, ok := ce.fields[fieldName]; ok {
		// Push field value to stack
		ce.evalFieldStack.Push(field)

		// Apply functions
		var fErr error // Function error
		for _, f := range node.children {
			// Evaluate function
			result, skipping, err := ce.visit(f)
			if result == nil {
				return nil, false, err
			} else if skipping {
				return nil, true, nil
			} else {
				fErr = err
			}
		}

		// Pop field value from stack
		result := ce.evalFieldStack.Pop().(types.IType)

		// Return result
		return result, false, fErr

	} else if ce.optMissingFields[fieldName] {
		// Skipping check because optional field is missing
		return nil, true, fmt.Errorf("skipping check because referenced optional field %s is missing", fieldName)
	}

	return nil, false, fmt.Errorf("field %s does not exist", fieldName)
}

func (ce *CheckEvaluator) visitOrExpr(node *cmclNode) (types.IType, bool, error) {
	// Evaluate left expression
	left, skipping, err := ce.visit(node.children[0])
	if left == nil {
		return nil, false, err
	} else if skipping {
		return nil, true, nil
	}

	// Check if the left expression is bool
	if left.TypeName() != "bool" {
		return nil, false, fmt.Errorf("or expression left expression must be a bool")
	}

	// Check if the left expression is true
	if left.Value().(bool) {
		// Make bool true to return
		t, _ := types.MakeType("bool", true)
		return t, false, nil
	}

	// Check if there is a right expression
	if len(node.children) < 2 {
		// Make bool false to return
		t, _ := types.MakeType("bool", false)
		return t, false, err
	}

	var errs []error
	errs = append(errs, err)

	// Evaluate right expression
	right, skipping, err := ce.visit(node.children[1])
	if right == nil {
		return nil, false, err
	} else if skipping {
		return nil, true, nil
	}

	// Check if the right expression is bool
	if right.TypeName() != "bool" {
		return nil, false, fmt.Errorf("or expression right expression must be a bool")
	}

	if right.Value().(bool) {
		// Make bool true to return
		t, _ := types.MakeType("bool", true)
		return t, false, nil
	}

	// Add errors
	errs = append(errs, err)

	// Return the right expression
	return right, false, fmt.Errorf("or expression failed: %v", errs)
}
