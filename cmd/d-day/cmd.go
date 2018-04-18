package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

//commandLine the arguments command line
type commandLine struct {
	date string

	stdout io.Writer
	stderr io.Writer
	stdin  *os.File
}

func (c *commandLine) init() {
	//flag
	log.SetPrefix("[d-day]\t")
	log.SetOutput(c.stderr)

	flag.StringVar(&c.date, "date", "2018-05-04", "date yyyy-dd-mm")

	flag.Parse()

}

func (c *commandLine) main() int {
	time1, err := time.Parse(time.RFC3339, c.date+"T00:00:00Z")
	if err != nil {
		fmt.Fprintf(c.stderr, fmt.Sprint(err))
		return 1
	}

	delta := time1.Sub(time.Now())
	days := int64(delta.Hours() / 24)
	fmt.Fprint(c.stdout, days)
	return 0
}
