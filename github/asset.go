package github

import (
	"io"
	"net/http"
	"os"
)

//Asset information on asset to upload https://developer.github.com/v3/repos/releases/#upload-a-release-asset
type Asset struct {
	//File the file of OS
	File string
	//Name paramter for github
	Name string
	//Label parameter for github
	Label string
	//ContentType value for HEADER
	ContentType string
}

//reader reader of asset
func (a *Asset) reader() (io.ReadCloser, error) {
	return os.Open(a.File)
}

//size size of asset
func (a *Asset) size() (int64, error) {
	fileStat, err := os.Stat(a.File)
	if err != nil {
		return -1, err
	}
	return fileStat.Size(), nil
}

//request the request to upload asset on github
func (a *Asset) request(urlPath string) (*http.Request, error) {

	const (
		name  = "name"
		label = "label"
	)

	body, err := a.reader()
	if err != nil {
		return nil, err
	}

	size, err := a.size()
	if err != nil {
		return nil, err
	}

	//body and size
	request, _ := http.NewRequest(http.MethodPost, urlPath, body)
	request.ContentLength = size
	//header
	request.Header.Add(contentType, a.ContentType)

	//query
	query := request.URL.Query()
	query.Add(name, a.Name)
	query.Add(label, a.Label)
	request.URL.RawQuery = query.Encode()
	return request, nil
}
