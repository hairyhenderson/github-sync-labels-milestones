package sync

import (
	"fmt"
	"net/http"
	"os"

	"github.com/google/go-github/github"
	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
	"golang.org/x/oauth2"
)

// GitHubClient -
type GitHubClient struct {
	client *github.Client
}

// NewGitHubClient - creates a client, authenticated by OAuth2 via a static token
func NewGitHubClient(cachePath string) *GitHubClient {
	token := os.Getenv("GH_SYNC_TOKEN")
	fmt.Printf("GH Token is %s\n", token)
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	c := diskcache.New(cachePath)
	t := httpcache.NewTransport(c)
	hc := &http.Client{
		Transport: &oauth2.Transport{
			Base:   t,
			Source: ts,
		},
	}
	client := github.NewClient(hc)

	return &GitHubClient{client: client}
}
