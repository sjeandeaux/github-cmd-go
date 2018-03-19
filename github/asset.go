package github

import (
	"io"
	"net/http"
	"os"
)

//Asset information on asset to upload https://developer.github.com/v3/repos/releases/#upload-a-release-asset
type Asset struct {
	File        string
	Name        string
	Label       string
	ContentType string
}

func (a *Asset) reader() (io.ReadCloser, error) {
	return os.Open(a.File)
}

func (a *Asset) size() (int64, error) {
	fileStat, err := os.Stat(a.File)
	if err != nil {
		return -1, err
	}
	return fileStat.Size(), nil
}

func (a *Asset) request(urlPath string) (*http.Request, error) {

	const (
		name  = "name"
		label = "label"
	)

	size, err := a.size()
	if err != nil {
		return nil, err
	}
	body, err := a.reader()
	if err != nil {
		return nil, err
	}

	request, _ := http.NewRequest(http.MethodPost, urlPath, body)
	request.ContentLength = size

	request.Header.Add(contentType, a.ContentType)
	query := request.URL.Query()
	query.Add(name, a.Name)
	query.Add(label, a.Label)
	request.URL.RawQuery = query.Encode()
	return request, nil
}
