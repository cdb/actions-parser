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

func parse(in string) interface{} {
	ast := &expr{}

	err := parser.Parse(strings.NewReader(in), ast)
	if err != nil {
		fmt.Printf("----  err %+v\n", err)

	}
	return ast
}

type expr struct {
	Symbol string `@Ident`
}

var (
	exprLexer = lexer.Must(regex.New(`
		Null = null
		Bool = false|true
		Number = -?(?:0|[1-9]\d*)(?:\.\d+)?(?:[eE][+-]?\d+)?
		String = '([^\\']|'')*'
		Whitespace = \s+
		Ident = [[:ascii:]][\w\d]*
	  `))

	parser = participle.MustBuild(&expr{},
		participle.Lexer(exprLexer),
	)
)
