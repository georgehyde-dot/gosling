package lexer

import (
	"gosling/token"
	"log"
	"os"
	"path/filepath"
)

type Lexer struct {
	input        string // file contents in a string
	Location     token.TokenLocation
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

// I split LexFile, LexRepl, and New out here,
// but they may be better as one later
func LexFile(f string) *Lexer {
	if filepath.Ext(f) != ".gos" {
		// not sure if this is best, but it can work for now
		return nil
	}
	contents, err := os.ReadFile(f)
	if err != nil {
		log.Fatalf("failed to open file %s\n", f)
	}
	input := string(contents)
	l := New(input)
	l.Location.Filename = f
	l.readChar()
	return l
}

func New(input string) *Lexer {
	l := &Lexer{
		input: input,
		Location: token.TokenLocation{
			Line:     0,
			LineCh:   1,
			Filename: "",
		},
	}
	return l
}

func LexRepl(in string) *Lexer {
	l := New(in)
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	if l.ch == byte('\n') {
		l.Location.Line++
		l.Location.LineCh = 0
	} else {
		l.Location.LineCh++
	}

	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}

}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhiteSpace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tokLocation := token.TokenLocation{
				Line:     l.Location.Line,
				LineCh:   l.Location.LineCh,
				Filename: l.Location.Filename,
			}
			tok = token.Token{Type: token.EQ, Literal: literal, Location: tokLocation}
		} else {
			tok = newToken(token.ASSIGN, l.ch, l.Location)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch, l.Location)
	case '(':
		tok = newToken(token.LPAREN, l.ch, l.Location)
	case ')':
		tok = newToken(token.RPAREN, l.ch, l.Location)
	case ',':
		tok = newToken(token.COMMA, l.ch, l.Location)
	case '+':
		tok = newToken(token.PLUS, l.ch, l.Location)
	case '{':
		tok = newToken(token.LBRACE, l.ch, l.Location)
	case '}':
		tok = newToken(token.RBRACE, l.ch, l.Location)
	case '-':
		tok = newToken(token.MINUS, l.ch, l.Location)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal, Location: tok.Location}
		} else {
			tok = newToken(token.BANG, l.ch, tok.Location)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch, l.Location)
	case '*':
		tok = newToken(token.ASTERISK, l.ch, l.Location)
	case '%':
		tok = newToken(token.MOD, l.ch, l.Location)
	case '<':
		tok = newToken(token.LT, l.ch, l.Location)
	case '>':
		tok = newToken(token.GT, l.ch, l.Location)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch, l.Location)
			log.Printf("illegal token at line: %d char: %d\n", l.Location.Line, l.Location.LineCh)
		}
	}

	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch byte, loc token.TokenLocation) token.Token {
	tokLocation := token.TokenLocation{
		Line:     loc.Line,
		LineCh:   loc.LineCh,
		Filename: loc.Filename,
	}
	return token.Token{Type: tokenType, Literal: string(ch), Location: tokLocation}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && '9' >= ch
}

func (l *Lexer) skipWhiteSpace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}
