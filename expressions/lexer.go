package expressions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/alecthomas/participle/lexer/regex"
)

// Lexer defines the token types that will be parsed out of an expression
var Lexer = lexer.Must(regex.New(`
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
