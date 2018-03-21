package main

import (
	"flag"
	"fmt"
	"io"
	"log"

	"github.com/sjeandeaux/github-cmd-go/semver"

	internalos "github.com/sjeandeaux/github-cmd-go/internal/os"
)

//commandLine the arguments command line
type commandLine struct {
	position string
	version  string
	stdout   io.Writer
	stderr   io.Writer
}

//increment the version
func (c *commandLine) increment() (*semver.Version, error) {
	println(c.version)
	v, err := semver.NewVersion(c.version)
	if err != nil {
		return nil, err
	}
	return v.Increment(c.position)
}

func (c *commandLine) init() {
	//flag
	log.SetPrefix("incrementator")
	log.SetOutput(c.stderr)

	//command line
	flag.StringVar(&c.position, "position", internalos.Getenv("INCREMENTOR_POSITION", "minor"), "The position major minor patch")
	flag.StringVar(&c.version, "version", "", "The version x.y.z use the git tag if not set")
	flag.Parse()
}

func (c *commandLine) main() int {
	value, err := c.increment()
	if err != nil {
		fmt.Fprintf(c.stderr, fmt.Sprint(err))
		return 1
	}
	fmt.Fprint(c.stdout, value)
	return 0
}