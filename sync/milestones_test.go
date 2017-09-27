package sync

import (
	"testing"

	"github.com/hairyhenderson/github-sync-labels-milestones/config"
	"github.com/stretchr/testify/assert"
)

func TestSearchMS(t *testing.T) {
	left := &config.Milestone{Title: "foo"}
	right := &config.Milestone{Title: "foo"}
	milestones := []*config.Milestone{right}
	result, _ := searchMS(left, milestones)
	assert.Equal(t, right, result)

	left = &config.Milestone{Title: "foo"}
	right = &config.Milestone{Title: "bar"}
	milestones = []*config.Milestone{right}
	result, _ = searchMS(left, milestones)
	assert.Nil(t, result)

	left = &config.Milestone{Title: "foo", PreviousTitles: []string{"baz"}}
	right = &config.Milestone{Title: "bar"}
	milestones = []*config.Milestone{right}
	result, _ = searchMS(left, milestones)
	assert.Nil(t, result)

	left = &config.Milestone{Title: "foo", PreviousTitles: []string{"baz", "bar"}}
	right = &config.Milestone{Title: "bar"}
	milestones = []*config.Milestone{right}
	result, _ = searchMS(left, milestones)
	assert.Equal(t, right, result)

	left = &config.Milestone{Title: "foo", PreviousTitles: []string{"baz", "bar"}}
	right = &config.Milestone{Title: "bar"}
	right2 := &config.Milestone{Title: "bar"}
	milestones = []*config.Milestone{right, right2}
	_, err := searchMS(left, milestones)
	assert.Error(t, err)
}
