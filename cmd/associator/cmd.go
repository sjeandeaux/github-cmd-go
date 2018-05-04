package main

import (
	"flag"
	"os"

	"github.com/sjeandeaux/toolators/github"

	internalcmd "github.com/sjeandeaux/toolators/internal/cmd"
)

//githubClient github interaction
type githubClient interface {
	CreateRelease(edit *github.EditRelease) (*github.Release, error)
	GetReleaseByTag(tag string) (*github.Release, error)
	Upload(urlPath string, u github.UploadInformation) error
}

//github token and path
type commandLine struct {
	internalcmd.CommandLine

	token       string
	owner       string
	repo        string
	create      bool
	path        string
	tag         string
	name        string
	label       string
	contentType string

	githubClient githubClient
}

func (c *commandLine) init() *commandLine {
	c.Init("[associator]")

	flag.StringVar(&c.token, "token", os.Getenv("GITHUB_TOKEN"), "The token")
	flag.StringVar(&c.owner, "owner", "", "The owner")
	flag.StringVar(&c.repo, "repo", "", "The repo")
	flag.StringVar(&c.tag, "tag", "", "The tag")
	flag.BoolVar(&c.create, "create", false, "Create tag")

	flag.StringVar(&c.path, "path", "", "The path")
	flag.StringVar(&c.name, "name", "", "The name")
	flag.StringVar(&c.label, "label", "", "The label")
	flag.StringVar(&c.contentType, "content-type", "", "The contentType")

	flag.Parse()

	c.githubClient = github.NewClient(c.token, c.owner, c.repo)
	return c
}

func (c *commandLine) main() int {

	var only *github.Release
	var err error
	if c.create {

		e := &github.EditRelease{
			TagName: c.tag,
			Name:    c.tag,
		}

		only, err = c.githubClient.CreateRelease(e)

		if err != nil {
			return c.Fatal(err)
		}
	} else {
		only, err = c.githubClient.GetReleaseByTag(c.tag)
		if err != nil {
			return c.Fatal(err)
		}
	}

	a := &github.Asset{
		File:        c.path,
		Name:        c.name,
		Label:       c.label,
		ContentType: c.contentType,
	}

	if err := c.githubClient.Upload(only.UploadURL(), a); err != nil {
		return c.Fatal(err)
	}
	return 0
}
