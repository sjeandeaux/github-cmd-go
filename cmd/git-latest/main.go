package main

import (
	"fmt"
	"log"

	"github.com/sjeandeaux/github-cmd-go/semver"
)

type commandLine struct {
}

func (c *commandLine) init() {
	//flag
	log.SetPrefix("git-latest")
}

func (c *commandLine) main() {
	if value, err := semver.NewGitVersion(); err != nil {
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
