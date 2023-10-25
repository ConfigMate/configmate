package analyzer

import (
	"fmt"
	"strconv"
	"strings"

	"go.uber.org/multierr"

	"github.com/ConfigMate/configmate/analyzer/types"
	"github.com/golang-collections/collections/stack"
)

type cmclNodeType int

const (
	cmclIfCheck cmclNodeType = iota
	cmclForeachCheck
	cmclForeachItemAlias
	cmclForeachListArg
	cmclFieldExpr
	cmclFuncExpr
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

func newCheckEvaluator(primaryField string, fields map[string]types.IType, optMissingFields map[string]bool) *CheckEvaluator {
	// Create evaluator
	evaluator := &CheckEvaluator{
		primaryField:     primaryField,
		fields:           fields,
		optMissingFields: optMissingFields,
	}

	return evaluator
}

func (ce *CheckEvaluator) evaluate(check string) (types.IType, bool, error) {
	// Parse check
	parser := &CheckParser{}
	node, err := parser.parse(check)
	if err != nil {
		return nil, false, err
	}

	// Get primary field value
	if pField, ok := ce.fields[ce.primaryField]; ok {
		ce.fields["this"] = pField
		ce.evalFieldStack.Push(pField)
	} else if ce.optMissingFields[ce.primaryField] {
		// Skipping check because optional field is missing
		// Make bool false to return
		t, _ := types.MakeType("bool", false)
		return t, true, fmt.Errorf("skipping check because primary field %s is optional and missing", ce.primaryField)
	} else {
		return nil, false, fmt.Errorf("primary field %s does not exist", ce.primaryField)
	}

	// Evaluate check
	res, skipping, err := ce.visit(node)
	if res == nil {
		return nil, false, err
	} else if skipping {
		return res, true, err
	}

	// Check if the result is a bool
	if res.TypeName() != "bool" {
		return nil, false, fmt.Errorf("check must evaluate to a bool")
	}

	return res, skipping, err
}

func (ce *CheckEvaluator) visit(node *cmclNode) (types.IType, bool, error) {
	switch node.nodeType {
	case cmclIfCheck:
		return ce.visitIfCheck(node)
	case cmclForeachCheck:
		return ce.visitForeachCheck(node)
	case cmclFieldExpr:
		return ce.visitFieldExpr(node)
	case cmclFuncExpr:
		return ce.visitFuncExpr(node)
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
		return condition, true, err
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
			return condition, true, err
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
		return nil, false, fmt.Errorf("list item alias %s in foreach conflicts with existing field", alias)
	}

	// Get list to iterate over
	listFieldName := node.children[1].value
	list, ok := ce.fields[listFieldName]
	if !ok {
		return nil, false, fmt.Errorf("field %s does not exist", listFieldName)
	}

	// Check if the list is a list
	if !strings.HasPrefix(list.TypeName(), "list:") {
		return nil, false, fmt.Errorf("foreach argument must be a list")
	}

	// Evaluate foreach statement. Overall result will be true
	// if all foreach body are true
	resultErrors := make([]error, 0)
	for i, value := range list.Value().([]types.IType) {
		// Add alias to for list item
		ce.fields[alias] = value

		// Evaluate foreach body
		result, skipping, err := ce.visit(node.children[2])
		if result == nil {
			return nil, false, err
		} else if skipping {
			return result, true, err
		}

		// Check if the result is a bool
		if _, ok := result.Value().(bool); !ok {
			return nil, false, fmt.Errorf("foreach body must evaluate to a bool")
		}

		// Collect error if the result is false
		if !result.Value().(bool) {
			resultErrors = append(resultErrors, fmt.Errorf("item %d: %v", i, err))
		}

		// Remove alias from for list item
		delete(ce.fields, alias)
	}

	if len(resultErrors) > 0 {
		// Make bool false to return
		t, _ := types.MakeType("bool", false)
		return t, false, fmt.Errorf("foreach body failed: %v", multierr.Combine(resultErrors...))
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
				return result, true, err
			}

			// Save error
			fErr = err

			// Update field value on stack
			ce.evalFieldStack.Pop()
			ce.evalFieldStack.Push(result)
		}

		// Pop field value from stack
		result := ce.evalFieldStack.Pop().(types.IType)

		// Return result
		return result, false, fErr

	} else if ce.optMissingFields[fieldName] {
		// Skipping check because optional field is missing
		// Make bool false to return
		t, _ := types.MakeType("bool", false)
		return t, true, fmt.Errorf("skipping check because referenced optional field %s is missing", fieldName)
	}

	return nil, false, fmt.Errorf("field %s does not exist", fieldName)
}

func (ce *CheckEvaluator) visitFuncExpr(node *cmclNode) (types.IType, bool, error) {
	// Place this as the node the function applies to
	node.value = "this"
	return ce.visitFieldExpr(node)
}

