package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"io"

	internalcmd "github.com/sjeandeaux/toolators/internal/cmd"
)

type commandLine struct {
	internalcmd.CommandLine

	data string
	file string
}

func cat(r io.Reader, w io.Writer) error {

	const (
		header  = "\x1b]1337;File=;inline=1:"
		footer  = "\a"
		termEnv = "TERM"
	)

	pr, pw := io.Pipe()
	go encode(r, pw)

	bHeader := bytes.NewBufferString(header)
	bFooter := bytes.NewBufferString(footer)

	_, err := io.Copy(w, io.MultiReader(bHeader, pr, bFooter))
	return err
}

func encode(r io.Reader, w *io.PipeWriter) {
	encoder := base64.NewEncoder(base64.StdEncoding, w)
	defer func() {
		if err := encoder.Close(); err != nil {
			w.CloseWithError(err)
		} else {
			w.Close()
		}
	}()
	_, err := io.Copy(encoder, r)
	w.CloseWithError(err)
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
