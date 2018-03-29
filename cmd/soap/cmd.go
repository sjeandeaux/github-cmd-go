package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	internalhttp "github.com/sjeandeaux/github-cmd-go/internal/http"
)

//commandLine the arguments command line
type commandLine struct {
	client http.Client
	action string
	url    string
	data   string
	file   string

	stdout io.Writer
	stderr io.Writer
	stdin  *os.File
}

func (c *commandLine) init() {
	//flag
	log.SetPrefix("[soap]\t")
	log.SetOutput(c.stderr)

	flag.StringVar(&c.action, "action", "", "Action SOAP")
	flag.StringVar(&c.url, "url", "", "URL SOAP")
	flag.StringVar(&c.data, "data", "", "Data SOAP")
	flag.StringVar(&c.file, "file", "", "File SOAP")
	flag.Parse()

}

func (c *commandLine) input() (io.ReadCloser, error) {
	if c.data != "" {
		return ioutil.NopCloser(strings.NewReader(c.data)), nil
	}

	if c.file != "" {
		return os.Open(c.file)
	}

	fi, err := c.stdin.Stat()
	if err != nil {
		return nil, err
	}
	if fi.Mode()&os.ModeNamedPipe == 0 {
		return nil, errors.New("No pipe")
	}
	return os.Stdin, nil
}

func (c *commandLine) main() int {
	const (
		SOAPHeaderAction           = "SOAPAction"
		SOAPHeaderContentType      = "Content-Type"
		SOAPHeaderContentTypeValue = "text/xml;charset=UTF-8"
	)

	input, err := c.input()
	defer input.Close()
	if err != nil {
		fmt.Fprintf(c.stderr, fmt.Sprint(err))
		return -1
	}

	req, err := http.NewRequest(http.MethodPost, c.url, input)
	if err != nil {
		fmt.Fprintf(c.stderr, fmt.Sprint(err))
		return -1
	}
	req.Header.Set(SOAPHeaderContentType, SOAPHeaderContentTypeValue)
	req.Header.Set(SOAPHeaderAction, c.action)

	resp, err := c.client.Do(req)
	defer internalhttp.Close(resp)
	if err != nil {
		fmt.Fprintf(c.stderr, fmt.Sprint(err))
		return -1
	}

	_, err = io.Copy(c.stdout, resp.Body)
	if err != nil {
		fmt.Fprintf(c.stderr, fmt.Sprint(err))
		return -1
	}

	return resp.StatusCode
}
