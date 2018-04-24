package main

import (
	"os"
)

var commandLineValue = commandLine{
	stdout: os.Stdout,
	stderr: os.Stderr,
	stdin:  os.Stdin,
}

func init() {
	commandLineValue.init()
}

func main() {
	os.Exit(commandLineValue.main())
}
