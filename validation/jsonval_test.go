package validation_test

import (
	"testing"

	"github.com/baez90/nurse/validation"
)

type testCase interface {
	run(t *testing.T)
	name() string
}

type jsonValueComparator_EqualsTestCase[V validation.Value] struct {
	testName string
	expected V
	got      any
	want     bool
}

func (tt jsonValueComparator_EqualsTestCase[V]) run(t *testing.T) {
	t.Parallel()
	t.Helper()
	comparator, err := validation.JSONValueComparatorFor(tt.expected)
	if err != nil {
		t.Fatalf("validation.JSONValueComparatorFor() err = %v", err)
	}

	if got := comparator.Equals(tt.got); got != tt.want {
		t.Errorf("Equals() = %v, want %v", got, tt.want)
	}
}

func (tt jsonValueComparator_EqualsTestCase[V]) name() string {
	return tt.testName
}

func TestJSONValueComparator_Equals(t *testing.T) {
	t.Parallel()
	tests := []testCase{
		jsonValueComparator_EqualsTestCase[int]{
			testName: "Test int equality",
			expected: 42,
			got:      42,
			want:     true,
		},
		jsonValueComparator_EqualsTestCase[int]{
			testName: "Test int equality - wrong value",
			expected: 42,
			got:      43,
			want:     false,
		},
		jsonValueComparator_EqualsTestCase[int]{
			testName: "Test int equality - string value",
			expected: 42,
			got:      "42",
			want:     true,
		},
		jsonValueComparator_EqualsTestCase[int]{
			testName: "Test int equality - []byte value",
			expected: 42,
			got:      []byte("42"),
			want:     true,
		},
		jsonValueComparator_EqualsTestCase[int]{
			testName: "Test int equality - float value",
			expected: 42,
			got:      42.0,
			want:     true,
		},
		jsonValueComparator_EqualsTestCase[int8]{
			testName: "Test int8 equality",
			expected: 42,
			got:      42,
			want:     true,
		},
		jsonValueComparator_EqualsTestCase[int8]{
			testName: "Test int8 equality - wrong value",
			expected: 42,
			got:      43,
			want:     false,
		},
		jsonValueComparator_EqualsTestCase[int8]{
			testName: "Test int8 equality - int16 value",
			expected: 42,
			got:      int16(42),
			want:     true,
		},
		jsonValueComparator_EqualsTestCase[int8]{
			testName: "Test int8 equality - uint16 value",
			expected: 42,
			got:      uint16(42),
			want:     true,
		},
		jsonValueComparator_EqualsTestCase[float32]{
			testName: "Test float32 equality - float value",
			expected: 42.0,
			got:      42.0,
			want:     true,
		},
		jsonValueComparator_EqualsTestCase[float32]{
			testName: "Test float32 equality - float value",
			expected: 42.0,
			got:      float64(42.0),
			want:     true,
		},
		jsonValueComparator_EqualsTestCase[float64]{
			testName: "Test float64 equality - float value",
			expected: 42.0,
			got:      42.0,
			want:     true,
		},
		jsonValueComparator_EqualsTestCase[float64]{
			testName: "Test float64 equality - int value",
			expected: 42.0,
			got:      42,
			want:     true,
		},
		jsonValueComparator_EqualsTestCase[float64]{
			testName: "Test float64 equality - []byte value",
			expected: 42.0,
			got:      []byte("42"),
			want:     true,
		},
		jsonValueComparator_EqualsTestCase[float64]{
			testName: "Test float64 equality - float32 value",
			expected: 42.0,
			got:      float32(42.0),
			want:     true,
		},
		jsonValueComparator_EqualsTestCase[float64]{
			testName: "Test float64 equality - string value",
			expected: 42.0,
			got:      "42.0",
			want:     true,
		},
		jsonValueComparator_EqualsTestCase[float64]{
			testName: "Test float64 equality - string value without dot",
			expected: 42.0,
			got:      "42",
			want:     true,
		},
		jsonValueComparator_EqualsTestCase[string]{
			testName: "Test string equality",
			expected: "hello",
			got:      "hello",
			want:     true,
		},
		jsonValueComparator_EqualsTestCase[string]{
			testName: "Test string equality - []byte value",
			expected: "hello",
			got:      []byte("hello"),
			want:     true,
		},
		jsonValueComparator_EqualsTestCase[string]{
			testName: "Test string equality - int value",
			expected: "1337",
			got:      1337,
			want:     true,
		},
		jsonValueComparator_EqualsTestCase[string]{
			testName: "Test string equality - float value",
			expected: "13.37",
			got:      13.37,
			want:     true,
		},
		jsonValueComparator_EqualsTestCase[string]{
			testName: "Test string equality - wrong case",
			expected: "hello",
			got:      "HELLO",
			want:     false,
		},
	}

	//nolint:paralleltest
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name(), tt.run)
	}
}
