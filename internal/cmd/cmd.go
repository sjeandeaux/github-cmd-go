package cmd

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
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
		fileReader, err := os.Open(file)
		if err == nil {
			return fileReader, err
		}

		_, err = url.Parse(file)
		if err != nil {
			return nil, err
		}
		println(file)
		req, _ := http.NewRequest(http.MethodGet, file, nil)
		//TODO client
		resp, err := http.DefaultClient.Do(req)
		return resp.Body, err

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
