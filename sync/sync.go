package sync

import (
	"fmt"

	"github.com/hairyhenderson/github-sync-labels-milestones/config"
	"github.com/pkg/errors"
)

// Options -
type Options struct {
	CachePath string
	NoCache   bool
	DryRun    bool
}

// Sync - do the thing
func Sync(c *config.Config, opts Options) error {
	g := NewGitHubClient(opts)
	for _, repo := range c.Repositories {
		fmt.Printf("Repo %s\n", repo)
		err := g.updateMilestones(repo, c.Milestones)
		if err != nil {
			return errors.Wrapf(err, "failed to update milestones on %s", repo)
		}
		err = g.updateLabels(repo, c.Labels)
		if err != nil {
			return errors.Wrapf(err, "failed to update labels on %s", repo)
		}
	}
	return nil
}
