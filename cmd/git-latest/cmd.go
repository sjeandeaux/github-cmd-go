package main

import (
	"fmt"
	"io"
	"log"

	"github.com/sjeandeaux/github-cmd-go/semver"
)

type commandLine struct {
	stdout io.Writer
	stderr io.Writer
}

func (c *commandLine) init() {
	//flag
	log.SetPrefix("git-latest")
	log.SetOutput(c.stderr)
}

func (c *commandLine) main() int {
	value, err := semver.NewGitVersion()
	if err != nil {
		fmt.Fprintf(c.stderr, fmt.Sprint(err))
		return 1
	}
	fmt.Fprint(c.stdout, value)
	return 0
}
