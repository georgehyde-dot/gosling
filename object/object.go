package object

import (
	"fmt"
	"gosling/token"
)

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	ERROR_OBJ        = "ERROR"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
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
