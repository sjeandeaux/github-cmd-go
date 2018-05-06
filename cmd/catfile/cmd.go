package main

import (
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	internalcmd "github.com/sjeandeaux/toolators/internal/cmd"
)

type commandLine struct {
	internalcmd.CommandLine

	data string
	file string
}

func cat(r io.Reader, w io.Writer) error {

	const (
		header = "\x1b]1337;File=;inline=1:"
		footer = '\a'

		termEnv = "TERM"
	)

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

func (c *commandLine) init() *commandLine {
	//flag
	c.Init("[catfile]")
	flag.StringVar(&c.data, "data", "", "Data")
	flag.StringVar(&c.file, "file", "", "File")
	flag.Parse()
	return c
}

func (c *commandLine) main() int {
	r, err := internalcmd.Input(c.data, c.file, c.Stdin)
	if err != nil {
		return c.Fatal(err)
	}
	err = cat(r, c.Stdout)
	if err != nil {
		return c.Fatal(err)
	}
	return 0
}
