package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"testing"
)

const realeseByTag = `{
	"url": "https://api.github.com/repos/octocat/Hello-World/releases/1",
	"html_url": "https://github.com/octocat/Hello-World/releases/v1.0.0",
	"assets_url": "https://api.github.com/repos/octocat/Hello-World/releases/1/assets",
	"upload_url": "https://uploads.github.com/repos/octocat/Hello-World/releases/1/assets{?name,label}",
	"tarball_url": "https://api.github.com/repos/octocat/Hello-World/tarball/v1.0.0",
	"zipball_url": "https://api.github.com/repos/octocat/Hello-World/zipball/v1.0.0",
	"id": 1,
	"tag_name": "v1.0.0",
	"target_commitish": "master",
	"name": "v1.0.0",
	"body": "Description of the release",
	"draft": false,
	"prerelease": false,
	"created_at": "2013-02-27T19:35:32Z",
	"published_at": "2013-02-27T19:35:32Z",
	"author": {
	  "login": "octocat",
	  "id": 1,
	  "avatar_url": "https://github.com/images/error/octocat_happy.gif",
	  "gravatar_id": "",
	  "url": "https://api.github.com/users/octocat",
	  "html_url": "https://github.com/octocat",
	  "followers_url": "https://api.github.com/users/octocat/followers",
	  "following_url": "https://api.github.com/users/octocat/following{/other_user}",
	  "gists_url": "https://api.github.com/users/octocat/gists{/gist_id}",
	  "starred_url": "https://api.github.com/users/octocat/starred{/owner}{/repo}",
	  "subscriptions_url": "https://api.github.com/users/octocat/subscriptions",
	  "organizations_url": "https://api.github.com/users/octocat/orgs",
	  "repos_url": "https://api.github.com/users/octocat/repos",
	  "events_url": "https://api.github.com/users/octocat/events{/privacy}",
	  "received_events_url": "https://api.github.com/users/octocat/received_events",
	  "type": "User",
	  "site_admin": false
	},
	"assets": [
	  {
		"url": "https://api.github.com/repos/octocat/Hello-World/releases/assets/1",
		"browser_download_url": "https://github.com/octocat/Hello-World/releases/download/v1.0.0/example.zip",
		"id": 1,
		"name": "example.zip",
		"label": "short description",
		"state": "uploaded",
		"content_type": "application/zip",
		"size": 1024,
		"download_count": 42,
		"created_at": "2013-02-27T19:35:32Z",
		"updated_at": "2013-02-27T19:35:32Z",
		"uploader": {
		  "login": "octocat",
		  "id": 1,
		  "avatar_url": "https://github.com/images/error/octocat_happy.gif",
		  "gravatar_id": "",
		  "url": "https://api.github.com/users/octocat",
		  "html_url": "https://github.com/octocat",
		  "followers_url": "https://api.github.com/users/octocat/followers",
		  "following_url": "https://api.github.com/users/octocat/following{/other_user}",
		  "gists_url": "https://api.github.com/users/octocat/gists{/gist_id}",
		  "starred_url": "https://api.github.com/users/octocat/starred{/owner}{/repo}",
		  "subscriptions_url": "https://api.github.com/users/octocat/subscriptions",
		  "organizations_url": "https://api.github.com/users/octocat/orgs",
		  "repos_url": "https://api.github.com/users/octocat/repos",
		  "events_url": "https://api.github.com/users/octocat/events{/privacy}",
		  "received_events_url": "https://api.github.com/users/octocat/received_events",
		  "type": "User",
		  "site_admin": false
		}
	  }
	]
  }`

