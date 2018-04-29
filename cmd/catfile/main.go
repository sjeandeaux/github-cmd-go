package main

import (
	"os"
)

func main() {
	commandLineValue := commandLine{}
	commandLineValue.init()
	os.Exit(commandLineValue.main())
}
