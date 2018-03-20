package github

import (
	"io"
	"os"
	"path/filepath"
	"testing"
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
