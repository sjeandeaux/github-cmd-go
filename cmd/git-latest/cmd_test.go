package main

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/sjeandeaux/toolators/semver"

	"github.com/stretchr/testify/assert"
)

func GetCurrentVersion(err string) func() (*semver.Version, error) {
	return func() (*semver.Version, error) {
		if err != "" {
			return nil, errors.New(err)
		}
		return &semver.Version{Major: 6, Minor: 6, Patch: 6}, nil
	}
}

func Test_commandLine_main(t *testing.T) {
	type fields struct {
		stdout     io.Writer
		stderr     io.Writer
		gitVersion func() (*semver.Version, error)
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
				gitVersion: GetCurrentVersion(""),
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
				gitVersion: GetCurrentVersion("Houston problem"),
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
				gitVersion: tt.fields.gitVersion,
			}
			c.Stdout = tt.fields.stdout
			c.Stderr = tt.fields.stderr
			assert.Equal(t, c.main(), tt.wants.exitCode)
			assert.Equal(t, tt.wants.stdout, c.Stdout.(*bytes.Buffer).String())
			assert.Equal(t, tt.wants.stderr, c.Stderr.(*bytes.Buffer).String())
		})
	}
}
