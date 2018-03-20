package main

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/sjeandeaux/github-cmd-go/semver"

	"github.com/stretchr/testify/assert"
)

type gitVersionTest struct {
	err string
}

func Test_defaultGitVersion_GetCurrentVersion(t *testing.T) {
	defaultGitVersion := defaultGitVersion{}
	value, err := defaultGitVersion.GetCurrentVersion()
	assert.NotNil(t, value)
	assert.Nil(t, err)
}

func (g *gitVersionTest) GetCurrentVersion() (*semver.Version, error) {
	if g.err != "" {
		return nil, errors.New(g.err)
	}
	return &semver.Version{Major: 6, Minor: 6, Patch: 6}, nil
}

func Test_commandLine_main(t *testing.T) {
	type fields struct {
		stdout     io.Writer
		stderr     io.Writer
		gitVersion gitVersion
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
			name: "ok",
			fields: fields{
				gitVersion: &gitVersionTest{err: ""},
				stdout:     bytes.NewBufferString(""),
				stderr:     bytes.NewBufferString(""),
			},
			wants: wants{
				exitCode: 0,
				stdout:   "6.6.6",
				stderr:   "",
			},
		},

		{
			name: "ok",
			fields: fields{
				gitVersion: &gitVersionTest{err: "Houston problem"},
				stdout:     bytes.NewBufferString(""),
				stderr:     bytes.NewBufferString(""),
			},
			wants: wants{
				exitCode: 1,
				stdout:   "",
				stderr:   "Houston problem",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &commandLine{
				stdout:     tt.fields.stdout,
				stderr:     tt.fields.stderr,
				gitVersion: tt.fields.gitVersion,
			}
			assert.Equal(t, c.main(), tt.wants.exitCode)
			assert.Equal(t, tt.wants.stdout, c.stdout.(*bytes.Buffer).String())
			assert.Equal(t, tt.wants.stderr, c.stderr.(*bytes.Buffer).String())
		})
	}
}
