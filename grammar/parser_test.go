package grammar_test

import (
	"testing"

	"github.com/maxatome/go-testdeep/td"

	"github.com/baez90/nurse/grammar"
	"github.com/baez90/nurse/internal/values"
)

var wantParsedScript = td.Struct(new(grammar.Script), td.StructFields{
	"Checks": td.Bag(
		grammar.Check{
			Initiator: &grammar.Call{
				Module: "http",
				Name:   "Get",
				Params: params(grammar.Param{String: values.StringP("https://www.gogol.com/")}),
			},
			Validators: &grammar.Filters{
				Chain: []grammar.Call{
					{
						Name: "Status",
						Params: []grammar.Param{
							{
								Int: values.IntP(404),
							},
						},
					},
				},
			},
		},
		grammar.Check{
			Initiator: &grammar.Call{
				Module: "http",
				Name:   "Get",
				Params: params(grammar.Param{String: values.StringP("https://www.microsoft.com/")}),
			},
			Validators: &grammar.Filters{
				Chain: []grammar.Call{
					{
						Name: "Status",
						Params: []grammar.Param{
							{
								Int: values.IntP(200),
							},
						},
					},
					{
						Name: "Header",
						Params: []grammar.Param{
							{
								String: values.StringP("Content-Type"),
							},
							{
								String: values.StringP("text/html"),
							},
						},
					},
				},
			},
		},
	),
})

func TestParser_Parse(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		rawRule string
		want    any
		wantErr bool
	}{
		{
			name:    "Check - Initiator only - string argument",
			rawRule: `http.Get("https://www.microsoft.com/")`,
			want: &grammar.Script{
				Checks: []grammar.Check{
					{
						Initiator: &grammar.Call{
							Module: "http",
							Name:   "Get",
							Params: params(grammar.Param{String: values.StringP("https://www.microsoft.com/")}),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "Check - Initiator only - raw string argument",
			rawRule: "http.Post(\"https://www.microsoft.com/\", `{\"Name\":\"Ted.Tester\"}`)",
			want: &grammar.Script{
				Checks: []grammar.Check{
					{
						Initiator: &grammar.Call{
							Module: "http",
							Name:   "Post",
							Params: []grammar.Param{
								{
									String: values.StringP("https://www.microsoft.com/"),
								},
								{
									String: values.StringP(`{"Name":"Ted.Tester"}`),
								},
							},
						},
					},
				},
			},
		},
		{
			name:    "Check - Initiator and single filter",
			rawRule: `http.Get("https://www.microsoft.com/") => Status(200)`,
			want: &grammar.Script{
				Checks: []grammar.Check{
					{
						Initiator: &grammar.Call{
							Module: "http",
							Name:   "Get",
							Params: params(grammar.Param{String: values.StringP("https://www.microsoft.com/")}),
						},
						Validators: &grammar.Filters{
							Chain: []grammar.Call{
								{
									Name: "Status",
									Params: []grammar.Param{
										{
											Int: values.IntP(200),
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:    "Check - Initiator and multiple filters",
			rawRule: `http.Get("https://www.microsoft.com/") => Status(200) -> Header("Content-Type", "text/html")`,
			want: &grammar.Script{
				Checks: []grammar.Check{
					{
						Initiator: &grammar.Call{
							Module: "http",
							Name:   "Get",
							Params: params(grammar.Param{String: values.StringP("https://www.microsoft.com/")}),
						},
						Validators: &grammar.Filters{
							Chain: []grammar.Call{
								{
									Name: "Status",
									Params: []grammar.Param{
										{
											Int: values.IntP(200),
										},
									},
								},
								{
									Name: "Header",
									Params: []grammar.Param{
										{
											String: values.StringP("Content-Type"),
										},
										{
											String: values.StringP("text/html"),
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "CheckScript without comments",
			rawRule: `
http.Get("https://www.gogol.com/") => Status(404)
http.Get("https://www.microsoft.com/") => Status(200) -> Header("Content-Type", "text/html")
`,
			want: wantParsedScript,
		},
		{
			name: "CheckScript without comments - single line",
			//nolint:lll // required at this point
			rawRule: `http.Get("https://www.gogol.com/") => Status(404) http.Get("https://www.microsoft.com/") => Status(200) -> Header("Content-Type", "text/html")`,
			want:    wantParsedScript,
		},
		{
			name: "CheckScript with comments",
			rawRule: `
# GET https://www.gogol.com/ expect a not found response
http.Get("https://www.gogol.com/") => Status(404)

// GET https://www.microsoft.com/ - expect status OK and HTML content
http.Get("https://www.microsoft.com/") => Status(200) -> Header("Content-Type", "text/html")
`,

			want: wantParsedScript,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			p, err := grammar.NewParser[grammar.Script]()
			if err != nil {
				t.Fatalf("NewParser() err = %v", err)
			}

			got, err := p.Parse(tt.rawRule)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			td.Cmp(t, got, tt.want)
		})
	}
}
