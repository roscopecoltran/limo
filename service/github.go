package service

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/oauth2"
	log "github.com/sirupsen/logrus"
	"github.com/google/go-github/github"
	"github.com/hoop33/entrevista"
	"github.com/roscopecoltran/sniperkit-limo/model"
	// tablib "github.com/agrison/go-tablib"
	// "github.com/davecgh/go-spew/spew"
	// jsoniter "github.com/json-iterator/go"
	// fuzz "github.com/google/gofuzz"
)

// Github represents the Github service
// GitHub holds specific information that is used for GitHub integration.
type Github struct {	
	/*
	Client 			*github.Client 						`json:"-"` 						// Holds the client instance details. Internal only.
	SearchRepo 		*github.RepositoriesSearchResult	`json:"-"`						//
	SearchCode 		*github.CodeSearchResult 			`json:"-"`						//
	SearchCommits 	*github.CommitsSearchResult 		`json:"-"`						//
	SearchIssues 	*github.IssuesSearchResult 			`json:"-"`						//
	SearchUsers 	*github.UsersSearchResult 			`json:"-"`						//
	User 			*github.User 						`json:"-"`						//
	*/
	IgnoreRepos 	[]string 							`json:"ignoreRepos,omitempty"`	// A list of URLs that the bot can ignore.
	Repo 			*github.Repository 		  			`json:"-"`						//
	Starred 		*github.StarredRepository 			`json:"-"`						//
	Catalog  		map[string][]github.Repository 		`json:"-"`
	Username 		string 								`json:"-"`
	Type     		string 								`json:"-"`
	PerPage  		int 								`json:"-"`
}

// Login logs in to Github
func (g *Github) Login(ctx context.Context) (string, error) {
	interview := createInterview()
	interview.Questions = []entrevista.Question{
		{
			Key:      "token",
			Text:     "Enter your GitHub API token",
			Required: true,
			Hidden:   true,
		},
	}
	answers, err := interview.Run()
	if err != nil {
		return "", err
	}
	return answers["token"].(string), nil
}

// about linting code: https://github.com/seiffert/ghrepos/blob/master/scripts/lint

// GetStars returns the stars for the specified user (empty string for authenticated user)
func (g *Github) GetStars(ctx context.Context, starChan chan<- *model.StarResult, token string, user string) {

	log.WithFields(log.Fields{"service": "GetStars"}).Infof("token: %#v", token)
	log.WithFields(log.Fields{"service": "GetStars"}).Infof("user: %#v", user)
	//"application/vnd.github.mercy-preview+json"
	client := g.getClient(token)

	// spew.Dump(client)

	// The first response will give us the correct value for the last page
	currentPage := 1
	lastPage := 1

	// , topics []string
	// https://github.com/seiffert/ghrepos/blob/master/ghrepos.go

	for currentPage <= lastPage {
		// https://github.com/dougt/githubwebpush/blob/master/src/githubpusher/frontend/main.go

		repos, response, err := client.Activity.ListStarred(ctx, user, &github.ActivityListStarredOptions{
			ListOptions: github.ListOptions{
				Page: currentPage,
			},
		})

		// repos, _, _ := client.Repositories.List(input, nil)
		// b, _ := jsoniter.Marshal(repos)
		// fmt.Print(string(b))

		/*
		log.WithFields(log.Fields{"service": "GetStars"}).Info("repos")
		// spew.Dump(repos)

		log.WithFields(log.Fields{"service": "GetStars"}).Info("response")
		// spew.Dump(response)

		log.WithFields(log.Fields{"service": "GetStars"}).Infof("response")
		// spew.Dump(ctx)


		dumpResponse, err := jsoniter.Marshal(&response)
		if err != nil {
			log.WithError(err).WithFields(log.Fields{"service": "GetStars"}).Warn("dump error on response with jsoniter")
		} else {
			fmt.Println(dumpResponse)
		}

		dumpPageResults, err := jsoniter.Marshal(&repos)
		if err != nil {
			log.WithError(err).WithFields(log.Fields{"service": "GetStars"}).Warn("dump error on repos with jsoniter")
		} else {
			fmt.Println(dumpPageResults)
		}
		*/

		// If we got an error, put it on the channel
		if err != nil {
			starChan <- &model.StarResult{
				Error: err,
				Star:  nil,
			}
		} else {
			// ds, _ := tablib.LoadJSON(response)
			// yaml, _ := ds.YAML()
			// fmt.Println(repos)
			// Set last page only if we didn't get an error
			lastPage = response.LastPage
			// Create a Star for each repository and put it on the channel
			for _, repo := range repos {
				/*
				dumpRepoDetails, err := jsoniter.Marshal(&repo)
				if err != nil {
					log.WithError(err).WithFields(log.Fields{"service": "GetStars"}).Warn("dump error on repo details with jsoniter")
				} else {
					fmt.Println(dumpRepoDetails)
				}
				*/
				// spew.Dump(repo)
				// *github.StarredRepository
				star, err := model.NewStarFromGithub(repo.StarredAt, *repo.Repository)
				starChan <- &model.StarResult{
					Error: err,
					Star:  star,
				}
			}
		}
		log.WithFields(log.Fields{"service": "GetStars"}).Infof("currentPage: %#v", currentPage)
		// Go to the next page
		currentPage++
	}
	close(starChan)
}

