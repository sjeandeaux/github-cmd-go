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

func (c *commandLine) increment() (string, error) {
	const format = "%d.%d.%d"
	version, err := semver.NewVersion(c.version)
	if err != nil {
		return "", err
	}

	switch c.kind {
	case "major":
		return fmt.Sprintf(format, version.Major+1, 0, 0), nil
	case "minor":
		return fmt.Sprintf(format, version.Major, version.Minor+1, 0), nil
	case "patch":
		return fmt.Sprintf(format, version.Major, version.Minor, version.Patch+1), nil
	default:
		return "", fmt.Errorf("%q is unknown", c.kind)
	}

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
