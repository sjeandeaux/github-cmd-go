package main

import (
	"fmt"
	"log"

	"github.com/sjeandeaux/github-cmd-go/semver"
)

func main() {
	if value, err := semver.NewGitVersion(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Print(value)
	}
}
