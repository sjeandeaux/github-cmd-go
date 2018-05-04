package main

import (
	"os"
)

func main() {
	os.Exit((&commandLine{}).init().main())
}
