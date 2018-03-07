package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/sjeandeaux/github-cmd-go/semver"

	internalos "github.com/sjeandeaux/github-cmd-go/internal/os"
)

type commandLine struct {
	position string
	version  string
}

func (c *commandLine) increment() (*semver.Version, error) {
	v, err := semver.NewVersion(c.version)
	if err != nil {
		return nil, err
	}
	return v.Increment(c.position)
}

var commandLineValue = new(commandLine)

func init() {
	flag.StringVar(&commandLineValue.position, "position", internalos.Getenv("INCREMENTOR_POSITION", "minor"), "The position major minor patch")
	flag.StringVar(&commandLineValue.version, "version", internalos.Getenv("INCREMENTOR_VERSION", "0.1.0"), "The version x.y.z")
	flag.Parse()
}

func main() {
	if value, err := commandLineValue.increment(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Print(value)
	}
}
