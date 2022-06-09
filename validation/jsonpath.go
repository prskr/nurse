package validation

import (
	"encoding/json"
	"fmt"

	"github.com/PaesslerAG/jsonpath"
)

var _ ValueComparator = (*JSONPathValidator)(nil)

func JSONPathValidatorFor[T Value](path string, want T) (*JSONPathValidator, error) {
	comparator, err := JSONValueComparatorFor(want)
	if err != nil {
		return nil, err
	}

	return &JSONPathValidator{
		Path:       path,
		Comparator: comparator,
	}, nil
}

type JSONPathValidator struct {
	Path       string
	Comparator ValueComparator
}

func (j JSONPathValidator) Equals(got any) error {
	parsed, err := parse(got)
	if err != nil {
		return err
	}
	val, err := jsonpath.Get(j.Path, parsed)
	if err != nil {
		return err
	}

	return j.Comparator.Equals(val)
}

func parse(in any) (any, error) {
	keyValue := make(map[string]any)
	arr := make([]any, 0)
	switch data := in.(type) {
	case []byte:
		if err := json.Unmarshal(data, &keyValue); err == nil {
			return keyValue, nil
		}

		err := json.Unmarshal(data, &arr)

		return arr, err
	case string:
		raw := []byte(data)
		if err := json.Unmarshal(raw, &keyValue); err == nil {
			return keyValue, nil
		}
		err := json.Unmarshal(raw, &arr)

		return arr, err
	}

	return nil, fmt.Errorf("cannot convert %v to JSON structure", in)
}
