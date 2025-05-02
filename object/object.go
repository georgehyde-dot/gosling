package object

import (
	"fmt"
)

type ObjectType string

const (
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
	ERROR_OBJ   = "ERROR"
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
	Value    string
	Line     int
	LineCh   int
	Filename string
}

// Integer Methods
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

// Boolean Methods
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }

// I'm not doing NULL values, but I will include and error return
// Error Methods
func (e *Error) Inspect() string {
	return fmt.Sprintf("file: %s line: %d char: %d %s", e.Filename, e.Line, e.LineCh, e.Value)
}
func (e *Error) Type() ObjectType { return ERROR_OBJ }
