package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	internalcmd "github.com/sjeandeaux/toolators/internal/cmd"
)

func Test_commandLine_main(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.RequestURI {
		case "/6.6.6.OK":
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer ts.Close()

	type fields struct {
		CommandLine internalcmd.CommandLine
		httpClient  *http.Client
		action      string
		url         string
		data        string
		file        string
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "ok",
			fields: fields{
				CommandLine: internalcmd.CommandLine{
					Stdout: os.Stdout,
					Stderr: os.Stderr,
				},
				httpClient: ts.Client(),
				url:        fmt.Sprint(ts.URL, "/6.6.6.OK"),
				action:     "ok",
				data:       "the payload",
			},
			want: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &commandLine{
				CommandLine: tt.fields.CommandLine,
				httpClient:  tt.fields.httpClient,
				action:      tt.fields.action,
				url:         tt.fields.url,
				data:        tt.fields.data,
				file:        tt.fields.file,
			}
			if got := c.main(); got != tt.want {
				t.Errorf("commandLine.main() = %v, want %v", got, tt.want)
			}
		})
	}
}
