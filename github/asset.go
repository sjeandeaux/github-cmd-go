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

func (a *Asset) request(urlPath string, body io.Reader, size int64) (*http.Request, error) {
	const (
		contentType = "Content-Type"
		name        = "name"
		label       = "label"
	)
	request, _ := http.NewRequest(http.MethodPost, urlPath, body)
	request.ContentLength = size
	query := request.URL.Query()
	query.Add(name, a.Name)
	query.Add(label, a.Label)
	request.URL.RawQuery = query.Encode()
	request.Header.Add(contentType, a.ContentType)
	return request, nil
}
