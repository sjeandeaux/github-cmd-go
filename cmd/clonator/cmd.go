package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"sync"

	internalcmd "github.com/sjeandeaux/toolators/internal/cmd"
	"golang.org/x/oauth2"
)

//Tokens in ~/.clonator
type Tokens struct {
	Bitbucket string
	Github    string
}

//commandLine the arguments command line
type commandLine struct {
	internalcmd.CommandLine

	provider string
	name     string
	filter   string
	file     string
	tokens   *Tokens
}

func (c *commandLine) init() *commandLine {

	//flag
	c.Init("[clonator]")

	user, _ := user.Current()
	flag.StringVar(&c.file, "file", fmt.Sprint(user.HomeDir, "/.clonator"), "")
	flag.StringVar(&c.name, "name", "sjeandeaux", "")
	flag.StringVar(&c.filter, "filter", "", "")
	flag.StringVar(&c.provider, "provider", "github-org", "github-user, github-org, bitbucket")
	flag.Parse()
	c.tokens = readToken(c.file)
	return c

}

func readToken(file string) *Tokens {
	value, _ := ioutil.ReadFile(file)
	var tmpConfig Tokens
	json.Unmarshal(value, &tmpConfig)
	return &tmpConfig
}

func (c *commandLine) main() int {

	provider, err := c.getProvider()
	if err != nil {
		return c.Fatal(err)
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: provider.getToken()},
	)

	//Get the UploadURL
	oauthClient := oauth2.NewClient(oauth2.NoContext, ts)

	for {
		r, err := oauthClient.Get(provider.getURL())
		if err != nil {
			return c.Fatal(err)
		}

		_, nextPage, err := provider.parse(r)
		if err != nil {
			return c.Fatal(err)
		}
		defer r.Body.Close()
		println(r.StatusCode)

		if nextPage == "" {
			return 0
		}
		provider.setURL(nextPage)
	}
}

//Provider github,...
type Provider interface {
	getURL() string
	setURL(string)
	parse(r *http.Response) ([]*Project, string, error)
	getToken() string
}

//Github provider
type Github struct {
	url   string
	token string
}

func (g *Github) parse(r *http.Response) ([]*Project, string, error) {
	if links, ok := r.Header["Link"]; ok && len(links) > 0 {
		for _, link := range strings.Split(links[0], ",") {
			segments := strings.Split(strings.TrimSpace(link), ";")
			// try to pull out page parameter
			url, err := url.Parse(segments[0][1 : len(segments[0])-1])
			if err == nil {
				if strings.TrimSpace(segments[1]) == `rel="next"` {
					return nil, url.String(), nil
				}
			}
		}
		return nil, "", nil
	}
	return nil, "", nil
}

func (g *Github) getToken() string {
	return g.token
}

func (g *Github) getURL() string {
	return g.url
}
func (g *Github) setURL(url string) {
	g.url = url
}

//Bitbucket provider
type Bitbucket struct {
	url   string
	token string
}

func (g *Bitbucket) parse(r *http.Response) ([]*Project, string, error) {
	println(g.getToken())
	type Response struct {
		Next string `json:"next"`
	}
	var tmp Response
	json.NewDecoder(r.Body).Decode(&tmp)
	return nil, tmp.Next, nil
}

func (g *Bitbucket) getToken() string {
	return g.token
}

func (g *Bitbucket) getURL() string {
	return g.url
}
func (g *Bitbucket) setURL(url string) {
	g.url = url
}

func (c *commandLine) getProvider() (Provider, error) {
	const (
		githubUser = "https://api.github.com/users/%s/repos?type=all"
		githubOrg  = "https://api.github.com/orgs/%s/repos?type=all"
		bitbucket  = "https://api.bitbucket.org/2.0/repositories/%s?full_name~%s"
	)
	switch c.provider {
	case "github-user":
		return &Github{url: fmt.Sprintf(githubUser, c.name), token: c.tokens.Github}, nil
	case "github-org":
		return &Github{url: fmt.Sprintf(githubOrg, c.name), token: c.tokens.Github}, nil
	case "bitbucket":
		return &Bitbucket{url: fmt.Sprintf(bitbucket, c.name, c.filter), token: c.tokens.Bitbucket}, nil
	default:
		return nil, fmt.Errorf("i don't know you %q", c.provider)
	}
}

//Project the name and url to clone
type Project struct {
	directory string
	name      string
	cloneURL  string
}

func (p *Project) clone(wg *sync.WaitGroup) {
	d := filepath.Join(p.directory, p.name)
	if _, err := os.Stat(d); os.IsNotExist(err) {
		wg.Add(1)
		cmd := exec.Command("git", "clone", p.cloneURL, d)

		go func(wg *sync.WaitGroup) {
			errCmd := cmd.Start()
			defer wg.Done()
			if errCmd != nil {
				fmt.Printf("error: %v\n\n", errCmd)
			}
			err := cmd.Wait()
			if err != nil {
				fmt.Printf("%s", err)
			}
			fmt.Printf("%q\t\t\t\t\t%q done\n", p.name, p.cloneURL)

		}(wg)
	} else {
		fmt.Printf("%q\t\t\t\t\t%q exists\n", p.name, p.cloneURL)
		//TODO stash and pull
	}
}
