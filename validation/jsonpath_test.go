package validation_test

import (
	"testing"

	"github.com/baez90/nurse/validation"
)

type jsonPathValidator_EqualsTestCase[V validation.Value] struct {
	testName string
	expected V
	jsonPath string
	json     string
	want     bool
}

func (tt jsonPathValidator_EqualsTestCase[V]) name() string {
	return tt.testName
}

func (tt jsonPathValidator_EqualsTestCase[V]) run(t *testing.T) {
	t.Parallel()
	t.Helper()
	validator, err := validation.JSONPathValidatorFor(tt.jsonPath, tt.expected)
	if err != nil {
		t.Fatalf("JSONPathValidatorFor() err = %v", err)
	}

	if validator.Equals(tt.json) != tt.want {
		t.Errorf("Failed to equal value in %s to %v", tt.json, tt.expected)
	}
}

func TestJSONPathValidator_Equals(t *testing.T) {
	t.Parallel()
	tests := []testCase{
		jsonPathValidator_EqualsTestCase[string]{
			testName: "Simple object navigation",
			expected: "hello",
			jsonPath: "$.greeting",
			json:     `{"greeting": "hello"}`,
			want:     true,
		},
		jsonPathValidator_EqualsTestCase[string]{
			testName: "Simple object navigation - number as string",
			expected: "42",
			jsonPath: "$.number",
			json:     `{"number": 42}`,
			want:     true,
		},
		jsonPathValidator_EqualsTestCase[string]{
			testName: "Simple array navigation",
			expected: "world",
			jsonPath: "$[1]",
			json:     `["hello", "world"]`,
			want:     true,
		},
		jsonPathValidator_EqualsTestCase[int]{
			testName: "Simple array navigation - string to int",
			expected: 37,
			jsonPath: "$[1]",
			json:     `["13", "37"]`,
			want:     true,
		},
		jsonPathValidator_EqualsTestCase[int]{
			testName: "Simple array navigation - string to int - wrong value",
			expected: 42,
			jsonPath: "$[1]",
			json:     `["13", "37"]`,
			want:     false,
		},
		jsonPathValidator_EqualsTestCase[string]{
			testName: "Simple array navigation - int to string",
			expected: "37",
			jsonPath: "$[1]",
			json:     `[13, 37]`,
			want:     true,
		},
	}
	//nolint:paralleltest
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name(), tt.run)
	}
}
