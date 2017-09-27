package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/github"
)

// Configs -
type Configs []Config

// Config -
type Config struct {
	Repositories []*Repository `json:"repositories"`
	Milestones   []*Milestone  `json:"milestones"`
	Labels       []*Label      `json:"labels"`
}

// Repository -
type Repository struct {
	User string
	Repo string
}

// FromGH - convert from github data model
func (r *Repository) FromGH(repo *github.Repository) {
	r.UnmarshalText([]byte(*repo.FullName))
}

// MarshalText -
// func (r Repository) MarshalText() (text []byte, err error) {
// 	return []byte(r.String()), nil
// }

// UnmarshalText -
func (r *Repository) UnmarshalText(text []byte) (err error) {
	s := strings.SplitN(string(text), "/", 2)
	if len(s) != 2 {
		return fmt.Errorf("error: wrong format for repo '%s' (%#v)", text, s)
	}
	*r = Repository{
		User: s[0],
		Repo: s[1],
	}
	return nil
}

func (r *Repository) String() string {
	return r.User + "/" + r.Repo
}

// Milestone -
type Milestone struct {
	Title          string    `json:"title"`
	State          string    `json:"state"`
	Description    string    `json:"description"`
	DueOn          time.Time `json:"due_on"`
	PreviousTitles []string  `json:"previousTitles,omitempty"`
	Number         int       `json:"number,omitempty"`
}

// Equals - determine whether or not two milestones are _mostly_ equal.
// The DueOn property must simply be within the same (UTC) day
func (m *Milestone) Equals(o *Milestone) bool {
	if o == nil || m == nil {
		return false
	}
	if o.Title != m.Title {
		return false
	}
	if o.State != m.State {
		return false
	}
	if o.Description != m.Description {
		return false
	}
	if o.Number != m.Number {
		return false
	}

	mDay := m.DueOn.Format("2006-01-02")
	oDay := o.DueOn.Format("2006-01-02")
	if mDay != oDay {
		return false
	}

	return true
}

// NewMilestonesFromGH - convert from github data model
func NewMilestonesFromGH(gms []*github.Milestone) []*Milestone {
	a := []*Milestone{}
	for _, g := range gms {
		a = append(a, NewMilestoneFromGH(g))
	}
	return a
}

// NewMilestoneFromGH - convert from github data model
func NewMilestoneFromGH(g *github.Milestone) *Milestone {
	m := &Milestone{
		Title:  *(g.Title),
		State:  *(g.State),
		Number: *(g.Number),
	}
	if g.Description != nil {
		m.Description = *(g.Description)
	}
	if g.DueOn != nil {
		m.DueOn = *(g.DueOn)
	}
	return m
}

// Label -
type Label struct {
	Name          string   `json:"name"`
	Color         string   `json:"color"`
	PreviousNames []string `json:"previousNames,omitempty"`
	State         string   `json:"state,omitempty"`
}

// Equals - determine whether or not two labels are equal.
func (l *Label) Equals(o *Label) bool {
	if o == nil || l == nil {
		return false
	}
	if l.Name != o.Name {
		return false
	}
	if l.Color != o.Color {
		return false
	}

	return true
}

// NewLabelsFromGH - convert from github data model
func NewLabelsFromGH(gl []*github.Label) []*Label {
	a := []*Label{}
	for _, g := range gl {
		a = append(a, NewLabelFromGH(g))
	}
	return a
}

// NewLabelFromGH - convert from github data model
func NewLabelFromGH(g *github.Label) *Label {
	return &Label{
		Name:  *(g.Name),
		Color: *(g.Color),
	}
}

// ParseFile -
func ParseFile(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	d := json.NewDecoder(f)
	c := &Config{}
	err = d.Decode(c)
	return c, err
}
