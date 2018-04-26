package main

import (
	"os"
)

func main() {
	var commandLineValue = commandLine{
		stdout: os.Stdout,
		stderr: os.Stderr,
	}
	commandLineValue.init()
	os.Exit(commandLineValue.main())
}
