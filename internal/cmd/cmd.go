package cmd

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/sjeandeaux/toolators/information"
)

//CommandLine the I/O command line
type CommandLine struct {
	Stdout io.Writer
	Stderr io.Writer
	Stdin  *os.File
}

//Init log and print version
func (c *CommandLine) Init(prefix string) *CommandLine {
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin

	log.SetPrefix(prefix)
	log.SetOutput(c.Stderr)
	log.Println(information.Print())
	return c
}

//Fatal return 1 and log in Err
func (c *CommandLine) Fatal(err error) int {
	fmt.Fprintf(c.Stderr, fmt.Sprint(err))
	return 1
}

//Input return the reader on the data or file (can be http) or pipe
func (c *CommandLine) Input(data, file string) (io.ReadCloser, error) {
	return Input(data, file, c.Stdin)
}

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
