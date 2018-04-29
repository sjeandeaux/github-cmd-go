package main

import (
	"os"
)

func main() {
	var commandLineValue = commandLine{}
	commandLineValue.init()
	os.Exit(commandLineValue.main())
}
