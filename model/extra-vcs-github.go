package model

/*
// https://github.com/hairyhenderson/github-sync-labels-milestones/blob/master/config/config.go
// https://github.com/hairyhenderson/github-sync-labels-milestones/blob/master/sync/sync.go
// https://github.com/hairyhenderson/github-sync-labels-milestones/blob/master/sync/labels.go
// https://github.com/hairyhenderson/github-sync-labels-milestones/blob/master/sync/milestones.go
// https://github.com/moul/as-a-service/blob/master/github.go
// https://github.com/moul/as-a-service/blob/master/flickr.go
// https://github.com/moul/as-a-service/blob/master/cache.go
// https://github.com/theothertomelliott/github-unwatch/blob/master/app/controllers/app.go
// https://github.com/qiuyesuifeng/community/blob/master/github.go
// https://github.com/crewjam/triggr
// https://github.com/electricbookworks/electric-book-gui

// https://github.com/shurcooL/notifications/blob/master/githubapi/githubapi.go
// https://github.com/m-lab/alertmanager-github-receiver/blob/master/alerts/handler.go

// HunterGate.cmake
// https://github.com/FINTprosjektet/fint-model

// issue-analyzer
// https://github.com/coreos/issue-analyzer/blob/master/repo.go

// pull-request-parser
// https://github.com/guywithnose/pull-request-parser/blob/master/command/github.go

// github-utils search
// https://github.com/parkr/github-utils/tree/master/search
// https://github.com/parkr/github-utils/blob/master/gh/gh.go

// github scores and metrics
// https://github.com/jgautheron/exago/tree/master/score
// https://github.com/jgautheron/exago/blob/master/score/rank.go
// https://github.com/jgautheron/exago/blob/master/gosearch/gosearch.go
// https://github.com/jgautheron/exago/blob/master/showcaser/showcaser.go

// job candidates
// https://github.com/dziemba/seeker


// reports
// https://github.com/shumipro/gh-report-crawler
// https://github.com/bassam/stargazers
// https://github.com/spencerkimball/repo-digest


// travis Repo
// https://github.com/Jimdo/repos/blob/master/github.go
// https://github.com/Jimdo/repos/blob/master/github.go#L57-L59
// https://github.com/Jimdo/repos/blob/master/main.go#L59

*/

/*

import (
	"github.com/SlyMarbo/rss"
	"github.com/google/go-github/github"
	"github.com/patrickmn/go-cache"
)

var githubFeed *rss.Feed

type ExtraGithub_Biblio struct {
	Cache  	map[string]*RepositoryInfo 					`yaml:"-" json:"-"`
	Client 	*github.Client 								`yaml:"-" json:"-"`
}

type ExtraGithub_RepositoryInfo struct {
	LastSyncedIssue struct {
		IssueNumber 			int 					`yaml:"issue_number,omitempty" json:"issue_number,omitempty"`
		Count       			int 					`yaml:"count,omitempty" json:"count,omitempty"`
	} `yaml:"last_synced_issue,omitempty" json:"last_synced_issue,omitempty"`
	Stargazers      			[]string  				`yaml:"stargazers,omitempty" json:"stargazers,omitempty"`
	Subscribers     			[]string 				`yaml:"suscribers,omitempty" json:"suscribers,omitempty"`
	LastPR        				int 					`yaml:"last_pr,omitempty" json:"last_pr,omitempty"`
	ReleasesCount 				int 					`yaml:"releases_count,omitempty" json:"releases_count,omitempty"`
	ForksCount    				int 					`yaml:"forks_count,omitempty" json:"forks_count,omitempty"`
}

// Configs -
type ExtraGithub_Configs []ExtraGithub_Config

// Config -
type ExtraGithub_Config struct {
	ExtraGithub_Repositories []*ExtraGithub_RepositoryInfo `json:"repositories" yaml:"repositories"`
	ExtraGithub_Milestones   []*ExtraGithub_Milestone  	`json:"milestones" yaml:"milestones"`
	ExtraGithub_Labels       []*ExtraGithub_Label      	`json:"labels" yaml:"labels"`
}

// Repository -
type ExtraGithub_Repository struct {
	User 					string  					`json:"user" yaml:"user"`
	Repo 					string  					`json:"repo" yaml:"repo"`
}

// FromGH - convert from github data model
func (r *ExtraGithub_Repository) FromGH(repo *github.Repository) {
	r.UnmarshalText([]byte(*repo.FullName))
}

// MarshalText -
// func (r ExtraGithub_Repository) MarshalText() (text []byte, err error) {
// 	return []byte(r.String()), nil
// }

// UnmarshalText -
func (r *ExtraGithub_Repository) UnmarshalText(text []byte) (err error) {
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

func (r *ExtraGithub_Repository) String() string {
	return r.User + "/" + r.Repo
}

// Milestone -
type ExtraGithub_Milestone struct {
	Title          string    `json:"title" yaml:"title"`
	State          string    `json:"state" yaml:"state"`
	Description    string    `json:"description" yaml:"description"`
	DueOn          time.Time `json:"due_on" yaml:"due_on"`
	PreviousTitles []string  `json:"previous_titles,omitempty" yaml:"previous_titles,omitempty"`
	Number         int       `json:"number,omitempty" yaml:"number,omitempty"`
}


// Equals - determine whether or not two milestones are _mostly_ equal.
// The DueOn property must simply be within the same (UTC) day
func (m *ExtraGithub_Milestone) Equals(o *ExtraGithub_Milestone) bool {
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
	if len(o.PreviousTitles) != len(m.PreviousTitles) {
		return false
	}
	for _, ot := range o.PreviousTitles {
		found := false
		for _, mt := range m.PreviousTitles {
			if ot == mt {
				found = true
			}
		}
		if !found {
			return false
		}
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
func NewMilestonesFromGH(gms []*github.Milestone) []*ExtraGithub_Milestone {
	a := []*ExtraGithub_Milestone{}
	for _, g := range gms {
		a = append(a, NewMilestoneFromGH(g))
	}
	return a
}

// NewMilestoneFromGH - convert from github data model
func NewMilestoneFromGH(g *github.Milestone) *ExtraGithub_Milestone {
	m := &ExtraGithub_Milestone{
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
type ExtraGithub_Label struct {
	Name          string   `json:"name" json:"name"`
	Color         string   `json:"color" json:"color"`
	PreviousNames []string `json:"previous_names,omitempty" json:"previous_names,omitempty"`
	State         string   `json:"state,omitempty" json:"state,omitempty"`
}

// Equals - determine whether or not two labels are equal.
func (l *ExtraGithub_Label) Equals(o *ExtraGithub_Label) bool {
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
func NewLabelsFromGH(gl []*github.Label) []*ExtraGithub_Label {
	a := []*ExtraGithub_Label{}
	for _, g := range gl {
		a = append(a, NewLabelFromGH(g))
	}
	return a
}

// NewLabelFromGH - convert from github data model
func NewLabelFromGH(g *github.Label) *ExtraGithub_Label {
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

*/