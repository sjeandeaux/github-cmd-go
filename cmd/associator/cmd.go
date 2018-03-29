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

//githubClient github interaction
type githubClient interface {
	CreateRelease(edit *github.EditRelease) (*github.Release, error)
	GetReleaseByTag(tag string) (*github.Release, error)
	Upload(urlPath string, u github.UploadInformation) error
}

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

	githubClient githubClient
	stdout       io.Writer
	stderr       io.Writer
}

func (c *commandLine) init() {
	log.SetPrefix("[associator]\t")
	log.SetOutput(c.stderr)

	flag.StringVar(&c.token, "token", os.Getenv("GITHUB_TOKEN"), "The token")
	flag.StringVar(&c.owner, "owner", "", "The owner")
	flag.StringVar(&c.repo, "repo", "", "The repo")
	flag.StringVar(&c.tag, "tag", "", "The tag")
	flag.BoolVar(&c.create, "create", false, "Create tag")

	flag.StringVar(&c.file, "file", "", "The file")
	flag.StringVar(&c.name, "name", "", "The name")
	flag.StringVar(&c.label, "label", "", "The label")
	flag.StringVar(&c.contentType, "content-type", "", "The contentType")

	flag.Parse()

	c.githubClient = github.NewClient(c.token, c.owner, c.repo)
}

func (c *commandLine) main() int {
	log.Println(information.Print())

	var only *github.Release
	var err error
	if c.create {

		e := &github.EditRelease{
			TagName: c.tag,
			Name:    c.tag,
		}

		only, err = c.githubClient.CreateRelease(e)

		if err != nil {
			fmt.Fprintf(c.stderr, fmt.Sprint(err))
			return 1
		}
	} else {
		only, err = c.githubClient.GetReleaseByTag(c.tag)
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

	if err := c.githubClient.Upload(only.UploadURL(), a); err != nil {
		fmt.Fprintf(c.stderr, fmt.Sprint(err))
		return 1
	}
	return 0
}
