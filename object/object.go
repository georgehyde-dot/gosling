package object

import (
	"bytes"
	"fmt"
	"gosling/ast"
	"gosling/token"
	"strings"
)

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	ERROR_OBJ        = "ERROR"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	FUNCTION_OBJ     = "FUNCTION"
	STRING_OBJ       = "STRING"
	BUILTIN_OBJ      = "BUILTIN"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

type Boolean struct {
	Value bool
}

type Error struct {
	Message  string
	Location token.TokenLocation
}

type Null struct {
	Value string
}

type ReturnValue struct {
	Value Object
}

type Environment struct {
	store map[string]Object
	outer *Environment
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

type String struct {
	Value string
}

type Builtin struct {
	Fn BuiltinFunction
}

type BuiltinFunction func(args ...Object) Object

// Integer Methods
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

// Boolean Methods
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }

// Error Methods
func (e *Error) Inspect() string {
	return fmt.Sprintf("file: %s line: %d char: %d %s", e.Location.Filename, e.Location.Line, e.Location.LineCh, e.Message)
}
func (e *Error) Type() ObjectType { return ERROR_OBJ }

// Null Methods
func (n *Null) Inspect() string  { return "null" }
func (n *Null) Type() ObjectType { return NULL_OBJ }

// ReturnValue Methods
func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

// Helper function to create new errors
func NewError(message string, location token.TokenLocation) *Error {
	return &Error{Message: message, Location: location}
}

// Environment Methods

// Function Methods
func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

// String Methods
func (s *String) Inspect() string  { return s.Value }
func (s *String) Type() ObjectType { return STRING_OBJ }

// Builtin Methods
func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }
