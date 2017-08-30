package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMilestoneEquals(t *testing.T) {
	left := &Milestone{}
	right := &Milestone{}
	assert.True(t, left.Equals(right))

	left = &Milestone{Title: "foo"}
	right = &Milestone{Title: "foo"}
	assert.True(t, left.Equals(right))

	left = &Milestone{State: "open"}
	right = &Milestone{State: "open"}
	assert.True(t, left.Equals(right))

	left = &Milestone{Description: "bar"}
	right = &Milestone{Description: "bar"}
	assert.True(t, left.Equals(right))

	left = &Milestone{Title: "foo"}
	right = &Milestone{Title: "bar"}
	assert.False(t, left.Equals(right))

	left = &Milestone{Description: "bar"}
	right = &Milestone{Description: "baz"}
	assert.False(t, left.Equals(right))

	left = &Milestone{State: "open"}
	right = &Milestone{State: "closed"}
	assert.False(t, left.Equals(right))

	left = &Milestone{PreviousTitles: nil}
	right = &Milestone{PreviousTitles: nil}
	assert.True(t, left.Equals(right))

	left = &Milestone{PreviousTitles: []string{}}
	right = &Milestone{PreviousTitles: []string{}}
	assert.True(t, left.Equals(right))

	left = &Milestone{PreviousTitles: []string{"foo"}}
	right = &Milestone{PreviousTitles: []string{}}
	assert.False(t, left.Equals(right))

	left = &Milestone{PreviousTitles: []string{"foo", "bar", "baz"}}
	right = &Milestone{PreviousTitles: []string{"bar", "foo", "baz"}}
	assert.True(t, left.Equals(right))

	left = &Milestone{PreviousTitles: []string{"foo", "bar", "baz"}}
	right = &Milestone{PreviousTitles: []string{"bar", "foo", "qux"}}
	assert.False(t, left.Equals(right))

	left = &Milestone{Number: 0}
	right = &Milestone{Number: 0}
	assert.True(t, left.Equals(right))

	left = &Milestone{Number: 42}
	right = &Milestone{Number: 43}
	assert.False(t, left.Equals(right))

	left = &Milestone{DueOn: time.Time{}}
	right = &Milestone{DueOn: time.Time{}}
	assert.True(t, left.Equals(right))

	today, _ := time.Parse("2006-01-02T15:04:05Z0700", "2017-01-10T23:59:59Z")
	morning, _ := time.Parse("2006-01-02T15:04:05Z0700", "2017-01-10T06:00:00Z")
	lastmonth, _ := time.Parse("2006-01-02T15:04:05Z0700", "2016-12-10T23:59:59Z")
	left = &Milestone{DueOn: today}
	right = &Milestone{DueOn: today}
	assert.True(t, left.Equals(right))

	left = &Milestone{}
	right = &Milestone{DueOn: today}
	assert.False(t, left.Equals(right))

	assert.False(t, left.Equals(nil))

	left = &Milestone{DueOn: today}
	right = &Milestone{DueOn: morning}
	assert.True(t, left.Equals(right))

	left = &Milestone{DueOn: today}
	right = &Milestone{DueOn: lastmonth}
	assert.False(t, left.Equals(right))
}
