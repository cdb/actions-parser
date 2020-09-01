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

type expr struct {
	Literal *Literal `@@`
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
		Ident = [[:ascii:]][\w\d]*
	  `))

	parser = participle.MustBuild(&expr{},
		participle.Lexer(exprLexer),
	)
)
