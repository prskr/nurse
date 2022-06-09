package validation

import "math"

const equalityThreshold = 0.00000001

var (
	_ ValueComparator = (*GenericComparator[int])(nil)
	_ ValueComparator = FloatComparator(0)
)

type ValueComparator interface {
	Equals(got any) bool
}

type GenericComparator[T int | string] struct {
	Want   T
	Parser func(got any) (T, error)
}

func (g GenericComparator[T]) Equals(got any) bool {
	parsed, err := g.Parser(got)
	if err != nil {
		return false
	}

	return parsed == g.Want
}

type FloatComparator float64

func (f FloatComparator) Equals(got any) bool {
	val, err := ParseJSONFloat(got)
	if err != nil {
		return false
	}

	return math.Abs(float64(f)-val) < equalityThreshold
}