// GetEvents returns the events for the authenticated user
func (g *Github) GetEvents(ctx context.Context, eventChan chan<- *model.EventResult, token, user string, page, count int) {
	client := g.getClient(token)

	currentPage := page
	lastPage := page + count - 1

	for currentPage <= lastPage {
		events, _, err := client.Activity.ListEventsReceivedByUser(ctx, user, false, &github.ListOptions{
			Page: currentPage,
		})

		if err != nil {
			eventChan <- &model.EventResult{
				Error: err,
				Event: nil,
			}
		} else {
			for _, event := range events {
				eventChan <- &model.EventResult{
					Error: nil,
					Event: model.NewEventFromGithub(event),
				}
			}
		}
		currentPage++
	}
	close(eventChan)
}

// GetTrending returns the trending repositories
func (g *Github) GetTrending(ctx context.Context, trendingChan chan<- *model.StarResult, token string, language string, verbose bool) {
	client := g.getClient(token)
	log.WithFields(log.Fields{"service": "GetTrending"}).Infof("token: %#v", token)
	// TODO perhaps allow them to specify multiple pages?
	// Might be overkill -- first page probably plenty

	// TODO Make this more configurable. Sort by stars, forks, default.
	// Search by number of stars, pushed, created, or whatever.
	// Lots of possibilities.

	q := g.getDateSearchString()

	if language != "" {
		q = fmt.Sprintf("language:%s %s", language, q)
		log.WithFields(log.Fields{"service": "GetTrending"}).Infof("language: %#v", language)
	}

	if verbose {
		fmt.Println("q =", q)
		log.WithFields(log.Fields{"service": "GetTrending"}).Infof("q: %#v", q)
	}

	result, _, err := client.Search.Repositories(ctx, q, &github.SearchOptions{
		Sort:  "stars",
		Order: "desc",
	})

	// If we got an error, put it on the channel
	if err != nil {
		trendingChan <- &model.StarResult{
			Error: err,
			Star:  nil,
		}
	} else {
		// Create a Star for each repository and put it on the channel
		for _, repo := range result.Repositories {
			star, err := model.NewStarFromGithub(nil, repo)
			trendingChan <- &model.StarResult{
				Error: err,
				Star:  star,
			}
		}
	}

	close(trendingChan)
}


// get CreatedAt from repo
func getCreatedAtFromRepo(ctx context.Context, client *github.Client, owner string, repo string) (createdAt time.Time, err error) {
	repoinfo, _, err := client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		fmt.Println(err)
		return
	}
	var shortForm = "2006-01-02 15:04:05 -0700 UTC"
	ctime, _ := time.Parse(shortForm, fmt.Sprintf("%s", repoinfo.CreatedAt))

	return ctime, nil
}

func (g *Github) getDateSearchString() string {
	// TODO make this configurable
	// Default should be in configuration file
	// and should be able to override from command line
	// TODO should be able to specify whether "created" or "pushed"
	date := time.Now().Add(-7 * (24 * time.Hour))
	log.WithFields(log.Fields{"service": "getDateSearchString"}).Infof("date > %#v", date)
	return fmt.Sprintf("created:>%s", date.Format("2006-01-02"))
}

func (g *Github) getClient(token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	return github.NewClient(tc)
}

