package evaluator

import (
	"fmt"
	"gosling/lexer"
	"gosling/object"
	"gosling/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"24", 24},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
		{"4 % 2", 0},
		{"4 % 3", 1},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()

	return Eval(program)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object not Integer, got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value, got=%d want=%d",
			result.Value, expected)
		return false
	}

	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 == 1", true},
		{"1 == 2", false},
		{"1 != 2", true},
		{"1 != 1", false},
		{"false != false", false},
		{"true != false", true},
		{"true == false", false},
		{"true == true", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object not Boolean, got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value, got=%t want=%t",
			result.Value, expected)
		return false
	}

	return true
}

func testErrorObject(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.Error)
	if !ok {
		t.Errorf("object not Error, got=%T (%+v)", obj, obj)
		return false
	}
	if result.Inspect() != expected {
		t.Errorf("object has wrong value, got=%s want=%s",
			result.Inspect(), expected)
		return false
	}

	return true
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!!true", true},
		{"!!false", false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestBangOperatorError(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"!5", fmt.Sprintf("file: %s line: %d char: %d %s", "", 0, 0, "unknown operator: !INTEGER")},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testErrorObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (false) { 10 } else { 20 }", 20},
		{"if (true) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testEmptyObject(t, evaluated)
		}
	}
}

func testEmptyObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL, got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{"if (10 > 1) { return 10; }", 10},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			fmt.Sprintf("file: %s line: %d char: %d %s", "", 0, 4, "unknown operator: INTEGER + BOOLEAN"),
		},
		{
			"5 + true; 5;",
			fmt.Sprintf("file: %s line: %d char: %d %s", "", 0, 4, "unknown operator: INTEGER + BOOLEAN"),
		},
		{
			"-true",
			fmt.Sprintf("file: %s line: %d char: %d %s", "", 0, 2, "unknown operator: -BOOLEAN"),
		},
		{
			"true + false;",
			fmt.Sprintf("file: %s line: %d char: %d %s", "", 0, 7, "unknown operator: BOOLEAN + BOOLEAN"),
		},
		{
			"5; true + false; 5",
			fmt.Sprintf("file: %s line: %d char: %d %s", "", 0, 10, "unknown operator: BOOLEAN + BOOLEAN"),
		},
		{
			"if (10 > 1) { true + false; }",
			fmt.Sprintf("file: %s line: %d char: %d %s", "", 0, 21, "unknown operator: BOOLEAN + BOOLEAN"),
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned, got=%T(%+v)",
				evaluated, evaluated)
			continue
		}

		if errObj.Inspect() != tt.expectedMessage {
			t.Errorf("wrong error message, expected=%q, got=%q",
				tt.expectedMessage, errObj.Inspect())
		}
	}
}
