package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/sjeandeaux/github-cmd-go/information"
	"golang.org/x/oauth2"
)

//onlyUploadURL we want only the upload url.
type onlyUploadURL struct {
	UploadURLTemplate string `json:"upload_url"`
	TagName           string `json:"tag_name"`
	URL               string `json:"url"`
}

func (o *onlyUploadURL) uploadURL() string {
	uriTempl := strings.Index(o.UploadURLTemplate, "{?")
	if uriTempl >= 0 {
		return o.UploadURLTemplate[:uriTempl]
	}
	return o.UploadURLTemplate
}

//github token and path
type github struct {
	//token to communicate with github. https://github.com/settings/tokens/new
	token  string
	owner  string
	repo   string
	tag    string
	create bool
}

type assertValue asset

//asset information on asset to upload https://developer.github.com/v3/repos/releases/#upload-a-release-asset
type asset struct {
	file        string
	name        string
	label       string
	contentType string
}

type editRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
}

func (a *asset) reader() (io.ReadCloser, error) {
	return os.Open(a.file)
}

func (a *asset) size() (int64, error) {
	fileStat, err := os.Stat(a.file)
	if err != nil {
		return -1, err
	}
	return fileStat.Size(), nil
}

func (a *asset) request(urlPath string, body io.Reader, size int64) (*http.Request, error) {
	const (
		contentType = "Content-Type"
	)
	request, _ := http.NewRequest(http.MethodPost, urlPath, body)
	request.ContentLength = size
	query := request.URL.Query()
	query.Add("name", a.name)
	query.Add("label", a.label)
	request.URL.RawQuery = query.Encode()
	request.Header.Add(contentType, a.contentType)
	return request, nil
}

func (a *asset) upload(client *http.Client, urlPath string) error {
	file, err := a.reader()
	defer file.Close()
	if err != nil {
		return err
	}

	size, err := a.size()
	if err != nil {
		return err
	}

	request, err := a.request(urlPath, file, size)
	if err != nil {
		return err
	}
	resp, err := client.Do(request)
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()
	if err != nil {
		return err
	}
	if resp.StatusCode != 201 {
		b, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("failed %d %s", resp.StatusCode, b)
	}
	return nil
}

var githubValue = new(github)
var assetValue = new(asset)

func init() {
	flag.StringVar(&githubValue.token, "token", os.Getenv("ASSOCIATOR_GITHUB_TOKEN"), "The token")
	flag.StringVar(&githubValue.owner, "owner", os.Getenv("ASSOCIATOR_OWNER"), "The owner")
	flag.StringVar(&githubValue.repo, "repo", os.Getenv("ASSOCIATOR_REPO"), "The repo")
	flag.StringVar(&githubValue.tag, "tag", os.Getenv("ASSOCIATOR_TAG"), "The tag")
	flag.BoolVar(&githubValue.create, "create", false, "Create tag")

	flag.StringVar(&assetValue.file, "file", os.Getenv("ASSOCIATOR_FILE"), "The file")
	flag.StringVar(&assetValue.name, "name", os.Getenv("ASSOCIATOR_NAME"), "The name")
	flag.StringVar(&assetValue.label, "label", os.Getenv("ASSOCIATOR_LABEL"), "The label")
	flag.StringVar(&assetValue.contentType, "content-type", os.Getenv("ASSOCIATOR_CONTENT_TYPE"), "The contentType")

	flag.Parse()
}

func main() {
	log.Println(information.Print())

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubValue.token},
	)

	//Get the UploadURL
	client := oauth2.NewClient(oauth2.NoContext, ts)

	var only *onlyUploadURL
	var err error
	if githubValue.create {
		only, err = githubValue.createRelease(client)

		if err != nil {
			log.Fatal(err)
		}
	} else {
		only, err = githubValue.getReleaseByTag(client)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := assetValue.upload(client, only.uploadURL()); err != nil {
		log.Fatal(err)
	}

}

func (g *github) createRelease(client *http.Client) (*onlyUploadURL, error) {
	const githubAPI = "https://api.github.com/repos/"
	const (
		contentType     = "Content-Type"
		applicationJSON = "application/json"
	)

	e := &editRelease{
		TagName: g.tag,
		Name:    g.tag,
	}

	url := fmt.Sprint(githubAPI, g.owner, "/", g.repo, "/releases")
	jsonValue, _ := json.Marshal(e)
	request, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonValue))
	request.Header.Add(contentType, applicationJSON)
	resp, err := client.Do(request)
	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 201 {
		b, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed %d %s", resp.StatusCode, b)
	}
	onlyURL := &onlyUploadURL{}
	err = json.NewDecoder(resp.Body).Decode(onlyURL)
	if err != nil {
		return nil, err
	}

	return onlyURL, nil
}

func (g *github) getReleaseByTag(client *http.Client) (*onlyUploadURL, error) {
	const githubAPI = "https://api.github.com/repos/"
	const (
		contentType     = "Content-Type"
		applicationJSON = "application/json"
	)

	url := fmt.Sprint(githubAPI, g.owner, "/", g.repo, "/releases/tags/", g.tag)

	request, _ := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Add(contentType, applicationJSON)
	resp, err := client.Do(request)
	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		b, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed %d %s", resp.StatusCode, b)
	}
	onlyURL := &onlyUploadURL{}
	err = json.NewDecoder(resp.Body).Decode(onlyURL)
	if err != nil {
		return nil, err
	}

	return onlyURL, nil
}