/*
func NewSyncer(conf Config, data DataStore) *Syncer {
	ts := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: conf.Token,
	})

	tc := oauth2.NewClient(oauth2.NoContext, ts)

	return &Syncer{
		conf:   conf,
		client: github.NewClient(tc),
		data:   data,
		//logger: level.New(log.NewContext(log.NewLogfmtLogger(os.Stderr)).With("ts", log.DefaultTimestamp)),
	}
}

// Syncer dog
type Syncer struct {
	conf   Config
	client *github.Client
	http   httpdown.Server
	data   DataStore
	//logger level.Option
}

func (s *Syncer) StartSyncer() error {
	if err := s.syncRepos(); err != nil {
		return err
	}

	if err := s.syncMembers(); err != nil {
		return err
	}

	if err := s.syncTeams(); err != nil {
		return err
	}

	if err := s.syncIssues(); err != nil {
		return err
	}

	if err := s.syncIssuesComments(); err != nil {
		return err
	}

	if err := s.syncReviewComments(); err != nil {
		return err
	}

	if err := s.syncCommitComments(); err != nil {
		return err
	}

	srv := http.Server{
		Addr:    s.conf.ListenAddr,
		Handler: http.HandlerFunc(s.handleHook),
	}

	hsrv, err := httpConfig.ListenAndServe(&srv)

	if err != nil {
		return err
	}

	////s.logger.Info().Log("msg", "listening", "addr", s.conf.ListenAddr)

	s.http = hsrv

	return nil
}

func (s *Syncer) Wait() error {
	return s.http.Wait()
}

func (s *Syncer) Stop() error {
	return s.http.Stop()
}

func (s *Syncer) handleHook(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != hookPath {
		http.NotFound(w, r)
		return
	}

	var (
		ev = r.Header.Get("X-Github-Event")
		// dest interface{}
	)

	switch ev {
	default:
		b, _ := httputil.DumpRequest(r, true)
		////s.logger.Warn().Log("msg", "Unhandle event", "event", ev, "request", string(b))
	}
}
*/

/*
func (s *Syncer) syncCommitComments() error {
	return s.data.ForEachRepo(s.syncCommitCommentsByRepo)
}

func (s *Syncer) syncCommitCommentsByRepo(ctx context.Context, r *github.Repository) error {
	var (
		page = 1
		size = 100
	)

	for {
		comments, resp, err := s.client.Repositories.ListComments(ctx, s.conf.Organization, *r.Name, &github.ListOptions{
			Page:    page,
			PerPage: size,
		})

		if err != nil {
			return err
		}

		for i := range comments {
			//s.logger.Debug().Log("msg", "updating commit comment", "repo", *r.Name, "id", comments[i].ID)

			if err := s.data.UpdateCommitComment(&comments[i]); err != nil {
				return err
			}
		}

		if resp.NextPage == 0 {
			break
		}

		page++
	}

	return nil
}
*/

/*
func (s *Syncer) syncIssuesComments() error {
	return s.data.ForEachRepo(s.syncIssuesCommentsByRepo)
}

func (s *Syncer) syncIssuesCommentsByRepo(ctx context.Context, r *github.Repository) error {
	last, err := s.data.LastUpdatedIssueComment(*r.Name)

	if err != nil {
		return err
	}

	var (
		page  = 1
		size  = 100
		since time.Time
	)

	if last != nil {
		since = *last.UpdatedAt
	}

	for {
		comments, resp, err := s.client.Issues.ListComments(ctx, s.conf.Organization, *r.Name, 0, &github.IssueListCommentsOptions{
			Sort:      "updated",
			Direction: "asc",
			Since:     since,
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: size,
			},
		})

		if err != nil {
			return err
		}

		for i := range comments {
			//s.logger.Debug().Log("msg", "updating issue comment", "repo", *r.Name, "id", comments[i].ID)

			if err := s.data.UpdateIssueComment(&comments[i]); err != nil {
				return err
			}
		}

		if resp.NextPage == 0 {
			break
		}

		page++
	}

	return nil
}
*/

/*

func (s *Syncer) syncIssues(ctx context.Context) error {
	last, err := s.data.LastUpdatedIssue()

	if err != nil {
		return err
	}

	var (
		page  = 1
		size  = 100
		since time.Time
	)

	if last != nil {
		since = *last.UpdatedAt
	}

	for {
		issues, resp, err := s.client.Issues.ListByOrg(ctx, s.conf.Organization, &github.IssueListOptions{
			Filter:    "all",
			State:     "all",
			Sort:      "updated",
			Direction: "asc",
			Since:     since,
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: size,
			},
		})

		if err != nil {
			return err
		}

		for i := range issues {
			//s.logger.Debug().Log("msg", "updating issue", "resp", *issues[i].Repository.Name, "issue", *issues[i].Number)

			if err := s.data.UpdateIssue(&issues[i]); err != nil {
				return err
			}
		}

		if resp.NextPage == 0 {
			break
		}

		page++
	}

	return nil
}

*/

