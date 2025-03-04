package ast

// expression is something that has a value.
type expression any

// Primitive represents a primitive value with generic type
type Primitive[T any] struct {
	Value T
}

// Primitives.
type (
	Int64   = Primitive[int64]
	Float64 = Primitive[float64]
	String  = Primitive[string]
	Bool    = Primitive[bool]
	Byte    = Primitive[byte]
	Rune    = Primitive[rune]
)

// Symbol is a name with a value in an environment.
type Symbol struct {
	Name string
}

// Call represents a function/method call.
type Call struct {
	Function  any
	Arguments []expression
}

// Special forms.
type (
	Assign struct {
		Target Symbol
		Value  expression
	}

	Binding struct {
		Variable Symbol
		Value    expression
	}

	Break struct {
		Value expression
	}

	Continue struct{}

	Def struct {
		Name  Symbol
		Value expression
	}

	Fun struct {
		Name       Symbol
		Parameters []Symbol
		Body       []expression
	}

	Lambda struct {
		Parameters []Symbol
		Body       []expression
	}

	Let struct {
		Bindings []Binding
		Body     []expression
	}

	Loop struct {
		Bindings  []Binding
		Condition expression
		Body      []expression
	}

	Struct struct {
		Name   Symbol
		Fields []Binding
	}

	Tie struct {
		Function any
		Args     []expression
	}

	When struct {
		Clauses []WhenClause
		Else    []expression
	}

	WhenClause struct {
		Condition expression
		Body      []expression
	}
)

// Collections.
type (
	Array []any
	Map   map[any]any
	Set   map[any]struct{}
)
