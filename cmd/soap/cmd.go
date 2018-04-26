package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	internalcmd "github.com/sjeandeaux/toolators/internal/cmd"
	internalhttp "github.com/sjeandeaux/toolators/internal/http"
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

	flag.StringVar(&c.action, "soap-action", "", "Action SOAP")
	flag.StringVar(&c.url, "soap-url", "", "URL SOAP")
	flag.StringVar(&c.data, "soap-data", "", "Data SOAP")
	flag.StringVar(&c.file, "soap-file", "", "File SOAP")
	flag.Parse()

}

func (c *commandLine) main() int {
	const (
		SOAPHeaderAction           = "SOAPAction"
		SOAPHeaderContentType      = "Content-Type"
		SOAPHeaderContentTypeValue = "text/xml;charset=UTF-8"
	)

	input, err := internalcmd.Input(c.data, c.file, c.stdin)
	if err != nil {
		fmt.Fprintf(c.stderr, fmt.Sprint(err))
		return -1
	}
	defer input.Close()

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
