package main

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/sjeandeaux/toolators/github"
	"github.com/stretchr/testify/assert"
)

type githubClientTest struct {
	errorCreateRelease   error
	errorGetReleaseByTag error
	errorUpload          error
}

func (g *githubClientTest) CreateRelease(edit *github.EditRelease) (*github.Release, error) {
	return &github.Release{
		UploadURLTemplate: "TODO check",
		TagName:           "TODO check",
		URL:               "TODO check",
	}, g.errorCreateRelease
}

func (g *githubClientTest) GetReleaseByTag(tag string) (*github.Release, error) {
	return &github.Release{
		UploadURLTemplate: "TODO check",
		TagName:           "TODO check",
		URL:               "TODO check",
	}, g.errorGetReleaseByTag
}

func (g *githubClientTest) Upload(urlPath string, u github.UploadInformation) error {
	//TODO check
	return g.errorUpload
}

func Test_commandLine_main(t *testing.T) {
	type fields struct {
		token        string
		owner        string
		repo         string
		create       bool
		path         string
		tag          string
		name         string
		label        string
		contentType  string
		githubClient githubClient
		stdout       io.Writer
		stderr       io.Writer
	}
	type wants struct {
		exitCode int
		stdout   string
		stderr   string
	}

	tests := []struct {
		name   string
		fields fields
		wants  wants
	}{
		{
			name: "ok create true",
			fields: fields{
				token:       "",
				owner:       "",
				repo:        "",
				create:      true,
				path:        "",
				tag:         "",
				name:        "",
				label:       "",
				contentType: "",
				githubClient: &githubClientTest{
					errorCreateRelease:   nil,
					errorGetReleaseByTag: nil,
					errorUpload:          nil,
				},
				stdout: bytes.NewBufferString(""),
				stderr: bytes.NewBufferString(""),
			},
			wants: wants{
				exitCode: 0,
				stdout:   "",
				stderr:   "",
			},
		},
		{
			name: "ok create false",
			fields: fields{
				token:       "",
				owner:       "",
				repo:        "",
				create:      false,
				path:        "",
				tag:         "",
				name:        "",
				label:       "",
				contentType: "",
				githubClient: &githubClientTest{
					errorCreateRelease:   nil,
					errorGetReleaseByTag: nil,
					errorUpload:          nil,
				},
				stdout: bytes.NewBufferString(""),
				stderr: bytes.NewBufferString(""),
			},
			wants: wants{
				exitCode: 0,
				stdout:   "",
				stderr:   "",
			},
		},
		{
			name: "ko create false",
			fields: fields{
				token:       "",
				owner:       "",
				repo:        "",
				create:      false,
				path:        "",
				tag:         "",
				name:        "",
				label:       "",
				contentType: "",
				githubClient: &githubClientTest{
					errorCreateRelease:   nil,
					errorGetReleaseByTag: errors.New("booooum"),
					errorUpload:          nil,
				},
				stdout: bytes.NewBufferString(""),
				stderr: bytes.NewBufferString(""),
			},
			wants: wants{
				exitCode: 1,
				stdout:   "",
				stderr:   "booooum",
			},
		},
		{
			name: "ko create true",
			fields: fields{
				token:       "",
				owner:       "",
				repo:        "",
				create:      true,
				path:        "",
				tag:         "",
				name:        "",
				label:       "",
				contentType: "",
				githubClient: &githubClientTest{
					errorCreateRelease:   errors.New("booooum"),
					errorGetReleaseByTag: nil,
					errorUpload:          nil,
				},
				stdout: bytes.NewBufferString(""),
				stderr: bytes.NewBufferString(""),
			},
			wants: wants{
				exitCode: 1,
				stdout:   "",
				stderr:   "booooum",
			},
		},
		{
			name: "ko upload",
			fields: fields{
				token:       "",
				owner:       "",
				repo:        "",
				create:      true,
				path:        "",
				tag:         "",
				name:        "",
				label:       "",
				contentType: "",
				githubClient: &githubClientTest{
					errorCreateRelease:   nil,
					errorGetReleaseByTag: nil,
					errorUpload:          errors.New("booooum"),
				},
				stdout: bytes.NewBufferString(""),
				stderr: bytes.NewBufferString(""),
			},
			wants: wants{
				exitCode: 1,
				stdout:   "",
				stderr:   "booooum",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &commandLine{
				token:        tt.fields.token,
				owner:        tt.fields.owner,
				repo:         tt.fields.repo,
				create:       tt.fields.create,
				path:         tt.fields.path,
				tag:          tt.fields.tag,
				name:         tt.fields.name,
				label:        tt.fields.label,
				contentType:  tt.fields.contentType,
				githubClient: tt.fields.githubClient,
			}
			c.Stdout = tt.fields.stdout
			c.Stderr = tt.fields.stderr
			assert.Equal(t, c.main(), tt.wants.exitCode)
			assert.Equal(t, tt.wants.stdout, c.Stdout.(*bytes.Buffer).String())
			assert.Equal(t, tt.wants.stderr, c.Stderr.(*bytes.Buffer).String())
		})
	}
}
