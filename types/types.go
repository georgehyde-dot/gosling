package types

type Type interface {
	IsType(Type) bool
	IsAssignableTo(Type) bool
	Kind() string // "INTEGER", "STRING", "BOOLEAN", "ARRAY", "STRUCT", "FUNCTION", "BUILTIN_FUNCTION", "ERROR"
}

// IntegerType represents the type of an integer value.
type IntegerType struct{}

func (t *IntegerType) IsType(other Type) bool {
	_, ok := other.(*IntegerType)
	return ok
}
func (t *IntegerType) IsAssignableTo(other Type) bool {
	return t.IsType(other)
}
func (t *IntegerType) Kind() string {
	return "INTEGER"
}

// StringType represents the type of a string value.
type StringType struct{}

func (t *StringType) IsType(other Type) bool {
	_, ok := other.(*StringType)
	return ok
}
func (t *StringType) IsAssignableTo(other Type) bool {
	return t.IsType(other)
}
func (t *StringType) Kind() string {
	return "STRING"
}

// BooleanType represents the type of a boolean value.
type BooleanType struct{}

func (t *BooleanType) IsType(other Type) bool {
	_, ok := other.(*BooleanType)
	return ok
}
func (t *BooleanType) IsAssignableTo(other Type) bool {
	return t.IsType(other)
}
func (t *BooleanType) Kind() string {
	return "BOOLEAN"
}

// ArrayType represents the type of an array value.
type ArrayType struct {
	ElementType Type
}

func (t *ArrayType) IsType(other Type) bool {
	otherArray, ok := other.(*ArrayType)
	if !ok {
		return false
	}
	return t.ElementType.IsType(otherArray.ElementType)
}
func (t *ArrayType) IsAssignableTo(other Type) bool {
	otherArray, ok := other.(*ArrayType)
	if !ok {
		return false
	}
	return t.ElementType.IsAssignableTo(otherArray.ElementType)
}
func (t *ArrayType) Kind() string {
	return "ARRAY"
}
