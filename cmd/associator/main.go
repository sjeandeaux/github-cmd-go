package main

import (
	"flag"
	"log"
	"os"

	"github.com/sjeandeaux/github-cmd-go/github"
	"github.com/sjeandeaux/github-cmd-go/information"
)

//github token and path
type commandLineArgs struct {
	token       string
	owner       string
	repo        string
	create      bool
	file        string
	tag         string
	name        string
	label       string
	contentType string
}

var commandLineArgsValue = new(commandLineArgs)

func init() {
	flag.StringVar(&commandLineArgsValue.token, "token", os.Getenv("ASSOCIATOR_GITHUB_TOKEN"), "The token")
	flag.StringVar(&commandLineArgsValue.owner, "owner", os.Getenv("ASSOCIATOR_OWNER"), "The owner")
	flag.StringVar(&commandLineArgsValue.repo, "repo", os.Getenv("ASSOCIATOR_REPO"), "The repo")
	flag.StringVar(&commandLineArgsValue.tag, "tag", os.Getenv("ASSOCIATOR_TAG"), "The tag")
	flag.BoolVar(&commandLineArgsValue.create, "create", false, "Create tag")

	flag.StringVar(&commandLineArgsValue.file, "file", os.Getenv("ASSOCIATOR_FILE"), "The file")
	flag.StringVar(&commandLineArgsValue.name, "name", os.Getenv("ASSOCIATOR_NAME"), "The name")
	flag.StringVar(&commandLineArgsValue.label, "label", os.Getenv("ASSOCIATOR_LABEL"), "The label")
	flag.StringVar(&commandLineArgsValue.contentType, "content-type", os.Getenv("ASSOCIATOR_CONTENT_TYPE"), "The contentType")

	flag.Parse()
}

func main() {
	log.Println(information.Print())

	//Get the UploadURL
	client := github.NewClient(commandLineArgsValue.token, commandLineArgsValue.owner, commandLineArgsValue.repo)

	var only *github.Release
	var err error
	if commandLineArgsValue.create {

		e := &github.EditRelease{
			TagName: commandLineArgsValue.tag,
			Name:    commandLineArgsValue.tag,
		}

		only, err = client.CreateRelease(e)

		if err != nil {
			log.Fatal(err)
		}
	} else {
		only, err = client.GetReleaseByTag(commandLineArgsValue.tag)
		if err != nil {
			log.Fatal(err)
		}
	}

	a := &github.Asset{
		File:        commandLineArgsValue.file,
		Name:        commandLineArgsValue.name,
		Label:       commandLineArgsValue.label,
		ContentType: commandLineArgsValue.contentType,
	}

	if err := client.Upload(only.UploadURL(), a); err != nil {
		log.Fatal(err)
	}
}
