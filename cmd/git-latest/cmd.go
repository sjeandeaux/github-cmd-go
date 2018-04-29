package main

import (
	"fmt"

	"github.com/sjeandeaux/toolators/semver"

	internalcmd "github.com/sjeandeaux/toolators/internal/cmd"
)

type commandLine struct {
	internalcmd.CommandLine
	gitVersion func() (*semver.Version, error)
}

func (c *commandLine) init() {
	//flag
	c.Init("[git-latest]")
	c.gitVersion = semver.NewGitVersion
}

func (c *commandLine) main() int {
	value, err := c.gitVersion()
	if err != nil {
		return c.Fatal(err)
	}
	fmt.Fprint(c.Stdout, value)
	return 0
}
