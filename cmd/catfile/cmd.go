package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

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
		firstLine = "\033]1337;File=;inline=1:"
		lastLine  = "\a"

		tmuxEnv       = "TMUX"
		tmuxFirstLine = "\033Ptmux;\033"
		tmuxLastLine  = "\033\\"
	)
	tmux := os.Getenv(tmuxEnv) != ""

	if tmux {
		fmt.Fprint(w, tmuxFirstLine)
	}
	fmt.Fprint(w, firstLine)

	encoder := base64.NewEncoder(base64.StdEncoding, w)
	defer encoder.Close()
	_, err := io.Copy(encoder, r)
	if err != nil {
		return err
	}

	fmt.Fprintln(w, lastLine)
	if tmux {
		fmt.Fprintln(w, tmuxLastLine)
	}
	return nil
}

func (c *commandLine) init() {
	//flag
	log.SetPrefix("[git]\t")
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
