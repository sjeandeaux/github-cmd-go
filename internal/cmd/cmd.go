package cmd

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

//ErrNoData in the pipe
var errNoData = errors.New("No pipe")

//IsNoData true because no data
func IsNoData(err error) bool {
	return errNoData == err
}

//Input get the good input
//if data is here returns a reader on data
//if file is here returns the file
//if we are in pipe mode return the os.Stdin
func Input(data, file string, stdin *os.File) (io.ReadCloser, error) {
	if data != "" {
		return ioutil.NopCloser(strings.NewReader(data)), nil
	}

	if file != "" {
		return os.Open(file)
	}

	fi, err := stdin.Stat()
	if err != nil {
		return nil, err
	}
	if fi.Mode()&os.ModeNamedPipe == 0 {
		return nil, errNoData
	}
	return os.Stdin, nil
}
