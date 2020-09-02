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

	type tokInfo struct {
		t rune
		v string
	}
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
		{"<", "Less", "<"},
		{"<=", "LessOrEqual", "<="},
		{">", "Greater", ">"},
		{">=", "GreaterOrEqual", ">="},
		{"==", "Equal", "=="},
		{"!=", "NotEqual", "!="},
		{"&&", "And", "&&"},
		{"||", "Or", "||"},
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
	type tokInfo struct {
		t rune
		v string
	}
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
			tc.testFn(s, res.Literal)
		})
	}
}
