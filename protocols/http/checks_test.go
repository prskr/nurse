package http_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/maxatome/go-testdeep/td"

	"github.com/baez90/nurse/grammar"
	httpcheck "github.com/baez90/nurse/protocols/http"
)

func TestChecks_Execute(t *testing.T) {
	t.Parallel()

	httpModule := httpcheck.Module()

	type serverResponse struct {
		status int
		body   io.Reader
		err    error
	}

	tests := []struct {
		name    string
		check   string
		resp    serverResponse
		wantErr bool
	}{
		{
			name:  "GET check without validation",
			check: `http.GET("%s/api/books")`,
			resp: serverResponse{
				status: 200,
			},
			wantErr: false,
		},
		{
			name:  "GET check - status validation",
			check: `http.GET("%s/api/books") => Status(200)`,
			resp: serverResponse{
				status: 200,
			},
			wantErr: false,
		},
		{
			name:  "GET check - JSON path validation",
			check: `http.GET("%s/api/books") => JSONPath("$.firstName", "Homer")`,
			resp: serverResponse{
				status: 200,
				body:   strings.NewReader(`{"firstName": "Homer"}`),
			},
			wantErr: false,
		},
		{
			name:  "GET check - Status and JSON path validation",
			check: `http.GET("%s/api/books") => Status(200) -> JSONPath("$.firstName", "Homer")`,
			resp: serverResponse{
				status: 200,
				body:   strings.NewReader(`{"firstName": "Homer"}`),
			},
			wantErr: false,
		},
		{
			name:  "POST check without validation",
			check: `http.POST("%s/api/books")`,
			resp: serverResponse{
				status: 204,
			},
			wantErr: false,
		},
		{
			name:  "POST check - Status validation",
			check: `http.POST("%s/api/books") => Status(204)`,
			resp: serverResponse{
				status: 204,
			},
			wantErr: false,
		},
		{
			name:  "PUT check without validation",
			check: `http.PUT("%s/api/books/1")`,
			resp: serverResponse{
				status: 200,
			},
			wantErr: false,
		},
		{
			name:  "DELETE check without validation",
			check: `http.DELETE("%s/api/books/1")`,
			resp: serverResponse{
				status: 200,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				if tt.resp.err != nil {
					writer.WriteHeader(http.StatusInternalServerError)
					_, _ = writer.Write([]byte(tt.resp.err.Error()))
					return
				}

				writer.WriteHeader(tt.resp.status)
				if tt.resp.body != nil {
					_, _ = io.Copy(writer, tt.resp.body)
				}
			}))

			t.Cleanup(testServer.Close)

			parser, err := grammar.NewParser[grammar.Check]()
			td.CmpNoError(t, err, "grammar.NewParser()")
			parsedCheck, err := parser.Parse(fmt.Sprintf(tt.check, testServer.URL))
			td.CmpNoError(t, err, "parser.Parse()")

			chk, err := httpModule.Lookup(*parsedCheck, nil)
			td.CmpNoError(t, err, "http.LookupCheck()")

			if clientInjectable, ok := chk.(httpcheck.ClientInjectable); !ok {
				t.Fatal("Failed to inject client to check")
			} else {
				clientInjectable.SetClient(testServer.Client())
			}

			if tt.wantErr {
				td.CmpError(t, chk.Execute(context.Background()))
			} else {
				td.CmpNoError(t, chk.Execute(context.Background()))
			}
		})
	}
}
