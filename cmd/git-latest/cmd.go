package main

import (
	"fmt"
	"io"
	"log"

	"github.com/sjeandeaux/github-cmd-go/semver"
)

type commandLine struct {
	stdout     io.Writer
	stderr     io.Writer
	gitVersion func() (*semver.Version, error)
}

func (c *commandLine) init() {
	//flag
	log.SetPrefix("git-latest")
	log.SetOutput(c.stderr)
	c.gitVersion = semver.NewGitVersion
}

func (c *commandLine) main() int {
	value, err := c.gitVersion()
	if err != nil {
		fmt.Fprintf(c.stderr, fmt.Sprint(err))
		return 1
	}
	fmt.Fprint(c.stdout, value)
	return 0
}
