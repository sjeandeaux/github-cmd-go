//Package github play with github.
package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
)

//Client the http connection
type Client struct {
	httpClient *http.Client
	token      string
	owner      string
	repo       string
}

//NewClient init the client
func NewClient(token, owner, repo string) *Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	//Get the UploadURL
	oauthClient := oauth2.NewClient(oauth2.NoContext, ts)
	return &Client{
		httpClient: oauthClient,
		token:      token,
		owner:      owner,
		repo:       repo,
	}

}

//Release we want only the upload url.
type Release struct {
	UploadURLTemplate string `json:"upload_url"`
	TagName           string `json:"tag_name"`
	URL               string `json:"url"`
}

//UploadURL the upload url for tag
func (o *Release) UploadURL() string {
	uriTempl := strings.Index(o.UploadURLTemplate, "{?")
	if uriTempl >= 0 {
		return o.UploadURLTemplate[:uriTempl]
	}
	return o.UploadURLTemplate
}

//EditRelease information to send to release edition https://developer.github.com/v3/repos/releases/#edit-a-release
type EditRelease struct {
	TagName         string `json:"tag_name"`
	Name            string `json:"name,omitempty"`
	TargetCommitish string `json:"target_commitish,omitempty"`
	Body            string `json:"body,omitempty"`
	Draft           bool   `json:"draft,omitempty"`
	Prerelease      bool   `json:"prerelease,omitempty"`
}

//Upload on urlPath
func (c *Client) Upload(urlPath string, a *Asset) error {
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
	resp, err := c.httpClient.Do(request)
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

//GetReleaseByTag get the information
func (c *Client) GetReleaseByTag(tag string) (*Release, error) {
	const githubAPI = "https://api.github.com/repos/"
	const (
		contentType     = "Content-Type"
		applicationJSON = "application/json"
	)

	url := fmt.Sprint(githubAPI, c.owner, "/", c.repo, "/releases/tags/", tag)

	request, _ := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Add(contentType, applicationJSON)
	resp, err := c.httpClient.Do(request)
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
	onlyURL := &Release{}
	err = json.NewDecoder(resp.Body).Decode(onlyURL)
	if err != nil {
		return nil, err
	}

	return onlyURL, nil
}

//CreateRelease create a release
func (c *Client) CreateRelease(tag string) (*Release, error) {
	const githubAPI = "https://api.github.com/repos/"
	const (
		contentType     = "Content-Type"
		applicationJSON = "application/json"
	)

	e := &EditRelease{
		TagName: tag,
		Name:    tag,
	}

	url := fmt.Sprint(githubAPI, c.owner, "/", c.repo, "/releases")
	jsonValue, _ := json.Marshal(e)
	request, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonValue))
	request.Header.Add(contentType, applicationJSON)
	resp, err := c.httpClient.Do(request)
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
	onlyURL := &Release{}
	err = json.NewDecoder(resp.Body).Decode(onlyURL)
	if err != nil {
		return nil, err
	}

	return onlyURL, nil
}
