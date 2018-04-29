package main

import (
	"flag"
	"fmt"

	"github.com/sjeandeaux/toolators/semver"

	internalcmd "github.com/sjeandeaux/toolators/internal/cmd"
	internalos "github.com/sjeandeaux/toolators/internal/os"
)

//commandLine the arguments command line
type commandLine struct {
	internalcmd.CommandLine
	position string
	version  string
}

//increment the version
func (c *commandLine) increment() (*semver.Version, error) {
	v, err := semver.NewVersion(c.version)
	if err != nil {
		return nil, err
	}
	return v.Increment(c.position)
}

func (c *commandLine) init() {
	//flag
	c.Init("[incrementator]")

	//command line
	flag.StringVar(&c.position, "position", internalos.Getenv("INCREMENTOR_POSITION", "minor"), "The position major minor patch")
	flag.StringVar(&c.version, "version", "", "The version x.y.z use the git tag if not set")
	flag.Parse()
}

func (c *commandLine) main() int {
	value, err := c.increment()
	if err != nil {
		return c.Fatal(err)
	}
	fmt.Fprint(c.Stdout, value)
	return 0
}
