package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/alecthomas/participle/lexer"
	"github.com/stretchr/testify/suite"
)

type mainSuite struct {
	suite.Suite
}

func TestMain(t *testing.T) {
	suite.Run(t, new(mainSuite))
}

func (s *mainSuite) SetupTest() {
}

func (s *mainSuite) Test_LexingSimple() {
	sym := exprLexer.Symbols()
	sbr := lexer.SymbolsByRune(exprLexer)

	tests := []struct {
		in  string
		typ string
		val string
	}{
		{"null", "Null", "null"},
		{"false", "Bool", "false"},
		{"true", "Bool", "true"},
		{"711", "Integer", "711"},
		{"-9.2", "Float", "-9.2"},
		// "0xff" Worry about hex later
		{"-2.99e-2", "Float", "-2.99e-2"},
		{"'Mona the Octocat'", "String", "'Mona the Octocat'"},
		{"'It''s open source!'", "String", "'It''s open source!'"},
		{" ", "Whitespace", " "},
		{"!", "Not", "!"},
		{"<", "Operator", "<"},
		{"<=", "Operator", "<="},
		{">", "Operator", ">"},
		{">=", "Operator", ">="},
		{"==", "Operator", "=="},
		{"!=", "Operator", "!="},
		{"&&", "Operator", "&&"},
		{"||", "Operator", "||"},
	}

	// fmt.Printf("---- exprLexer exprLexer.Symbols() %+v\n", exprLexer.Symbols())

	for _, tc := range tests {
		s.Run(tc.in, func() {
			lx, err := exprLexer.Lex(strings.NewReader(tc.in))
			s.NoError(err)

			toks, err := lexer.ConsumeAll(lx)
			s.NoError(err)

			s.Len(toks, 2)
			tok := toks[0] // Only care about the first, the second should be EOF

			s.Equal(sym[tc.typ], tok.Type, fmt.Sprintf("Expected %s got %v", tc.typ, sbr[tok.Type]))
			s.Equal(tc.val, tok.Value)
		})
	}
}

func (s *mainSuite) Test_ParsingSimple() {
	tests := []struct {
		in     string
		testFn func(*mainSuite, *Literal)
	}{
		{"null", func(s *mainSuite, res *Literal) {
			s.NotNil(res.Nil)
			s.True(res.IsNil())
		}},
		{"false", func(s *mainSuite, res *Literal) {
			s.NotNil(res.Bool)
			s.False(res.True())
		}},
		{"true", func(s *mainSuite, res *Literal) {
			s.NotNil(res.Bool)
			s.True(res.True())
		}},
		{"'qwebe'", func(s *mainSuite, res *Literal) {
			s.NotNil(res.Str)
			s.Equal(*res.Str, "'qwebe'")
		}},
		{"711", func(s *mainSuite, res *Literal) {
			s.NotNil(res.Int)
			s.Equal(*res.Int, int64(711))
		}},
		{"-9.2", func(s *mainSuite, res *Literal) {
			s.NotNil(res.Float)
			s.Equal(*res.Float, float64(-9.2))
		}},
		// // "0xff" Worry about hex later
		{"-2.99e-2", func(s *mainSuite, res *Literal) {
			s.NotNil(res.Float)
			s.Equal(*res.Float, float64(-2.99e-2))
		}},
		{"'Mona the Octocat'", func(s *mainSuite, res *Literal) {
			s.NotNil(res.Str)
			s.Equal(*res.Str, "'Mona the Octocat'")
		}},
		{"'It''s open source!'", func(s *mainSuite, res *Literal) {
			s.NotNil(res.Str)
			s.Equal(*res.Str, "'It''s open source!'")
		}},
	}

	for _, tc := range tests {
		s.Run(tc.in, func() {
			res := parse(tc.in)
			tc.testFn(s, res.LHS)
		})
	}
}

func (s *mainSuite) Test_ParsingExpression() {
	tests := []struct {
		in     string
		testFn func(*mainSuite, *expr)
	}{
		{"true == true", func(s *mainSuite, res *expr) {
			s.True(res.LHS.True())
			s.Equal("==", res.OpString())
			s.True(res.RHS.True())
		}},
		{"1 < 2", func(s *mainSuite, res *expr) {
			lhs, err := res.LHS.Evaluate(nil)
			s.NoError(err)
			s.Equal(int64(1), lhs)
			s.NotNil(res.Op)
		}},
	}

	for _, tc := range tests {
		s.Run(tc.in, func() {
			res := parse(tc.in)
			tc.testFn(s, res)
		})
	}
}

func (s *mainSuite) Test_BasicEvaluation() {
	tests := []struct {
		in  string
		out interface{}
	}{
		{"null", nil},
		{"'bob'", "'bob'"},
		{"123", int64(123)},
		{"1.23", float64(1.23)},
		{"true", true},
		{"'hi' == 'hi'", true},
		{"'hi' == 'hello'", false},
		{"123 == 123", true},
		{"123 == 321", false},
		{"1.23 == 1.23", true},
		{"1.23 == 3.21", false},
		{"true == true", true},
		{"true == false", false},
		{"1 < 2", true},
		{"2 < 1", false},
		{"1 != 2", true},
		{"1 != 'asdf'", true},
		{"'asdf' != 'asdf'", false},
	}

	for _, tc := range tests {
		s.Run(tc.in, func() {
			ast := parse(tc.in)
			out, err := evaluate(ast, nil)

			s.NoError(err)
			s.Equal(tc.out, out, printTokens(tc.in))
		})
	}
}

func (s *mainSuite) Test_ObjectLiterals() {
	tests := []struct {
		in      string
		context Context
		out     interface{}
	}{
		{"github.token", Context{"github": {"token": "i-am-a-token"}}, "i-am-a-token"},
		{"github.event.base_ref", Context{"github": {"event": map[string]interface{}{"base_ref": "i-am-the-base-ref"}}}, "i-am-the-base-ref"},
		{"github.some-sha == github.other-sha", Context{"github": {"some-sha": "asdf", "other-sha": "asdf"}}, true},
		{"github.some-sha == github.other-sha", Context{"github": {"some-sha": "asdf", "other-sha": "qwer"}}, false},
	}

	for _, tc := range tests {
		s.Run(tc.in, func() {
			ast := parse(tc.in)
			out, err := evaluate(ast, tc.context)

			s.NoError(err)
			s.Equal(tc.out, out, printTokens(tc.in))
		})
	}
}

func printTokens(in string) string {
	lx, _ := exprLexer.Lex(strings.NewReader(in))
	toks, _ := lexer.ConsumeAll(lx)
	return fmt.Sprintf("%d tokens found: %+v\n", len(toks), toks)
}
