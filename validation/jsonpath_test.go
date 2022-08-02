package validation_test

import (
	"testing"

	"code.1533b4dc0.de/prskr/nurse/validation"
)

type jsonPathValidator_EqualsTestCase[V validation.Value] struct {
	testName string
	expected V
	jsonPath string
	json     string
	wantErr  bool
}

func (tt jsonPathValidator_EqualsTestCase[V]) name() string {
	return tt.testName
}

//nolint:thelper // is not a helper
func (tt jsonPathValidator_EqualsTestCase[V]) run(t *testing.T) {
	t.Parallel()
	t.Helper()
	validator, err := validation.JSONPathValidatorFor(tt.jsonPath, tt.expected)
	if err != nil {
		t.Fatalf("JSONPathValidatorFor() err = %v", err)
	}

	if err := validator.Equals(tt.json); err != nil {
		if !tt.wantErr {
			t.Errorf("Failed to equal value in %s to %v: %v", tt.json, tt.expected, err)
		}
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
			wantErr:  false,
		},
		jsonPathValidator_EqualsTestCase[string]{
			testName: "Simple object navigation - number as string",
			expected: "42",
			jsonPath: "$.number",
			json:     `{"number": 42}`,
			wantErr:  false,
		},
		jsonPathValidator_EqualsTestCase[string]{
			testName: "Simple array navigation",
			expected: "world",
			jsonPath: "$[1]",
			json:     `["hello", "world"]`,
			wantErr:  false,
		},
		jsonPathValidator_EqualsTestCase[int]{
			testName: "Simple array navigation - string to int",
			expected: 37,
			jsonPath: "$[1]",
			json:     `["13", "37"]`,
			wantErr:  false,
		},
		jsonPathValidator_EqualsTestCase[int]{
			testName: "Simple array navigation - string to int - wrong value",
			expected: 42,
			jsonPath: "$[1]",
			json:     `["13", "37"]`,
			wantErr:  true,
		},
		jsonPathValidator_EqualsTestCase[string]{
			testName: "Simple array navigation - int to string",
			expected: "37",
			jsonPath: "$[1]",
			json:     `[13, 37]`,
			wantErr:  false,
		},
	}
	//nolint:paralleltest
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name(), tt.run)
	}
}
