package validation

import (
	"fmt"
)

var _ ValueComparator = (*JSONValueComparator)(nil)

type Value interface {
	float32 | float64 | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | string | []byte
}

func JSONValueComparatorFor[T Value](want T) (*JSONValueComparator, error) {
	ti := any(want)
	switch in := ti.(type) {
	case float32, float64:
		return &JSONValueComparator{
			Comparator: FloatComparator(ToFloat64(in)),
		}, nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return &JSONValueComparator{
			Comparator: GenericComparator[int]{
				Want:   ToInt(in),
				Parser: ParseJSONInt,
			},
		}, nil
	case string:
		return &JSONValueComparator{
			Comparator: GenericComparator[string]{
				Want:   in,
				Parser: ParseJSONString,
			},
		}, nil
	case []byte:
		return &JSONValueComparator{
			Comparator: GenericComparator[string]{
				Want:   string(in),
				Parser: ParseJSONString,
			},
		}, nil
	default:
		return nil, fmt.Errorf("no matching type detected for %v", want)
	}
}

type JSONValueComparator struct {
	Comparator ValueComparator
}

func (j JSONValueComparator) Equals(got any) bool {
	return j.Comparator.Equals(got)
}