func (ce *CheckEvaluator) visitOrExpr(node *cmclNode) (types.IType, bool, error) {
	// Evaluate left expression
	left, skipping, err := ce.visit(node.children[0])
	if left == nil {
		return nil, false, err
	} else if skipping {
		return left, true, err
	}

	// If there is no right expression, return the left expression
	if len(node.children) < 2 {
		return left, false, err
	}

	// If there is a right expression, the left expression must be a bool
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

	var errs []error
	errs = append(errs, err)

	// Evaluate right expression
	right, skipping, err := ce.visit(node.children[1])
	if right == nil {
		return nil, false, err
	} else if skipping {
		return right, true, err
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
	return right, false, multierr.Combine(errs...)
}

func (ce *CheckEvaluator) visitAndExpr(node *cmclNode) (types.IType, bool, error) {
	// Evaluate left expression
	left, skipping, err := ce.visit(node.children[0])
	if left == nil {
		return nil, false, err
	} else if skipping {
		return left, true, err
	}

	// If there is no right expression, return the left expression
	if len(node.children) < 2 {
		return left, false, err
	}

	// If there is a right expression, the left expression must be a bool
	// Check if the left expression is bool
	if left.TypeName() != "bool" {
		return nil, false, fmt.Errorf("and expression left expression must be a bool")
	}

	// Check if the left expression is false
	if !left.Value().(bool) {
		// Make bool false to return
		t, _ := types.MakeType("bool", false)
		return t, false, err
	}

	// Evaluate right expression
	right, skipping, err := ce.visit(node.children[1])
	if right == nil {
		return nil, false, err
	} else if skipping {
		return right, true, err
	}

	// Check if the right expression is bool
	if right.TypeName() != "bool" {
		return nil, false, fmt.Errorf("and expression right expression must be a bool")
	}

	if !right.Value().(bool) {
		// Make bool false to return
		t, _ := types.MakeType("bool", false)
		return t, false, err
	}

	return right, false, nil
}

func (ce *CheckEvaluator) visitNotExpr(node *cmclNode) (types.IType, bool, error) {
	// Evaluate expression
	expr, skipping, err := ce.visit(node.children[0])
	if expr == nil {
		return nil, false, err
	} else if skipping {
		return expr, true, err
	}

	// Check if the expression is bool
	if expr.TypeName() != "bool" {
		return nil, false, fmt.Errorf("not expression value must be a bool")
	}

	// Check if the expression is true (undesired condition in this case because we are negating it)
	if expr.Value().(bool) {
		// Make bool false to return
		t, _ := types.MakeType("bool", false)
		return t, false, err
	}

	// Returns false with no error (because we are negating the expression)
	// Make bool true to return
	t, _ := types.MakeType("bool", true)
	return t, false, nil
}

func (ce *CheckEvaluator) visitParenExpr(node *cmclNode) (types.IType, bool, error) {
	// Evaluate expression
	expr, skipping, err := ce.visit(node.children[0])
	if expr == nil {
		return nil, false, err
	} else if skipping {
		return expr, true, err
	}

	return expr, false, err
}

func (ce *CheckEvaluator) visitFunction(node *cmclNode) (types.IType, bool, error) {
	// Get function name
	functionName := node.value

	// Get arguments
	args := make([]types.IType, 0)
	for _, arg := range node.children {
		// Evaluate argument
		result, skipping, err := ce.visit(arg)
		if result == nil {
			return nil, false, err
		} else if skipping {
			return result, true, err
		}

		args = append(args, result)
	}

	// Get field value
	field := ce.evalFieldStack.Peek().(types.IType)

	// Apply function
	result, err := field.GetMethod(functionName)(args)
	if result == nil {
		return nil, false, err
	}

	return result, false, err
}

func (ce *CheckEvaluator) visitString(node *cmclNode) (types.IType, bool, error) {
	// Remove quotes
	value := node.value[1 : len(node.value)-1]

	// Make string to return
	t, err := types.MakeType("string", value)
	return t, false, err
}

func (ce *CheckEvaluator) visitInt(node *cmclNode) (types.IType, bool, error) {
	// Parse int
	intValue, err := strconv.Atoi(node.value)
	if err != nil {
		return nil, false, err
	}

	// Make int to return
	t, err := types.MakeType("int", intValue)
	return t, false, err
}

func (ce *CheckEvaluator) visitFloat(node *cmclNode) (types.IType, bool, error) {
	// Parse float
	floatValue, err := strconv.ParseFloat(node.value, 64)
	if err != nil {
		return nil, false, err
	}

	// Make float to return
	t, err := types.MakeType("float", floatValue)
	return t, false, err
}

func (ce *CheckEvaluator) visitBool(node *cmclNode) (types.IType, bool, error) {
	// Parse bool
	boolValue, err := strconv.ParseBool(node.value)
	if err != nil {
		return nil, false, err
	}

	// Make bool to return
	t, err := types.MakeType("bool", boolValue)
	return t, false, err
}
