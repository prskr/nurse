package grammar

type CheckUnmarshaler interface {
	UnmarshalCheck(c Check) error
}

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

type Script struct {
	Checks []Check `parser:"@@*"`
}
