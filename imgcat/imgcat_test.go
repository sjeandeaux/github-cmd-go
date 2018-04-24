package imgcat

import (
	"io"
	"os"
	"path/filepath"
	"testing"
)

var fileOK = filepath.Join("testdata", "giphy.gif")

func TestPrint(t *testing.T) {
	t.Skip("imgcat is for iterm")
	gif, _ := os.Open(fileOK)

	type args struct {
		read io.ReadCloser
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "it should be OK",
			args: args{read: gif},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Print(tt.args.read)
		})
	}
}
