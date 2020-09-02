package main

import (
	"fmt"
	"regexp"
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

func evaluate(ast *expr, context Context) (interface{}, error) {
	return ast.Evaluate(context)
}

type expr struct {
	LHS *Literal `@@`
	Op  *string  `( @Operator`
	RHS *Literal `@@ )?`
}

func (e *expr) OpString() string {
	return *e.Op
}

func (e *expr) Evaluate(context Context) (interface{}, error) {
	if e.Op == nil {
		return e.LHS.Evaluate(context)
	}

	lhs, err := e.LHS.Evaluate(context)
	if err != nil {
		return nil, fmt.Errorf("Error on LSH evaluation: %w", err)
	}

	rhs, err := e.RHS.Evaluate(context)
	if err != nil {
		return nil, fmt.Errorf("Error on RHC evaluation: %w", err)
	}

	return e.evaluate(lhs, rhs)
}

func (e *expr) evaluate(lhs, rhs interface{}) (bool, error) {
	switch *e.Op {
	case "==":
		return lhs == rhs, nil
	case "!=":
		return lhs != rhs, nil
	case "<":
		l, _ := lhs.(int64)
		r, _ := rhs.(int64) // TODO: Handle errors, and floats
		return l < r, nil
	default:
		panic(fmt.Sprintf("op '%s' not implemented", *e.Op))
	}
}

// Literal is a "union" type, where only one matching value will be present.
type Literal struct {
	Pos lexer.Position

	Nil   *string  `@Null`
	Str   *string  `| @String`
	Int   *int64   `| @Integer`
	Float *float64 `| @Float`
	Bool  *myBool  `| @Bool`
	Obj   *Object  `| @@`
}

type Object struct {
	Head  string   `@Context`
	Props []string `{ @ContextProperty }`
}

var propCleaner = regexp.MustCompile(`[\.\[\]']`)

func (o *Object) Evaluate(context Context) (interface{}, error) {
	var ret interface{}
	ret, ok := context[o.Head]
	if !ok {
		return nil, fmt.Errorf("No root context named %s", o.Head)
	}

	for _, prop := range o.Props {
		prop = propCleaner.ReplaceAllString(prop, "") // Strip the leading . or wrapping [' ']
		switch t := ret.(type) {
		case map[string]interface{}:
			if v, ok := t[prop]; ok {
				ret = v
			}
		default:
			panic("Context not structured as expected")
		}
	}
	return ret, nil
}

type Context map[string]map[string]interface{}

func (l *Literal) Evaluate(context Context) (interface{}, error) {
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
	case l.Obj != nil:
		o := *l.Obj
		return o.Evaluate(context)
	default:
		panic("empty literal")
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

		Context = (github|env|job|steps|runner|secrets|strategy|matrix|needs)
		ContextProperty = (\.[\w-]+|\['[\w-]+'\])

		Operator = (<=?|>=?|==|!=|&&|\|\|)

		Not = !

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
