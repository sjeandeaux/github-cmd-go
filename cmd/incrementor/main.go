package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/sjeandeaux/github-cmd-go/semver"
)

type commandLine struct {
	kind    string
	version string
}

func (c *commandLine) increment() (*semver.Version, error) {
	v, err := semver.NewVersion(c.version)
	if err != nil {
		return nil, err
	}
	return v.IncrementString(c.kind)
}

var commandLineValue = new(commandLine)

func init() {
	flag.StringVar(&commandLineValue.kind, "kind", os.Getenv("INCREMENTOR_KIND"), "The kind major minor patch")
	flag.StringVar(&commandLineValue.version, "version", os.Getenv("INCREMENTOR_VERSION"), "The version x.y.z")
	flag.Parse()
}

func main() {
	if value, err := commandLineValue.increment(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Print(value)
	}
}