/*

func (s *Syncer) syncReviewComments() error {
	return s.data.ForEachRepo(s.syncReviewCommentsByRepo)
}

func (s *Syncer) syncReviewCommentsByRepo(ctx context.Context, r *github.Repository) error {
	last, err := s.data.LastUpdatedReviewComment(*r.Name)

	if err != nil {
		return err
	}

	var (
		page  = 1
		size  = 100
		since time.Time
	)

	if last != nil {
		since = *last.UpdatedAt
	}

	for {
		comments, resp, err := s.client.PullRequests.ListComments(ctx, s.conf.Organization, *r.Name, 0, &github.PullRequestListCommentsOptions{
			Sort:      "updated",
			Direction: "asc",
			Since:     since,
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: size,
			},
		})

		if err != nil {
			return err
		}

		for i := range comments {
			//s.logger.Debug().Log("msg", "updating review comment", "repo", *r.Name, "id", comments[i].ID)

			if err := s.data.UpdateReviewComment(&comments[i]); err != nil {
				return err
			}
		}

		if resp.NextPage == 0 {
			break
		}

		page++
	}

	return nil
}
*/

/*

func (s *Syncer) syncStarred(ctx context.Context) error {
	var (
		page = 1
		size = 100
	)

	for {
		repos, resp, err := s.client.Activity.ListStarred(ctx, s.conf.Organization, &github.RepositoryListByOrgOptions{
			Type: "all",
			ListOptions: github.ListOptions{
				PerPage: size,
				Page:    page,
			},
		})

		if err != nil {
			return err
		}

		for i := range repos {
			//s.logger.Debug().Log("msg", "updating repo", "repo", *repos[i].Name)

			if err := s.data.UpdateRepo(&repos[i]); err != nil {
				return err
			}
		}

		if resp.NextPage == 0 {
			break
		}

		page++
	}

	return nil
}

*/

/*

func (s *Syncer) syncRepos(ctx context.Context) error {
	var (
		page = 1
		size = 100
	)

	for {
		repos, resp, err := s.client.Repositories.ListByOrg(ctx, s.conf.Organization, &github.RepositoryListByOrgOptions{
			Type: "all",
			ListOptions: github.ListOptions{
				PerPage: size,
				Page:    page,
			},
		})

		if err != nil {
			return err
		}

		for i := range repos {
			//s.logger.Debug().Log("msg", "updating repo", "repo", *repos[i].Name)

			if err := s.data.UpdateRepo(&repos[i]); err != nil {
				return err
			}
		}

		if resp.NextPage == 0 {
			break
		}

		page++
	}

	return nil
}

*/

/*
func (s *Syncer) syncMembers() error {
	var (
		page = 1
		size = 100
	)

	for {
		members, resp, err := s.client.Organizations.ListMembers(s.conf.Organization, &github.ListMembersOptions{
			PublicOnly: false,
			Filter:     "all",
			Role:       "all",
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: size,
			},
		})

		if err != nil {
			return err
		}

		for i := range members {
			//s.logger.Debug().Log("msg", "updating member", "member", *members[i].Login)

			if err := s.data.UpdateUser(&members[i]); err != nil {
				return err
			}
		}

		if resp.NextPage == 0 {
			break
		}

		page++
	}

	return nil
}
*/

/*
func (s *Syncer) syncTeams() error {
	var (
		page = 1
		size = 100
	)

	for {
		teams, resp, err := s.client.Organizations.ListTeams(s.conf.Organization, &github.ListOptions{
			Page:    page,
			PerPage: size,
		})

		if err != nil {
			return err
		}

		for i := range teams {
			//s.logger.Debug().Log("msg", "updating team", "team", *teams[i].Name)

			if err := s.data.UpdateTeam(&teams[i]); err != nil {
				return err
			}
		}

		if resp.NextPage == 0 {
			break
		}

		page++
	}

	return nil
}
*/

func init() {
	registerService(&Github{})
}
