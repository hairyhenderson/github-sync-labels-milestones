package sync

import (
	"fmt"
	"net/http"
	"os"

	"github.com/hairyhenderson/github-sync-labels-milestones/config"

	"github.com/google/go-github/github"
	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
	"golang.org/x/oauth2"
)

// GitHubClient -
type GitHubClient struct {
	client *github.Client
}

func (g *GitHubClient) updateMilestones(repo *config.Repository, milestones []*config.Milestone) error {
	openMSes, err := g.getAllMilestones(repo, "open")
	if err != nil {
		return err
	}
	closedMSes, err := g.getAllMilestones(repo, "closed")
	if err != nil {
		return err
	}
	existingMSes := append(openMSes, closedMSes...)
	if err != nil {
		return err
	}
	// 3. for each _given_ MS:
	for _, ms := range milestones {
		// 3.1. is this MS in (1+2)? (include Title + PreviousTitles in search)
		existingMS, err := g.searchMS(ms, existingMSes)
		if err != nil {
			return err
		}
		if existingMS != nil {
			ms.Number = existingMS.Number
		}

		if existingMS.Equals(ms) {
			// fmt.Printf("%s up-to-date, nothing to do!\n", ms.Title)
		} else if existingMS == nil && ms.State == "absent" {
			// fmt.Printf("%s absent, nothing to do!\n", ms.Title)
		} else if existingMS == nil && ms.State != "absent" {
			// ms doesn't exist yet, create it
			fmt.Printf("%s not found, creating!\n", ms.Title)
			err := g.createMilestone(repo, ms)
			if err != nil {
				return err
			}
		} else if existingMS != nil && ms.State == "absent" {
			fmt.Printf("%s found, deleting!\n", ms.Title)
			err := g.deleteMilestone(repo, ms)
			if err != nil {
				return err
			}
		} else {
			// update milestone in GitHub
			fmt.Printf("%s found, updating!\n\told: %+v\n\tnew: %+v\n", ms.Title, existingMS, ms)
			err := g.updateMilestone(repo, ms)
			if err != nil {
				return err
			}
		}
	}
	// 3.1.
	return nil
}

// getAllMilestones - Retrieve a list of all milestones from GH
func (g *GitHubClient) getAllMilestones(repo *config.Repository, state string) ([]*config.Milestone, error) {
	// get from GitHub
	opt := &github.MilestoneListOptions{
		State:       state,
		ListOptions: github.ListOptions{PerPage: 100},
	}
	var allMilestones []*github.Milestone
	for {
		milestones, resp, err := g.client.Issues.ListMilestones(repo.User, repo.Repo, opt)
		if err != nil {
			return nil, err
		}
		allMilestones = append(allMilestones, milestones...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	ms := config.NewMilestonesFromGH(allMilestones)
	return ms, nil
}

func (g *GitHubClient) searchMS(ms *config.Milestone, milestones []*config.Milestone) (*config.Milestone, error) {
	found := []*config.Milestone{}
	for _, m := range milestones {
		if ms.Title == m.Title {
			found = append(found, m)
		}
		for _, t := range m.PreviousTitles {
			if ms.Title == t {
				found = append(found, m)
			}
		}
	}
	// if multiples are found, panic!
	if len(found) > 1 {
		return nil, fmt.Errorf("error: more than 1 matching remote milestone for %s", ms.Title)
	}
	if len(found) == 1 {
		return found[0], nil
	}
	return nil, nil
}

func (g *GitHubClient) createMilestone(repo *config.Repository, ms *config.Milestone) error {
	m := &github.Milestone{
		Title:       &ms.Title,
		State:       &ms.State,
		Description: &ms.Description,
		DueOn:       &ms.DueOn,
	}
	_, resp, err := g.client.Issues.CreateMilestone(repo.User, repo.Repo, m)
	if err != nil {
		return err
	}
	fmt.Printf("create milestone: %d (%+v)\n", resp.StatusCode, resp)
	return nil
}

func (g *GitHubClient) deleteMilestone(repo *config.Repository, ms *config.Milestone) error {
	resp, err := g.client.Issues.DeleteMilestone(repo.User, repo.Repo, ms.Number)
	if err != nil {
		return err
	}
	fmt.Printf("delete milestone: %d (%+v)\n", resp.StatusCode, resp)
	return nil
}

func (g *GitHubClient) updateMilestone(repo *config.Repository, ms *config.Milestone) error {
	m := &github.Milestone{
		Title:       &ms.Title,
		State:       &ms.State,
		Description: &ms.Description,
		DueOn:       &ms.DueOn,
	}
	_, resp, err := g.client.Issues.EditMilestone(repo.User, repo.Repo, ms.Number, m)
	if err != nil {
		return err
	}
	fmt.Printf("edited milestone: %d (%+v)\n", resp.StatusCode, resp)
	return nil
}

// NewGitHubClient - creates a client, authenticated by OAuth2 via a static token
func NewGitHubClient() *GitHubClient {
	token := os.Getenv("GH_SYNC_TOKEN")
	fmt.Printf("GH Token is %s\n", token)
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	c := diskcache.New(".cache")
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
