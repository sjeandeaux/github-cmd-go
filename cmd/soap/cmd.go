package main

import (
	"flag"
	"io"
	"net/http"

	internalcmd "github.com/sjeandeaux/toolators/internal/cmd"
	internalhttp "github.com/sjeandeaux/toolators/internal/http"
)

//commandLine the arguments command line
type commandLine struct {
	internalcmd.CommandLine
	httpClient *http.Client
	action     string
	url        string
	data       string
	file       string
}

func (c *commandLine) init() *commandLine {
	//flag
	c.Init("[soap]")

	flag.StringVar(&c.action, "action", "", "Action SOAP")
	flag.StringVar(&c.url, "url", "", "URL SOAP")
	flag.StringVar(&c.data, "data", "", "Data SOAP")
	flag.StringVar(&c.file, "file", "", "File SOAP")
	flag.Parse()

	return c

}

func (c *commandLine) main() int {
	const (
		SOAPHeaderAction           = "SOAPAction"
		SOAPHeaderContentType      = "Content-Type"
		SOAPHeaderContentTypeValue = "text/xml;charset=UTF-8"
	)

	input, err := c.Input(c.data, c.file)
	if err != nil {
		return c.Fatal(err)
	}
	defer input.Close()

	req, err := http.NewRequest(http.MethodPost, c.url, input)
	if err != nil {
		return c.Fatal(err)
	}
	req.Header.Set(SOAPHeaderContentType, SOAPHeaderContentTypeValue)
	req.Header.Set(SOAPHeaderAction, c.action)

	resp, err := c.httpClient.Do(req)
	defer internalhttp.Close(resp)
	if err != nil {
		return c.Fatal(err)
	}

	_, err = io.Copy(c.Stdout, resp.Body)
	if err != nil {
		return c.Fatal(err)
	}

	return resp.StatusCode
}
