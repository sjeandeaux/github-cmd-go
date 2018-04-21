package main

import (
	"flag"
	"io"
	"log"
	"os"
)

//commandLine the arguments command line
type commandLine struct {
	stdout io.Writer
	stderr io.Writer
	stdin  *os.File
}

func (c *commandLine) init() {

	//flag
	log.SetPrefix("[aws]\t")
	log.SetOutput(c.stderr)
	flag.Parse()

}

func (c *commandLine) main() int {

	return 0
}
