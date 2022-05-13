package grammar

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

func NewParser[T any]() (*Parser[T], error) {
	def, err := lexer.NewSimple([]lexer.SimpleRule{
		{Name: "Comment", Pattern: `(?:#|//)[^\n]*\n?`},
		{Name: `Module`, Pattern: `[a-z]{1}[A-z0-9]+`},
		{Name: `Ident`, Pattern: `[A-Z][a-zA-Z0-9_]*`},
		{Name: `CIDR`, Pattern: `(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}/(3[0-2]|[1-2][0-9]|[1-9])`},
		{Name: `IP`, Pattern: `(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`},
		{Name: `Float`, Pattern: `\d+\.\d+`},
		{Name: `Int`, Pattern: `[-]?\d+`},
		{Name: `RawString`, Pattern: "`[^`]*`"},
		{Name: `String`, Pattern: `'[^']*'|"[^"]*"`},
		{Name: `Arrows`, Pattern: `(->|=>)`},
		{Name: "whitespace", Pattern: `\s+`},
		{Name: "Punct", Pattern: `[-[!@#$%^&*()+_={}\|:;\."'<,>?/]|]`},
	})
	if err != nil {
		return nil, err
	}

	grammarParser, err := participle.Build(
		new(T),
		participle.Lexer(def),
		participle.Unquote("String", "RawString"),
		participle.Elide("Comment"),
	)
	if err != nil {
		return nil, err
	}

	return &Parser[T]{grammarParser: grammarParser}, nil
}

type Parser[T any] struct {
	grammarParser *participle.Parser
}

func (p Parser[T]) Parse(rawRule string) (*T, error) {
	into := new(T)
	if err := p.grammarParser.ParseString("", rawRule, into); err != nil {
		return nil, err
	}

	return into, nil
}

func (p Parser[T]) ParseBytes(data []byte) (*T, error) {
	into := new(T)
	if err := p.grammarParser.ParseBytes("", data, into); err != nil {
		return nil, err
	}

	return into, nil
}
