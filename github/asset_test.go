package github

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAsset_size(t *testing.T) {
	type fields struct {
		File        string
		Name        string
		Label       string
		ContentType string
	}
	tests := []struct {
		name    string
		fields  fields
		want    int64
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				File: filepath.Join("testdata", "data"),
			},
			want:    4,
			wantErr: false,
		},
		{
			name: "ko",
			fields: fields{
				File: "not found",
			},
			want:    -1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Asset{
				File:        tt.fields.File,
				Name:        tt.fields.Name,
				Label:       tt.fields.Label,
				ContentType: tt.fields.ContentType,
			}
			got, err := a.size()
			if (err != nil) != tt.wantErr {
				t.Errorf("Asset.size() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Asset.size() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAsset_reader(t *testing.T) {
	type fields struct {
		File        string
		Name        string
		Label       string
		ContentType string
	}

	fileNotFound, _ := os.Open("not found")
	fileFound, _ := os.Open(filepath.Join("testdata", "data"))
	tests := []struct {
		name    string
		fields  fields
		want    io.ReadCloser
		wantErr bool
	}{
		{
			name: "ko",
			fields: fields{
				File: "not found",
			},
			want:    fileNotFound,
			wantErr: true,
		},
		{
			name: "ok",
			fields: fields{
				File: filepath.Join("testdata", "data"),
			},
			want:    fileFound,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Asset{
				File:        tt.fields.File,
				Name:        tt.fields.Name,
				Label:       tt.fields.Label,
				ContentType: tt.fields.ContentType,
			}
			got, err := a.reader()
			if (err != nil) != tt.wantErr {
				t.Errorf("Asset.reader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != nil) != (tt.want != nil) {
				t.Errorf("Asset.reader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAsset_request(t *testing.T) {
	type fields struct {
		File        string
		Name        string
		Label       string
		ContentType string
	}
	type args struct {
		urlPath string
	}

	fi, _ := os.Open(filepath.Join("testdata", "data"))
	defer fi.Close()
	request, _ := http.NewRequest(http.MethodPost, "urlPath", fi)
	request.ContentLength = 4

	//header
	request.Header.Add("Content-Type", "application/binary")

	//query
	query := request.URL.Query()
	query.Add("name", "fileName")
	query.Add("label", "Label")
	request.URL.RawQuery = query.Encode()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *http.Request
		wantErr bool
	}{
		{
			name: "failed reader",
			fields: fields{
				File: "not found",
			},
			args: args{
				urlPath: "urlPath",
			},
			want:    nil,
			wantErr: true,
		},

		{
			name: "ok",
			fields: fields{
				File:        filepath.Join("testdata", "data"),
				ContentType: "application/binary",
				Name:        "fileName",
				Label:       "Label",
			},
			args: args{
				urlPath: "urlPath",
			},
			want:    request,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Asset{
				File:        tt.fields.File,
				Name:        tt.fields.Name,
				Label:       tt.fields.Label,
				ContentType: tt.fields.ContentType,
			}
			got, err := a.request(tt.args.urlPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Asset.request() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//TODO test!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
			if tt.want != nil {
				tt.want.Body = nil
			}
			if got != nil {
				got.Body = nil
			}
			assert.Equal(t, got, tt.want)
		})
	}
}
