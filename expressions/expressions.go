package expressions

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle"
)

// TODO:
// #( )	Logical grouping
// Lots more I'm sure?

var (
	parser = participle.MustBuild(&Expression{},
		participle.Lexer(Lexer),
		participle.Elide("Whitespace"),
	)
)

func parse(in string) *Expression {
	ast := &Expression{}

	err := parser.Parse(strings.NewReader(in), ast)
	if err != nil {
		fmt.Printf("----  err %+v\n", err)
	}
	return ast
}

func evaluate(ast *Expression, context Context) (interface{}, error) {
	return ast.Evaluate(context)
}
