package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type commandLine struct {
	kind    string
	version string
}

func (c *commandLine) convert(value string) (int64, error) {
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return -1, fmt.Errorf("%q is bad", c.version)
	}
	return parsed, nil
}

func (c *commandLine) decomposeVersion() (int64, int64, int64, error) {

	maMiPa := strings.Split(c.version, ".")
	if len(maMiPa) != 3 {
		return -1, -1, -1, fmt.Errorf("%q is bad", c.version)
	}
	major, err := c.convert(maMiPa[0])
	if err != nil {
		return -1, -1, -1, err
	}

	minor, err := c.convert(maMiPa[1])
	if err != nil {
		return -1, -1, -1, err
	}

	patch, err := c.convert(maMiPa[2])
	if err != nil {
		return -1, -1, -1, err
	}

	return major, minor, patch, nil

}

func (c *commandLine) increment() (string, error) {
	const format = "%d.%d.%d"
	major, minor, patch, err := c.decomposeVersion()
	if err != nil {
		return "", err
	}

	switch c.kind {
	case "major":
		return fmt.Sprintf(format, major+1, 0, 0), nil
	case "minor":
		return fmt.Sprintf(format, major, minor+1, 0), nil
	case "patch":
		return fmt.Sprintf(format, major, minor, patch+1), nil
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
