package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/sjeandeaux/toolators/notification/hipchat"

	internalcmd "github.com/sjeandeaux/toolators/internal/cmd"
	internalos "github.com/sjeandeaux/toolators/internal/os"
)

type commandLine struct {
	internalcmd.CommandLine

	token    string
	hostname string
	room     string
	data     string
	file     string

	message string
	format  string
	from    string
	notify  bool
}

func (c *commandLine) init() {
	//flag
	c.Init("[notifier]")

	flag.StringVar(&c.data, "data", "", "Data Message")
	flag.StringVar(&c.file, "file", "", "File Message")

	c.flagHipChat()

	flag.Parse()

}

func (c *commandLine) flagHipChat() {
	flag.StringVar(&c.token, "hipchat-token", internalos.Getenv("hipchat_token", ""), "token")
	flag.StringVar(&c.hostname, "hipchat-hostname", internalos.Getenv("hipchat_hostname", ""), "hostname")
	flag.StringVar(&c.room, "hipchat-room", internalos.Getenv("hipchat_room", ""), "room")

	flag.StringVar(&c.from, "hipchat-from", "notifier", "from")
	flag.BoolVar(&c.notify, "hipchat-notify", true, "notifiy")
	flag.StringVar(&c.message, "hipchat-message", "", "message")
	flag.StringVar(&c.format, "hipchat-format", "text", "message")
}

func (c *commandLine) main() int {
	data, err := c.Input(c.data, c.file)

	//TODO ugly have to change this part
	if internalcmd.IsNoData(err) {
		jsonMap := map[string]interface{}{
			"message":        c.message,
			"notify":         c.notify,
			"from":           c.from,
			"message_format": c.format,
		}
		b, err := json.Marshal(jsonMap)
		if err != nil {
			return c.Fatal(err)
		}
		data = ioutil.NopCloser(bytes.NewReader(b))
	} else {
		return c.Fatal(err)
	}
	defer data.Close()

	notifier := hipchat.NewNotifier(fmt.Sprintf(hipchat.URLRoom, c.hostname, c.room), c.token)
	err = notifier.Send(data)
	if err != nil {
		return c.Fatal(err)
	}
	return 0
}
