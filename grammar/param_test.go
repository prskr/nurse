package grammar_test

import (
	"testing"

	"github.com/baez90/nurse/grammar"
	"github.com/baez90/nurse/internal/values"
)

func TestParam_AsString(t *testing.T) {
	t.Parallel()
	type fields struct {
		String *string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "Empty string",
			fields: fields{
				String: values.StringP(""),
			},
			want: "",
		},
		{
			name: "Any string",
			fields: fields{
				String: values.StringP("Hello, world!"),
			},
			want: "Hello, world!",
		},
		{
			name:    "nil value",
			wantErr: true,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			p := grammar.Param{
				String: tt.fields.String,
			}
			got, err := p.AsString()
			if (err != nil) != tt.wantErr {
				t.Errorf("AsString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AsString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParam_AsInt(t *testing.T) {
	t.Parallel()
	type fields struct {
		Int *int
	}
	tests := []struct {
		name    string
		fields  fields
		want    int
		wantErr bool
	}{
		{
			name: "zero value",
			fields: fields{
				Int: values.IntP(0),
			},
			want: 0,
		},
		{
			name: "Any int",
			fields: fields{
				Int: values.IntP(42),
			},
			want: 42,
		},
		{
			name:    "nil value",
			wantErr: true,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			p := grammar.Param{
				Int: tt.fields.Int,
			}
			got, err := p.AsInt()
			if (err != nil) != tt.wantErr {
				t.Errorf("AsInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AsInt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParam_AsFloat(t *testing.T) {
	t.Parallel()
	type fields struct {
		Float *float64
	}
	tests := []struct {
		name    string
		fields  fields
		want    float64
		wantErr bool
	}{
		{
			name: "Zero value",
			fields: fields{
				Float: values.FloatP(0),
			},
			want: 0,
		},
		{
			name: "Any value",
			fields: fields{
				Float: values.FloatP(13.37),
			},
			want: 13.37,
		},
		{
			name:    "nil value",
			wantErr: true,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			p := grammar.Param{
				Float: tt.fields.Float,
			}
			got, err := p.AsFloat()
			if (err != nil) != tt.wantErr {
				t.Errorf("AsFloat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AsFloat() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func params(p ...grammar.Param) []grammar.Param {
	return p
}
