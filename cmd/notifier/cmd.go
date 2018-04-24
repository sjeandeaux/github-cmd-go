package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/sjeandeaux/toolators/notification/hipchat"

	internalos "github.com/sjeandeaux/toolators/internal/os"
)

type commandLine struct {
	stdout io.Writer
	stderr io.Writer
	stdin  *os.File

	token    string
	hostname string
	room     string
	data     string
	file     string
}

func (c *commandLine) init() {
	//flag
	log.SetPrefix("[notifier]\t")
	log.SetOutput(c.stderr)

	flag.StringVar(&c.token, "hipchat-token", internalos.Getenv("hipchat_token", ""), "token")
	flag.StringVar(&c.hostname, "hipchat-hostname", internalos.Getenv("hipchat_hostname", ""), "hostname")
	flag.StringVar(&c.room, "hipchat-room", internalos.Getenv("hipchat_room", ""), "room")
	flag.StringVar(&c.data, "data", "", "Data Message")
	flag.StringVar(&c.file, "file", "", "File Message")
	flag.Parse()

}

// TODO avoid the copy/paster
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
	data, err := c.input()
	if err != nil {
		fmt.Fprintf(c.stderr, fmt.Sprint(err))
		return 1
	}
	defer data.Close()

	notifier := hipchat.NewNotifier(fmt.Sprintf(hipchat.URLRoom, c.hostname, c.room), c.token)
	err = notifier.Send(data)
	if err != nil {
		fmt.Fprintf(c.stderr, fmt.Sprint(err))
		return 1
	}
	return 0
}
