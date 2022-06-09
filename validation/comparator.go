package validation

import (
	"fmt"
	"math"
)

const equalityThreshold = 0.00000001

var (
	_ ValueComparator = (*GenericComparator[int])(nil)
	_ ValueComparator = FloatComparator(0)
)

type ValueComparator interface {
	Equals(got any) error
}

type GenericComparator[T int | string] struct {
	Want   T
	Parser func(got any) (T, error)
}

func (g GenericComparator[T]) Equals(got any) error {
	parsed, err := g.Parser(got)
	if err != nil {
		return err
	}

	if parsed != g.Want {
		return fmt.Errorf("want %v but got %v", g.Want, parsed)
	}

	return nil
}

type FloatComparator float64

func (f FloatComparator) Equals(got any) error {
	val, err := ParseJSONFloat(got)
	if err != nil {
		return err
	}

	if math.Abs(float64(f)-val) > equalityThreshold {
		return fmt.Errorf("want %f but got %f", float64(f), val)
	}

	return nil
}
