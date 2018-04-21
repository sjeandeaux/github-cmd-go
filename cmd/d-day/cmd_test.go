package main

import (
	"bytes"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	timeNow = func() time.Time {
		return time.Date(2018, time.April, 20, 0, 0, 0, 0, time.UTC)
	}
}

func Test_commandLine_main(t *testing.T) {
	type fields struct {
		date   string
		stdout io.Writer
		stderr io.Writer
		stdin  *os.File
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
				date:   "2018-05-04",
				stdout: bytes.NewBufferString(""),
				stderr: bytes.NewBufferString(""),
			},
			wants: wants{
				exitCode: 0,
				stdout:   "14",
				stderr:   "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &commandLine{
				date:   tt.fields.date,
				stdout: tt.fields.stdout,
				stderr: tt.fields.stderr,
				stdin:  tt.fields.stdin,
			}
			assert.Equal(t, c.main(), tt.wants.exitCode)
			assert.Equal(t, tt.wants.stdout, c.stdout.(*bytes.Buffer).String())
			assert.Equal(t, tt.wants.stderr, c.stderr.(*bytes.Buffer).String())
		})
	}
}
