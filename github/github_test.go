package github

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

const getRealeseByTag = `{
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
		case "/Owner/Repo/releases/tags/6.6.6":
			fmt.Fprintln(w, getRealeseByTag)
		case "/Owner/Repo/releases/tags/StatusNotFound":
			w.WriteHeader(http.StatusNotFound)
		case "/Owner/Repo/releases/tags/BadPayload":
			fmt.Fprintln(w, "no...")
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
				tag: "6.6.6",
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
				tag: "StatusNotFound",
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
				tag: "BadPayload",
			},
			want:    &Release{},
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
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, `{
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
		  }`)
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
					TagName: "v1.0.0",
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
