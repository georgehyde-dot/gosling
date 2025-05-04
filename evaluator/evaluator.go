package evaluator

import (
	"gosling/ast"
	"gosling/object"
	"gosling/token"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right, node.Token.Location)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right, node.Token.Location)
	default:
		result, ok := node.(*ast.ExpressionStatement)
		if !ok {
			loc := token.TokenLocation{
				Line:     -1,
				LineCh:   -1,
				Filename: "",
			}
			return &object.Error{
				Value:    "invalid operation",
				Location: loc,
			}
		}
		return &object.Error{
			Value:    "invalid operation",
			Location: result.Token.Location,
		}
	}
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement)
	}

	return result
}

func evalPrefixExpression(operator string, right object.Object, loc token.TokenLocation) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right, loc)
	case "-":
		return evalMinusPrefixOperatorExpression(right, loc)
	default:
		return &object.Error{
			Value:    "invalid operation",
			Location: loc,
		}
	}
}

func evalInfixExpression(operator string, left, right object.Object, loc token.TokenLocation) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right, loc)
	default:
		return &object.Error{
			Value:    "invalid operation",
			Location: loc,
		}
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
		return &object.Integer{Value: leftVal / rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "%":
		return &object.Integer{Value: leftVal % rightVal}
	default:
		return &object.Error{
			Value:    "invalid operation",
			Location: loc,
		}
	}
}

func evalBangOperatorExpression(right object.Object, loc token.TokenLocation) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	default:
		return &object.Error{
			Value:    "invalid operation",
			Location: loc,
		}
	}
}

func evalMinusPrefixOperatorExpression(right object.Object, loc token.TokenLocation) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return &object.Error{
			Value:    "invalid operation",
			Location: loc,
		}
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

// rather than creating a new instance of the boolean object
// just evaluate it and return a predefined var from startup time
func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}

	return FALSE
}
