package grammar

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

var (
	_ json.Unmarshaler = (*Check)(nil)
	_ yaml.Unmarshaler = (*Check)(nil)
)

type Call struct {
	Module string  `parser:"(@Module'.')?"`
	Name   string  `parser:"@Ident"`
	Params []Param `parser:"'(' @@? ( ',' @@ )*')'"`
}

type Filters struct {
	Chain []Call `parser:"@@ ('->' @@)*"`
}

type Check struct {
	Initiator  *Call    `parser:"@@"`
	Validators *Filters `parser:"( '=>' @@)?"`
}

func (c *Check) UnmarshalYAML(value *yaml.Node) error {
	parser, err := NewParser[Check]()
	if err != nil {
		return err
	}
	chk, err := parser.Parse(value.Value)
	if err != nil {
		return err
	}

	*c = *chk
	return nil
}

func (c *Check) UnmarshalJSON(bytes []byte) error {
	parser, err := NewParser[Check]()
	if err != nil {
		return err
	}
	chk, err := parser.ParseBytes(bytes)
	if err != nil {
		return err
	}

	*c = *chk
	return nil
}

type Script struct {
	Checks []Check `parser:"(@@';'?)*"`
}
