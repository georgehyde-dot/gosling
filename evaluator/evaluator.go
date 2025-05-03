package evaluator

import (
	"gosling/ast"
	"gosling/object"
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
		return evalPrefixExpression(node.Operator, right, node.Token.Line, node.Token.LineCh, node.Token.Filename)

	default:
		result, ok := node.(*ast.ExpressionStatement)
		if !ok {
			return &object.Error{
				Value:    "invalid operation",
				Line:     -1,
				LineCh:   -1,
				Filename: "",
			}
		}
		return &object.Error{
			Value:    "invalid operation",
			Line:     result.Token.Line,
			LineCh:   result.Token.LineCh,
			Filename: result.Token.Filename,
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

func evalPrefixExpression(operator string, right object.Object, line int, ch int, file string) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right, line, ch, file)
	case "-":
		return evalMinusPrefixOperatorExpression(right, line, ch, file)
	default:
		return &object.Error{
			Value:    "invalid operation",
			Line:     line,
			LineCh:   ch,
			Filename: file,
		}
	}
}

func evalBangOperatorExpression(right object.Object, line int, ch int, file string) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	default:
		return &object.Error{
			Value:    "invalid operation",
			Line:     line,
			LineCh:   ch,
			Filename: file,
		}
	}
}

func evalMinusPrefixOperatorExpression(right object.Object, line int, ch int, file string) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return &object.Error{
			Value:    "invalid operation",
			Line:     line,
			LineCh:   ch,
			Filename: file,
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
