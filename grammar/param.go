package grammar

import (
	"fmt"
)

func ValidateParameterCount(params []Param, expected int) error {
	if len(params) < expected {
		return fmt.Errorf("%w: expected %d got %d", ErrAmbiguousParamCount, expected, len(params))
	}
	return nil
}

type Param struct {
	String *string  `parser:"@String | @RawString"`
	Int    *int     `parser:"| @Int"`
	Float  *float64 `parser:"| @Float"`
}

func (p Param) AsString() (string, error) {
	if p.String == nil {
		return "", fmt.Errorf("string is nil %w", ErrTypeMismatch)
	}
	return *p.String, nil
}

func (p Param) AsInt() (int, error) {
	if p.Int == nil {
		return 0, fmt.Errorf("int is nil %w", ErrTypeMismatch)
	}
	return *p.Int, nil
}

func (p Param) AsFloat() (float64, error) {
	if p.Float == nil {
		return 0, fmt.Errorf("float is nil %w", ErrTypeMismatch)
	}
	return *p.Float, nil
}