func TestClient_GetReleaseByTag(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.RequestURI {
		case "/Owner/Repo/releases/tags/6.6.6.OK":
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, realeseByTag)
		case "/Owner/Repo/releases/tags/6.6.6.NotFound":
			w.WriteHeader(http.StatusNotFound)
		case "/Owner/Repo/releases/tags/6.6.6.BadPayload":
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "no...")
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

	}))
	defer ts.Close()

	type fields struct {
		httpClient *http.Client
		owner      string
		repo       string
		baseURL    string
	}
	type args struct {
		tag string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Release
		wantErr bool
	}{
		{
			name: "OK",
			fields: fields{
				httpClient: ts.Client(),
				owner:      "Owner",
				repo:       "Repo",
				baseURL:    ts.URL,
			},
			args: args{
				tag: "6.6.6.OK",
			},
			want: &Release{
				TagName:           "v1.0.0",
				URL:               "https://api.github.com/repos/octocat/Hello-World/releases/1",
				UploadURLTemplate: "https://uploads.github.com/repos/octocat/Hello-World/releases/1/assets{?name,label}",
			},
			wantErr: false,
		},
		{
			name: "StatusNotFound",
			fields: fields{
				httpClient: ts.Client(),
				owner:      "Owner",
				repo:       "Repo",
				baseURL:    ts.URL,
			},
			args: args{
				tag: "6.6.6.NotFound",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "BadPayload",
			fields: fields{
				httpClient: ts.Client(),
				owner:      "Owner",
				repo:       "Repo",
				baseURL:    ts.URL,
			},
			args: args{
				tag: "6.6.6.BadPayload",
			},
			want:    &Release{},
			wantErr: true,
		},
		{
			name: "no server",
			fields: fields{
				httpClient: ts.Client(),
				owner:      "Owner",
				repo:       "Repo",
				baseURL:    "http://localhost:666",
			},
			args: args{
				tag: "6.6.6.BadPayload",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				httpClient: tt.fields.httpClient,
				owner:      tt.fields.owner,
				repo:       tt.fields.repo,
				baseURL:    tt.fields.baseURL,
			}
			got, err := c.GetReleaseByTag(tt.args.tag)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetReleaseByTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetReleaseByTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_CreateRelease(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		input := &EditRelease{}
		b := r.Body
		if b != nil {
			defer b.Close()
			json.NewDecoder(b).Decode(input)
		}

		switch input.TagName {
		case "6.6.6.OK":
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintln(w, realeseByTag)
		case "6.6.6.NotFound":
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, realeseByTag)
		case "6.6.6.BadPayload":
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintln(w, "no...")
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer ts.Close()
	type fields struct {
		httpClient *http.Client
		owner      string
		repo       string
		baseURL    string
	}
	type args struct {
		edit *EditRelease
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Release
		wantErr bool
	}{
		{
			name: "OK",
			args: args{
				&EditRelease{
					TagName: "6.6.6.OK",
				},
			},
			fields: fields{
				httpClient: ts.Client(),
				owner:      "Owner",
				repo:       "Repo",
				baseURL:    ts.URL,
			},
			want: &Release{
				TagName:           "v1.0.0",
				URL:               "https://api.github.com/repos/octocat/Hello-World/releases/1",
				UploadURLTemplate: "https://uploads.github.com/repos/octocat/Hello-World/releases/1/assets{?name,label}",
			},
			wantErr: false,
		},
		{
			name: "BadPayload",
			args: args{
				&EditRelease{
					TagName: "6.6.6.BadPayload",
				},
			},
			fields: fields{
				httpClient: ts.Client(),
				owner:      "Owner",
				repo:       "Repo",
				baseURL:    ts.URL,
			},
			want:    &Release{},
			wantErr: true,
		},
		{
			name: "NotFound",
			args: args{
				&EditRelease{
					TagName: "6.6.6.NotFound",
				},
			},
			fields: fields{
				httpClient: ts.Client(),
				owner:      "Owner",
				repo:       "Repo",
				baseURL:    ts.URL,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "no server",
			args: args{
				&EditRelease{
					TagName: "6.6.6.NotFound",
				},
			},
			fields: fields{
				httpClient: ts.Client(),
				owner:      "Owner",
				repo:       "Repo",
				baseURL:    "http://localhost:666",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				httpClient: tt.fields.httpClient,
				owner:      tt.fields.owner,
				repo:       tt.fields.repo,
				baseURL:    tt.fields.baseURL,
			}
			got, err := c.CreateRelease(tt.args.edit)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.CreateRelease() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.CreateRelease() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRelease_UploadURL(t *testing.T) {
	type fields struct {
		UploadURLTemplate string
		TagName           string
		URL               string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			fields: fields{
				UploadURLTemplate: "http://bob-the-sponge.com{?label, name}",
			},
			want: "http://bob-the-sponge.com",
		},
		{
			fields: fields{
				UploadURLTemplate: "http://bob-the-sponge.com",
			},
			want: "http://bob-the-sponge.com",
		},
		{
			fields: fields{
				UploadURLTemplate: "",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Release{
				UploadURLTemplate: tt.fields.UploadURLTemplate,
				TagName:           tt.fields.TagName,
				URL:               tt.fields.URL,
			}
			if got := o.UploadURL(); got != tt.want {
				t.Errorf("Release.UploadURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_Upload(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	type fields struct {
		httpClient *http.Client
		owner      string
		repo       string
		baseURL    string
	}
	type args struct {
		urlPath string
		a       *Asset
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				httpClient: ts.Client(),
				owner:      "Owner",
				repo:       "Repo",
				baseURL:    ts.URL,
			},
			args: args{
				urlPath: "/path",
				a: &Asset{
					File:        filepath.Join("testdata", "data"),
					ContentType: "application/binary",
					Name:        "fileName",
					Label:       "Label",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				httpClient: tt.fields.httpClient,
				owner:      tt.fields.owner,
				repo:       tt.fields.repo,
				baseURL:    tt.fields.baseURL,
			}
			if err := c.Upload(tt.args.urlPath, tt.args.a); (err != nil) != tt.wantErr {
				t.Errorf("Client.Upload() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
