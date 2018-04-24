package github

import (
	"io"
	"os"

	internalhttp "github.com/sjeandeaux/toolators/internal/http"
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

var _ UploadInformation = &Asset{}

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

//headers headers http
func (a *Asset) headers() map[string]string {
	return map[string]string{
		internalhttp.ContentType: a.ContentType,
	}
}

//parameters for upload
func (a *Asset) parameters() map[string]string {
	const (
		name  = "name"
		label = "label"
	)
	return map[string]string{
		name:  a.Name,
		label: a.Label,
	}
}
