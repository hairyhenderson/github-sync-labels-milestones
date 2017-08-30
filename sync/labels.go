package sync

import (
	"fmt"

	"github.com/hairyhenderson/github-sync-labels-milestones/config"

	"github.com/google/go-github/github"
)

func (g *GitHubClient) updateLabels(repo *config.Repository, labels []*config.Label) error {
	existingLabels, err := g.getAllLabels(repo)
	if err != nil {
		return err
	}
	for _, label := range labels {
		existing, err := g.searchLabels(label, existingLabels)
		if err != nil {
			return err
		}

		if label.State == "absent" && existing != nil {
			fmt.Printf("%s found, deleting!\n", label.Name)
			err := g.deleteLabel(repo, label)
			if err != nil {
				return err
			}
			continue
		}
		if label.State == "absent" && existing == nil {
			// fmt.Printf("%s absent, nothing to do!\n", label.Name)
			continue
		}

		if existing.Color == label.Color {
			// fmt.Printf("%s up-to-date, nothing to do!\n", label.Name)
			continue
		}

		if existing == nil {
			fmt.Printf("%s not found, creating!\n", label.Name)
			err := g.createLabel(repo, label)
			if err != nil {
				return err
			}
			continue
		}

		{
			fmt.Printf("%s found, updating!\n\told: %+v\n\tnew: %+v\n", label.Name, existing, label)
			err := g.updateLabel(repo, label)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// getAllLabels - Retrieve a list of all labels from GH
func (g *GitHubClient) getAllLabels(repo *config.Repository) ([]*config.Label, error) {
	opt := &github.ListOptions{PerPage: 100}

	var allLabels []*github.Label
	for {
		labels, resp, err := g.client.Issues.ListLabels(repo.User, repo.Repo, opt)
		if err != nil {
			return nil, err
		}
		allLabels = append(allLabels, labels...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	ms := config.NewLabelsFromGH(allLabels)
	return ms, nil
}

func (g *GitHubClient) searchLabels(label *config.Label, labels []*config.Label) (*config.Label, error) {
	found := []*config.Label{}
	for _, l := range labels {
		if label.Name == l.Name {
			found = append(found, l)
		}
		for _, n := range l.PreviousNames {
			if label.Name == n {
				found = append(found, l)
			}
		}
	}
	// if multiples are found, panic!
	if len(found) > 1 {
		return nil, fmt.Errorf("error: more than 1 matching remote label for %s", label.Name)
	}
	if len(found) == 1 {
		return found[0], nil
	}
	return nil, nil
}

func (g *GitHubClient) createLabel(repo *config.Repository, label *config.Label) error {
	l := &github.Label{Name: &label.Name, Color: &label.Color}
	_, resp, err := g.client.Issues.CreateLabel(repo.User, repo.Repo, l)
	if err != nil {
		return err
	}
	fmt.Printf("create label: %d (%+v)\n", resp.StatusCode, resp)
	return nil
}

func (g *GitHubClient) deleteLabel(repo *config.Repository, label *config.Label) error {
	resp, err := g.client.Issues.DeleteLabel(repo.User, repo.Repo, label.Name)
	if err != nil {
		return err
	}
	fmt.Printf("delete label: %d (%+v)\n", resp.StatusCode, resp)
	return nil
}

func (g *GitHubClient) updateLabel(repo *config.Repository, label *config.Label) error {
	l := &github.Label{Name: &label.Name, Color: &label.Color}
	_, resp, err := g.client.Issues.EditLabel(repo.User, repo.Repo, label.Name, l)
	if err != nil {
		return err
	}
	fmt.Printf("edited label: %d (%+v)\n", resp.StatusCode, resp)
	return nil
}
