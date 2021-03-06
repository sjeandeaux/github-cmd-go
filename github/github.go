//Package github play with github.
package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	internalhttp "github.com/sjeandeaux/toolators/internal/http"
	"golang.org/x/oauth2"
)

const (
	githubAPI = "https://api.github.com/repos"
)

//Release we want only the upload url.
type Release struct {
	//UploadURLTemplate the upload url template
	UploadURLTemplate string `json:"upload_url"`
	//TagName the tag name
	TagName string `json:"tag_name"`
	//URL the URL
	URL string `json:"url"`
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

//Client the http connection
type Client struct {
	httpClient *http.Client
	owner      string
	repo       string
	baseURL    string
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
		owner:      owner,
		repo:       repo,
		baseURL:    githubAPI,
	}

}

//UploadURL the upload url for tag
func (o *Release) UploadURL() string {
	uriTempl := strings.Index(o.UploadURLTemplate, "{?")
	if uriTempl >= 0 {
		return o.UploadURLTemplate[:uriTempl]
	}
	return o.UploadURLTemplate
}

//UploadInformation we need file and size
type UploadInformation interface {
	reader() (io.ReadCloser, error)
	size() (int64, error)
	headers() map[string]string
	parameters() map[string]string
}

//Upload on urlPath
func (c *Client) Upload(urlPath string, u UploadInformation) error {
	body, err := u.reader()
	defer body.Close()
	if err != nil {
		return err
	}

	size, err := u.size()
	if err != nil {
		return err
	}

	//body and size
	request, _ := http.NewRequest(http.MethodPost, urlPath, body)
	request.ContentLength = size
	//header
	for k, v := range u.headers() {
		request.Header.Add(k, v)
	}

	//query
	query := request.URL.Query()
	for k, v := range u.parameters() {
		query.Add(k, v)
	}
	request.URL.RawQuery = query.Encode()

	resp, err := c.httpClient.Do(request)
	defer internalhttp.Close(resp)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated {
		b, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("failed %d %s", resp.StatusCode, b)
	}
	return nil
}

//GetReleaseByTag get the information
func (c *Client) GetReleaseByTag(tag string) (*Release, error) {
	url := fmt.Sprint(c.baseURL, "/", c.owner, "/", c.repo, "/releases/tags/", tag)
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Add(internalhttp.ContentType, internalhttp.ApplicationJSON)
	resp, err := c.httpClient.Do(request)
	defer internalhttp.Close(resp)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed %d %s", resp.StatusCode, b)
	}

	onlyURL := &Release{}
	return onlyURL, decode(resp.Body, onlyURL)
}

//CreateRelease create a release
func (c *Client) CreateRelease(edit *EditRelease) (*Release, error) {
	url := fmt.Sprint(c.baseURL, "/", c.owner, "/", c.repo, "/releases")
	jsonValue, _ := json.Marshal(edit)
	request, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonValue))
	request.Header.Add(internalhttp.ContentType, internalhttp.ApplicationJSON)
	resp, err := c.httpClient.Do(request)
	defer internalhttp.Close(resp)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusCreated {
		b, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed %d %s", resp.StatusCode, b)
	}
	onlyURL := &Release{}
	return onlyURL, decode(resp.Body, onlyURL)
}

//decode decode the reader in i
func decode(r io.Reader, i interface{}) error {
	return json.NewDecoder(r).Decode(i)
}
