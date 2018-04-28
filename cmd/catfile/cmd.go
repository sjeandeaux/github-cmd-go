package main

import (
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	internalcmd "github.com/sjeandeaux/toolators/internal/cmd"
)

type commandLine struct {
	stdout io.Writer
	stderr io.Writer
	stdin  *os.File

	process string
	data    string
	file    string
}

func cat(r io.Reader, w io.Writer) error {

	const (
		header = "\x1b]1337;File=;inline=1:"
		footer = '\a'

		termEnv = "TERM"
	)
	//catfile -file https://media.giphy.com/media/OSQoA7hdnoIuQchEBS/giphy.gif
	screen := strings.HasPrefix(os.Getenv(termEnv), "screen")

	if screen {
		return errors.New("sorry for tmux")
	}
	fmt.Fprint(w, header)

	encoder := base64.NewEncoder(base64.StdEncoding, w)
	defer encoder.Close()
	_, err := io.Copy(encoder, r)
	if err != nil {
		return err
	}
	fmt.Fprint(w, footer)

	return nil
}

func (c *commandLine) init() {
	//flag
	log.SetPrefix("[catfile]\t")
	log.SetOutput(c.stderr)

	flag.StringVar(&c.data, "data", "", "Data")
	flag.StringVar(&c.file, "file", "", "File")
	flag.Parse()
}

func (c *commandLine) main() int {
	r, err := internalcmd.Input(c.data, c.file, c.stdin)
	if err != nil {
		fmt.Fprintf(c.stderr, fmt.Sprint(err))
		return 1
	}
	err = cat(r, c.stdout)
	if err != nil {
		fmt.Fprintf(c.stderr, fmt.Sprint(err))
		return 1
	}
	return 0
}
