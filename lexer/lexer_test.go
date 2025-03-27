package lexer

import (
	"testing"

	"my_interpreter/token"
)

func TestNextToken(t *testing.T) {
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "3"},
		{token.LT, "<"},
		{token.INT, "9"},
		{token.GT, ">"},
		{token.INT, "2"},
		{token.SEMICOLON, ";"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "3"},
		{token.LT, "<"},
		{token.INT, "9"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.IDENT, "five"},
		{token.EQ, "=="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "ten"},
		{token.NOT_EQ, "!="},
		{token.INT, "11"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}
	l := LexFile("./testfile.gos")

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - tokenliteral wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestIllegalChar(t *testing.T) {
	tests := []struct {
		failLine int
		failChar int
	}{
		{1, 16},
	}
	l := LexFile("./testIllegal.gos")

	for i, tt := range tests {
	inner:
		for {
			tok := l.NextToken()
			if tok.Type == token.ILLEGAL && tt.failChar != l.lineCh && tt.failLine == l.line {
				t.Fatalf("tests[%d] - wrong char number expected=%d, got=%d", i, tt.failChar, l.lineCh)
			}
			if tok.Type == token.ILLEGAL && tt.failLine != l.line && tt.failChar == l.lineCh {
				t.Fatalf("tests[%d] - wrong line number expected=%d, got=%d", i, tt.failLine, l.line)
			}
			// fix tests
			if tok.Type != token.ILLEGAL && tt.failChar == l.lineCh && tt.failLine == l.line {
				t.Fatalf("tests[%d] - wrong token type expected=%s, got=%s", i, token.ILLEGAL, tok.Type)
			}
			if tok.Type == token.ILLEGAL {
				break inner
			}
		}
	}
}
