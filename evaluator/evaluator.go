package evaluator

import (
	"fmt"
	"gosling/ast"
	"gosling/object"
	"gosling/token"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{Value: "null"}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right, node.Token.Location)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		if isError(left) {
			return left
		}
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right, node.Token.Location)
	case *ast.IfExpression:
		return evalIfExpression(node)
	case *ast.BlockStatement:
		return evalBlockStatement(node)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	default:
		return object.NewError("unknown node type", token.TokenLocation{
			Line:     -1,
			LineCh:   -1,
			Filename: "",
		})
	}
}

func evalProgram(program *ast.Program) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalIfExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	} else {
		return NULL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func evalPrefixExpression(operator string, right object.Object, loc token.TokenLocation) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right, loc)
	case "-":
		return evalMinusPrefixOperatorExpression(right, loc)
	default:
		return object.NewError(fmt.Sprintf("unknown prefix operator: %s", operator), loc)
	}
}

func evalInfixExpression(operator string, left, right object.Object, loc token.TokenLocation) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right, loc)
	case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
		return evalBooleanInfixExpression(operator, left, right, loc)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	default:
		return object.NewError(fmt.Sprintf("unknown operator: %s %s %s", left.Type(), operator, right.Type()), loc)
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object, loc token.TokenLocation) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "/":
		if rightVal == 0 {
			return object.NewError("division by zero", loc)
		}
		return &object.Integer{Value: leftVal / rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "%":
		if rightVal == 0 {
			return object.NewError("modulo by zero", loc)
		}
		return &object.Integer{Value: leftVal % rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return object.NewError(fmt.Sprintf("unknown operator: %s", operator), loc)
	}
}

func evalBooleanInfixExpression(operator string, left, right object.Object, loc token.TokenLocation) object.Object {
	leftVal := left.(*object.Boolean).Value
	rightVal := right.(*object.Boolean).Value

	switch operator {
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return object.NewError(fmt.Sprintf("unknown operator: %s %s %s", object.BOOLEAN_OBJ, operator, object.BOOLEAN_OBJ), loc)
	}
}

func evalBangOperatorExpression(right object.Object, loc token.TokenLocation) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return object.NewError(fmt.Sprintf("unknown operator: !%s", right.Type()), loc)
	}
}

func evalMinusPrefixOperatorExpression(right object.Object, loc token.TokenLocation) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return object.NewError(fmt.Sprintf("unknown operator: -%s", right.Type()), loc)
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

// rather than creating a new instance of the boolean object
// just evaluate it and return a predefined var from startup time
func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}

	return FALSE
}
