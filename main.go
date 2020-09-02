package main

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
	"github.com/alecthomas/participle/lexer/regex"
)

func main() {
	fmt.Println("hello")
}

func parse(in string) *expr {
	ast := &expr{}

	err := parser.Parse(strings.NewReader(in), ast)
	if err != nil {
		fmt.Printf("----  err %+v\n", err)
	}
	return ast
}

func evaluate(ast *expr) (interface{}, error) {
	return ast.Evaluate()
}

type expr struct {
	LHS      *Literal  `@@`
	Operator *Operator `( @@`
	RHS      *Literal  `@@ )?`
}

func (e *expr) Evaluate() (interface{}, error) {
	if e.Operator == nil {
		return e.LHS.Evaluate()
	}

	lhs, err := e.LHS.Evaluate()
	if err != nil {
		return nil, fmt.Errorf("Error on LSH evaluation: %w", err)
	}

	rhs, err := e.RHS.Evaluate()
	if err != nil {
		return nil, fmt.Errorf("Error on RHC evaluation: %w", err)
	}
	return e.Operator.Evaluate(lhs, rhs)
}

// Literal is a "union" type, where only one matching value will be present.
type Literal struct {
	Pos lexer.Position

	Nil   *string  `@Null`
	Str   *string  `| @String`
	Int   *int64   `| @Integer`
	Float *float64 `| @Float`
	Bool  *myBool  `| @Bool`
}

func (l *Literal) Evaluate() (interface{}, error) {
	switch {
	case l.Nil != nil:
		return nil, nil
	case l.Str != nil:
		return *l.Str, nil
	case l.Int != nil:
		return *l.Int, nil
	case l.Float != nil:
		return *l.Float, nil
	case l.Bool != nil:
		return bool(*l.Bool), nil
	default:
		panic("not implemented")
	}
}

type Operator struct {
	LessOrEqual *string `@LessOrEqual`
	Equal       *string `| @Equal`
}

func (o *Operator) Kind() string {
	switch {
	case o.LessOrEqual != nil:
		return "LessOrEqual"
	case o.Equal != nil:
		return "Equal"
	default:
		panic("Missing operator kind")
	}
}

func (o *Operator) Evaluate(lhs, rhs interface{}) (bool, error) {
	switch {
	case o.Equal != nil:
		return lhs == rhs, nil
	default:
		return false, nil
	}
}

type myBool bool

func (mb *myBool) Capture(values []string) error {
	out := strings.Join(values, "") == "true"
	*mb = myBool(out)
	return nil
}

func (l *Literal) IsNil() bool {
	return l.Nil != nil
}

func (l *Literal) True() bool {
	if l.Bool != nil {
		return *l.Bool == myBool(true)
	}
	return false
}

var (
	exprLexer = lexer.Must(regex.New(`
		Null = null
		Bool = (false|true)
		Float = -?(?:0|[1-9]\d*)(?:\.\d+)(?:[eE][+-]?\d+)?
		Integer = -?\d+
		String = '([^\\']|'')*'
		Whitespace = \s+

		LessOrEqual = <=
		Less = <[^=]*
		GreaterOrEqual = >=
		Greater = >
		Equal = ==
		NotEqual = !=
		Not = !
		And = &&
		Or = \|\|

		Ident = [[:ascii:]][\w\d]*
		`))

	// TODO:
	// #( )	Logical grouping
	// [ ]	Index
	// .	Property dereference

	parser = participle.MustBuild(&expr{},
		participle.Lexer(exprLexer),
		participle.Elide("Whitespace"),
	)
)
