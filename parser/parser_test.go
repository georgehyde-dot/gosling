package parser

import (
	"gosling/ast"
	"gosling/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	l := lexer.LexFile("./testlet.gos")
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatal("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got %d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foo"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.Tokenliteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let', got %q", s.Tokenliteral())
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s no *ast.LetStatement, got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not %s, got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.Tokenliteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral not %s, got=%s", name, letStmt.Name.Tokenliteral())
		return false
	}
	return true
}
