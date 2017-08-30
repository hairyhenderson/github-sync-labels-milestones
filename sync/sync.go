package sync

import (
	"fmt"

	"github.com/hairyhenderson/github-sync-labels-milestones/config"
)

// Sync - do the thing
func Sync(c *config.Config) {
	g := NewGitHubClient()
	for _, repo := range c.Repositories {
		fmt.Printf("Repo %s\n", repo)
		g.updateMilestones(repo, c.Milestones)
		g.updateLabels(repo, c.Labels)
	}
}
