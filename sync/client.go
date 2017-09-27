package sync

import (
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
	dryRun bool
}

// NewGitHubClient - creates a client, authenticated by OAuth2 via a static token
func NewGitHubClient(opts Options) *GitHubClient {
	token := os.Getenv("GH_SYNC_TOKEN")
	// fmt.Printf("GH Token is %s\n", token)
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	var hc *http.Client
	if !opts.NoCache {
		c := diskcache.New(opts.CachePath)
		t := httpcache.NewTransport(c)
		hc = &http.Client{Transport: &oauth2.Transport{Base: t, Source: ts}}
	} else {
		hc = &http.Client{Transport: &oauth2.Transport{Source: ts}}
	}
	client := github.NewClient(hc)

	return &GitHubClient{client: client, dryRun: opts.DryRun}
}
