package main

import (
	"flag"
	"fmt"
	"time"

	internalcmd "github.com/sjeandeaux/toolators/internal/cmd"
)

type now func() time.Time

var timeNow now

func init() {
	timeNow = time.Now
}

//commandLine the arguments command line
type commandLine struct {
	internalcmd.CommandLine
	date string
}

func (c *commandLine) init() *commandLine {

	//flag
	c.Init("[d-day]")

	flag.StringVar(&c.date, "date", "2018-05-04", "date yyyy-dd-mm")

	flag.Parse()

	return c

}

func (c *commandLine) main() int {
	const hoursByDay = 24
	time1, err := time.Parse(time.RFC3339, c.date+"T00:00:00Z")
	if err != nil {
		return c.Fatal(err)
	}

	delta := int64(time1.Sub(timeNow()).Hours())
	days := (delta / hoursByDay)
	if delta%hoursByDay == 0 {
		fmt.Fprint(c.Stdout, days)
	} else {
		fmt.Fprint(c.Stdout, days+1)
	}
	return 0
}
