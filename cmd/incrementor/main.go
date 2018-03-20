package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/sjeandeaux/github-cmd-go/semver"

	internalos "github.com/sjeandeaux/github-cmd-go/internal/os"
)

//commandLine the arguments command line
type commandLine struct {
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
	log.SetPrefix("incrementator")

	//command line
	flag.StringVar(&commandLineValue.position, "position", internalos.Getenv("INCREMENTOR_POSITION", "minor"), "The position major minor patch")
	flag.StringVar(&commandLineValue.version, "version", "", "The version x.y.z use the git tag if not set")
	flag.Parse()
}

func (c *commandLine) main() {
	if value, err := commandLineValue.increment(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Print(value)
	}
}

var commandLineValue = new(commandLine)

func init() {
	commandLineValue.init()
}

func main() {
	commandLineValue.main()
}
