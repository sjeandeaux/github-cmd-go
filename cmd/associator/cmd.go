package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/sjeandeaux/github-cmd-go/github"
	"github.com/sjeandeaux/github-cmd-go/information"
)

//github token and path
type commandLine struct {
	token       string
	owner       string
	repo        string
	create      bool
	file        string
	tag         string
	name        string
	label       string
	contentType string

	stdout io.Writer
	stderr io.Writer
}

func (c *commandLine) init() {
	log.SetPrefix("associator")
	log.SetOutput(c.stderr)

	flag.StringVar(&c.token, "token", os.Getenv("ASSOCIATOR_GITHUB_TOKEN"), "The token")
	flag.StringVar(&c.owner, "owner", os.Getenv("ASSOCIATOR_OWNER"), "The owner")
	flag.StringVar(&c.repo, "repo", os.Getenv("ASSOCIATOR_REPO"), "The repo")
	flag.StringVar(&c.tag, "tag", os.Getenv("ASSOCIATOR_TAG"), "The tag")
	flag.BoolVar(&c.create, "create", false, "Create tag")

	flag.StringVar(&c.file, "file", os.Getenv("ASSOCIATOR_FILE"), "The file")
	flag.StringVar(&c.name, "name", os.Getenv("ASSOCIATOR_NAME"), "The name")
	flag.StringVar(&c.label, "label", os.Getenv("ASSOCIATOR_LABEL"), "The label")
	flag.StringVar(&c.contentType, "content-type", os.Getenv("ASSOCIATOR_CONTENT_TYPE"), "The contentType")

	flag.Parse()
}

func (c *commandLine) main() int {
	log.Println(information.Print())

	//Get the UploadURL
	client := github.NewClient(c.token, c.owner, c.repo)

	var only *github.Release
	var err error
	if c.create {

		e := &github.EditRelease{
			TagName: c.tag,
			Name:    c.tag,
		}

		only, err = client.CreateRelease(e)

		if err != nil {
			fmt.Fprintf(c.stderr, fmt.Sprint(err))
			return 1
		}
	} else {
		only, err = client.GetReleaseByTag(c.tag)
		if err != nil {
			fmt.Fprintf(c.stderr, fmt.Sprint(err))
			return 1
		}
	}

	a := &github.Asset{
		File:        c.file,
		Name:        c.name,
		Label:       c.label,
		ContentType: c.contentType,
	}

	if err := client.Upload(only.UploadURL(), a); err != nil {
		fmt.Fprintf(c.stderr, fmt.Sprint(err))
		return 1
	}
	return 0
}
