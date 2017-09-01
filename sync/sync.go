package sync

import (
	"fmt"

	"github.com/hairyhenderson/github-sync-labels-milestones/config"
	"github.com/pkg/errors"
)

// Sync - do the thing
func Sync(cachePath string, c *config.Config) error {
	g := NewGitHubClient(cachePath)
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
