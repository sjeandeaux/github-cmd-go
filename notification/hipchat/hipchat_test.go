package hipchat

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNotifier_Send(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.RequestURI {
		case "/v2/room/666/notification":
			w.WriteHeader(http.StatusNoContent)
		case "/v2/room/boom/notification":
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		//TODO assert the payload
		//TODO assert headers

	}))
	defer ts.Close()

	type fields struct {
		token      string
		url        string
		httpClient *http.Client
	}
	type args struct {
		message io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "it should be ok",
			fields: fields{
				token:      "my toto ken",
				url:        ts.URL + "/v2/room/666/notification",
				httpClient: ts.Client(),
			},
			args: args{
				message: nil,
			},
			wantErr: false,
		},
		{
			name: "it should be ok",
			fields: fields{
				token:      "my toto ken",
				url:        ts.URL + "/v2/room/boom/notification",
				httpClient: ts.Client(),
			},
			args: args{
				message: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Notifier{
				token:      tt.fields.token,
				url:        tt.fields.url,
				httpClient: tt.fields.httpClient,
			}
			if err := n.Send(tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("Notifier.Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
