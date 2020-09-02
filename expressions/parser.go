package expressions

import (
	"strings"
)

// Expression is the top level of the AST, containing a literal and an optional right hand side
// TODO: Sub-expressions to get 1 + 2 + 3
type Expression struct {
	LHS *Literal `@@`
	Op  *string  `( @Operator`
	RHS *Literal `@@ )?`
}

// Literal is a "union" type, only one value with be non-nil, depending on the token type passed in.
type Literal struct {
	Nil   *string   `@Null`
	Str   *string   `| @String`
	Int   *int64    `| @Integer`
	Float *float64  `| @Float`
	Bool  *myBool   `| @Bool`
	Func  *Function `| @@`
	Obj   *Object   `| @@`
}

// Object represents one of the context objects available in expressions (like github or env)
type Object struct {
	Head  string   `@Context`
	Props []string `{ @ContextProperty }`
}

type Function struct {
	Name string     `@Ident`
	Args []*Literal `"(" (@@ ("," @@)*)? ")"`
}

type myBool bool

// Helper functions for the types above:

func (e *Expression) opString() string {
	return *e.Op
}

// Capture interface allows for custom handling of bool values
func (mb *myBool) Capture(values []string) error {
	out := strings.Join(values, "") == "true"
	*mb = myBool(out)
	return nil
}

func (l *Literal) isNil() bool {
	return l.Nil != nil
}

func (l *Literal) isTrue() bool {
	if l.Bool != nil {
		return *l.Bool == myBool(true)
	}
	return false
}
