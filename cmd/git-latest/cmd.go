package main

import (
	"fmt"
	"io"
	"log"

	"github.com/sjeandeaux/github-cmd-go/semver"
)

//GitVersion get the git version
type gitVersion interface {
	GetCurrentVersion() (*semver.Version, error)
}

type defaultGitVersion struct{}

//GetCurrentVersion get the current version git
func (*defaultGitVersion) GetCurrentVersion() (*semver.Version, error) {
	return semver.NewGitVersion()
}

type commandLine struct {
	stdout     io.Writer
	stderr     io.Writer
	gitVersion gitVersion
}

func (c *commandLine) init() {
	//flag
	log.SetPrefix("git-latest")
	log.SetOutput(c.stderr)
}

func (c *commandLine) main() int {
	value, err := c.gitVersion.GetCurrentVersion()
	if err != nil {
		fmt.Fprintf(c.stderr, fmt.Sprint(err))
		return 1
	}
	fmt.Fprint(c.stdout, value)
	return 0
}
