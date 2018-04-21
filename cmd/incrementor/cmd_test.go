package main

import (
	"bytes"
	"io"
	"testing"

	"github.com/sjeandeaux/toolators/semver"
	"github.com/stretchr/testify/assert"
)

func TestCommandLineIncrementOk(t *testing.T) {
	var inputs = []struct {
		input    commandLine
		expected *semver.Version
	}{
		{
			input: commandLine{
				position: semver.PositionMajor,
				version:  "0.1.0",
			},
			expected: &semver.Version{
				Major: 1,
				Minor: 0,
				Patch: 0,
			},
		},
	}

	for _, data := range inputs {
		actual, _ := data.input.increment()
		assert.Equal(t, data.expected, actual)
	}
}

func TestCommandLineIncrementKo(t *testing.T) {

	var inputs = []struct {
		input       commandLine
		expectedErr string
	}{
		{
			input: commandLine{
				position: semver.PositionMajor,
				version:  "bad version",
			},
			expectedErr: "",
		},
		{
			input: commandLine{
				position: "bad position",
				version:  "0.1.0",
			},
			expectedErr: "non",
		},
	}

	for _, data := range inputs {
		_, actual := data.input.increment()
		assert.Error(t, actual, data.expectedErr)
	}

}

func Test_commandLine_main(t *testing.T) {
	type fields struct {
		position string
		version  string
		stdout   io.Writer
		stderr   io.Writer
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
				position: "minor",
				version:  "0.1.0",
				stdout:   bytes.NewBufferString(""),
				stderr:   bytes.NewBufferString(""),
			},
			wants: wants{
				exitCode: 0,
				stdout:   "0.2.0",
				stderr:   "",
			},
		},
		{
			name: "ko",
			fields: fields{
				position: "ko",
				version:  "0.1.0",
				stdout:   bytes.NewBufferString(""),
				stderr:   bytes.NewBufferString(""),
			},
			wants: wants{
				exitCode: 1,
				stdout:   "",
				stderr:   "\"ko\" is unknown",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &commandLine{
				position: tt.fields.position,
				version:  tt.fields.version,
				stdout:   tt.fields.stdout,
				stderr:   tt.fields.stderr,
			}
			assert.Equal(t, c.main(), tt.wants.exitCode)
			assert.Equal(t, tt.wants.stdout, c.stdout.(*bytes.Buffer).String())
			assert.Equal(t, tt.wants.stderr, c.stderr.(*bytes.Buffer).String())
		})
	}
}
