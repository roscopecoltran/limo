package service

import (
	"context"
	"fmt"
	"time"
	"crypto/md5"
	"strings"
	//"regexp"
	"path"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	//log "github.com/sirupsen/logrus"
	"github.com/google/go-github/github"
	"github.com/hoop33/entrevista"
	"github.com/roscopecoltran/sniperkit-limo/model"
	// tablib "github.com/agrison/go-tablib"
	// "github.com/davecgh/go-spew/spew"
	// jsoniter "github.com/json-iterator/go"
	// fuzz "github.com/google/gofuzz"
)

const GITHUB__RAW_URL = "https://raw.githubusercontent.com/"

/*
const (
	NOTSTART = "NotStart"
	FETCHING = "Fetching"
	INDEXING = "Indexing"
	INDEXED  = "Indexed"
	ERROR    = "Error"
)

type GitHubUser struct {
	Account string `json:"account"`
	// AccessToken       string `json:"accessToken"`
	Tokens            []string `json:"tokens"`
	Status            string   `json:"status"`
	NumOfStarred      int
	IndicesOfStarrerd int
}
*/

type pjson struct {
  Name string
  Description string
  Keywords []string
}

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
	//SubChannels 	bool 								`json:"-"`
	//SubChannelsJobs uint 								`json:"-"`
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

/*
// https://raw.githubusercontent.com/blevesearch/bleve-wiki-indexer/master/git.go
func (g *Github) OpenGitRepo(path string) *github.Repository {
	repo, err := github.OpenRepository(path)
	if err != nil {
		log.Fatal(err)
	}

	return repo
}
*/

// about linting code: https://github.com/seiffert/ghrepos/blob/master/scripts/lint

// ctype, _, err := mime.ParseMediaType(res.Header.Get("Content-Type"))
// if err != nil {
// 	return nil, err
// }

// switch ctype {
// case "application/json", "text/javascript":
// 	var data map[string]interface{}
// 	json.Unmarshal(b, &data)
// 	return data, nil
// }

func gravatarHashFromEmail(email string) string {
	input := strings.ToLower(strings.TrimSpace(email))
	return fmt.Sprintf("%x", md5.Sum([]byte(input)))
}

func searchForFile(files []string, file string) bool {
  for _, b := range files {
    if b == file {
      return true
    }
  }
  return false
}

/*
func (g *Github) GetCommits(ctx context.Context, token string, user string) (*repo.RepositoryContent, string, error) {
	commits, _, err := client.Repositories.ListCommits(ctx, "google", "go-github", nil)
	if err != nil {
		log.Errorf(ctx, "ListCommits: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	for _, commit := range commits {
		fmt.Fprintln(w, commit.GetHTMLURL())
	}
}
*/

/*
func (g *Github) GetUserInfo(ctx context.Context, token string, user string) (*repo.RepositoryContent, string, error) {
	log.WithFields(logrus.Fields{"service": "GetReadme", "token": token, "owner": user, "repo": repo}).Infoln("token: ", token, ", owner: ", user, ", repo: ", repo)
	client := g.getClient(token)
	readme, response, err := client.Users.GetReadme(ctx, user, nil)
	content, err := readme.GetContent()
	if err != nil {
		log.WithFields(logrus.Fields{"service": "GetReadme", "token": token, "owner": owner, "repo": repo}).Warn("error while getting the content of the readme.")
		return nil, "", err
	}
	return readme, content, nil
}
*/

func (g *Github) GetReadme(ctx context.Context, token string, user string, name string) (*github.RepositoryContent, error) {
	log.WithFields(logrus.Fields{"service": "GetReadme", "token": token, "user": user, "repo": name}).Info("")
	client := g.getClient(token)
	readme, _, err := client.Repositories.GetReadme(ctx, user, name, nil)
	// lastPage = response.LastPage
	// content, err := readme.GetContent()
	if err != nil {
		log.WithFields(logrus.Fields{"service": "GetReadme", "token": token, "owner": user, "repo": name}).Warn("error while getting the content of the readme.")
		return nil, err
	}
	return readme, nil
}

// GetStars returns the stars for the specified user (empty string for authenticated user)
func (g *Github) GetStars(ctx context.Context, starChan chan<- *model.StarResult, token string, user string, subChannels bool, subChannelsJobs uint) {
	log.WithFields(logrus.Fields{"service": "GetStars", "token": token, "subChannels": subChannels, "subChannelsJobs": subChannelsJobs}).Info("")
	//log.WithFields(logrus.Fields{"service": "GetStars", "user": user}).Infof("user: %#v", user)
	//"application/vnd.github.mercy-preview+json"
	client := g.getClient(token)
	// The first response will give us the correct value for the last page
	currentPage := 1
	lastPage := 1
	currentStar := 1

	// , topics []string
	// https://github.com/seiffert/ghrepos/blob/master/ghrepos.go
	for currentPage <= lastPage {
		// https://github.com/dougt/githubwebpush/blob/master/src/githubpusher/frontend/main.go
		repos, response, err := client.Activity.ListStarred(ctx, user, &github.ActivityListStarredOptions{
			ListOptions: github.ListOptions{
				Page: currentPage,
			},
		})
		// If we got an error, put it on the channel
		if err != nil {
			starChan <- &model.StarResult{
				Error: err,
				Star:  nil,
			}
		} else {
			// Set last page only if we didn't get an error
			lastPage = response.LastPage
			// Create a Star for each repository and put it on the channel
			for _, repo := range repos {

				// readme
				ownerName 		:= string(*repo.Repository.Owner.Login)
				repoLanguage 	:= string(*repo.Repository.Language)
				repoName 		:= string(*repo.Repository.Name)
				repoUri 		:= path.Join("github.com", *repo.Repository.Owner.Login, *repo.Repository.Name)
				readmeInfo, err := g.GetReadme(ctx, token, ownerName, repoName)
				if err != nil {
					log.WithError(err).WithFields(logrus.Fields{"service": "GetReadme", "star_owner": ownerName, "star_owner_id": *repo.Repository.Owner.ID, "star_name": repoName, "star_id": *repo.Repository.ID}).Warn("error while getting the readme content.")
					readmeInfo = &github.RepositoryContent{}
				} else {
					if _, err := model.NewReadmeFromGithub(*readmeInfo, *repo.Repository.Owner.ID, *repo.Repository.ID, repoUri); err != nil {
						log.WithFields(logrus.Fields{"service": "GetReadme", "step": "NewReadmeFromGithub", "repoUri": repoUri}).Warn("could not fetched new readme.")
					}
				}

				// language
				log.WithFields(logrus.Fields{"service": "GetReadme", "step": "NewLanguageFromGithub", "repoUri": repoUri, "repoLanguage": repoLanguage}).Warn("")
				if *repo.Repository.Language != "" {
					if _, err := model.NewLanguageFromGithub(repoLanguage); err != nil {
						log.WithError(err).WithFields(logrus.Fields{"service": "GetReadme", "step": "NewStarFromGithub", "repoUri": repoUri, "Language": *repo.Repository.Language}).Warn("")
					}
				}

				// userinfo
				/*
				user, err := model.NewStarFromGithub(repo.Owner.ID)
				userChan <- &model.UserResult{
					Error: err,
					User:  user,
				}				
				*/

				// star
				star, err := model.NewStarFromGithub(repo.StarredAt, *repo.Repository, *readmeInfo)
				starChan <- &model.StarResult{
					Error: err,
					Star:  star,
				}

				currentStar++
				log.WithFields(logrus.Fields{"service": "GetStars"}).Infoln("star=",currentStar,",page=",currentPage)
			}
		}
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
	log.WithFields(logrus.Fields{"service": "GetTrending"}).Infof("token: %#v", token)
	// TODO perhaps allow them to specify multiple pages?
	// Might be overkill -- first page probably plenty
	// TODO Make this more configurable. Sort by stars, forks, default.
	// Search by number of stars, pushed, created, or whatever.
	// Lots of possibilities.
	q := g.getDateSearchString()
	if language != "" {
		q = fmt.Sprintf("language:%s %s", language, q)
		log.WithFields(logrus.Fields{"service": "GetTrending"}).Infof("language: %#v", language)
	}
	if verbose {
		fmt.Println("q =", q)
		log.WithFields(logrus.Fields{"service": "GetTrending"}).Infof("q: %#v", q)
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
			star, err := model.NewStarFromGithub(nil, repo, github.RepositoryContent{})
			trendingChan <- &model.StarResult{
				Error: err,
				Star:  star,
			}
		}
	}
	close(trendingChan)
}

// get CreatedAt from repo
// func (g *Github) getCreatedAtFromRepo(owner string, repo string) (createdAt time.Time, err error) {
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
	log.WithFields(logrus.Fields{"service": "getDateSearchString"}).Infof("date > %#v", date)
	return fmt.Sprintf("created:>%s", date.Format("2006-01-02"))
}

func (g *Github) getClient(token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	return github.NewClient(tc)
}

func init() {
	registerService(&Github{})
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

/*
// https://github.com/glena/github-starred-catalog/blob/master/lib/ghclient.go
func (g *Github) GetUsersRepositories(ctx context.Context, starChan chan<- *model.StarResult, token string, user string) {

	log.WithFields(logrus.Fields{"service": "GetStars", "token": token}).Infof("token: %#v", token)
	log.WithFields(logrus.Fields{"service": "GetStars", "user": user}).Infof("user: %#v", user)
	//"application/vnd.github.mercy-preview+json"
	client := g.getClient(token)
	g.Catalog = make(map[string][]github.Repository)
	page := 1
	g.Username = Username
	for me.loadRepos(page) {
		page++
	}
}

func (g *Github) GetReposReadme(ctx context.Context, starChan chan<- *model.StarResult, token string, user string) {
// func GetReadme(token string, repoList []*GitHubRepo, j int, sendWg *sync.WaitGroup) {
	// log.Println("debug log:", j)
	// repo := *repoList[j]
	readmeURL := repoList[j].APIURL + "/readme"

	// log.Println("try to get readme:", readmeURL)
	req, err := http.NewRequest("GET", readmeURL, nil)
	if err != nil {
		log.Println("new request error :", err)
		// channel <- j
		sendWg.Done()

		return
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github.raw")
	c := http.Client{}
	res, err := c.Do(req)
	if err != nil {
		log.Println("res error to readme:", err)

		// channel <- j
		sendWg.Done()

		return
	}
	status := res.Header.Get("Status")
	if status == "404 Not Found" {
		log.Println("404 Not Found")
		// 	body: {"message":"Not Found","documentation_url":"https://developer.github.com/v3"}

	} else {
		b, err := ioutil.ReadAll(res.Body)
		b2 := ""
		_ = b2
		if err != nil {
			log.Println("read body error:", err)
		} else {
			b2 := string(b)
			repoList[j].Readme = b2
			// log.Println("got readme:", repoList[j].Readme)
		}
		res.Body.Close()
	}
	// log.Println("try to get readme done:", readmeURL)
	sendWg.Done()
}

func (g *Github) GetReposReadme(ctx context.Context, starChan chan<- *model.StarResult, token string, user string) {
// func GetReposReadme(token string, repoList []*GitHubRepo) error {
	lenList := len(repoList)
	log.Println("try getting all readme:", lenList)
	// c := make(chan int, lenList)
	// checkList := make([]int, lenList)
	var sendWg *sync.WaitGroup
	sendWg = new(sync.WaitGroup)
	for i := 0; i < lenList; i++ {
		sendWg.Add(1)
		go GetReadme(token, repoList, i, sendWg)
	}
	log.Println("start to wait")
	sendWg.Wait()
	log.Println("end to wait, after getting all readme")
	return nil
}
*/

/*
func getFileContents(client *github.Client, file, owner, repo string) ([]byte, error) {
  repoUrl := fmt.Sprintf("%v%v/%v/master/%v", GITHUB__RAW_URL, owner, repo, file)
  resp, err := http.Get(repoUrl)
  if resp.StatusCode != 200 {
    return nil, errors.New("Couldn't read file " + repoUrl)
  }
  content, err := ioutil.ReadAll(resp.Body)
  if err == nil {
    return content, nil
  } else {
	log.WithError(err).WithFields(logrus.Fields{"service": "getFileContents"}).Warnln("Couldn't get contents of file", file, " for ", owner, "/", repo, ": ", err)
    //glog.Errorln("Couldn't get contents of file", file, " for ", owner, "/", repo, ": ", error)
    return nil, err
  }
}

func (g *Github) getJsonFileContents(ctx context.Context, g *github.Client, file, owner, repo string, i interface{}) error {
//func getJsonFileContents(g *github.Client, file, owner, repo string, i interface{}) error {
  contents, err := g.getFileContents(client, file, owner, repo)
  if err == nil {
    err = json.Unmarshal(contents, &i)
    if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"service": "getJsonFileContents"}).Warnln("Couldn't decode json of ", file, " for ", owner, "/", repo, ": ", err)
      	//glog.Errorln("Couldn't decode json of ", file, " for ", owner, "/", repo, ": ", error)
      	return err
    }
  } else {
    return error
  }
  return nil
}


func (g *Github) checkResponse(ctx context.Context, g *github.Client) (r *github.Response) {
//func checkResponse(r *github.Response) {
}

func (g *Github) ExtractFileNames(ctx context.Context, owner string, repo string) ([]string, error) {
// func (g *Github) ExtractFileNames(ctx context.Context, filesChan chan<- *model.StarResult, owner string, repo string) ([]string, error) {
// func ExtractFileNames(g *github.Client, owner string, repo string) ([]string, error) {
  _, dir, response, err := g.Repositories.GetContents(owner, repo, "/", &github.RepositoryContentGetOptions{})
  g.checkResponse(response)
  if err != nil {
	log.WithError(err).WithFields(logrus.Fields{"service": "GetStars", "token": token}).Warnln("Couldn't get list of files for ", owner, "/", repo, ": ", err)
    return nil, err
  }
  fileNames := make([]string, len(dir))
  for i, file := range dir {
    fileNames[i] = *file.Name
  }
  return fileNames, nil
}

func (g *Github) ExtractGemspec(ctx context.Context, owner string, repo string, files []string, out chan string) {
// // func (g *Github) ExtractGemspec(ctx context.Context, genSpecChan chan<- *model.PackageJsonResult, owner string, repo string, files []string, out chan string) {
// func (g *Github) ExtractGemspec(ctx context.Context, owner string, repo string, files []string) {
// func ExtractGemspec(client *github.Client, owner string, repo string, files []string, out chan string) {
  file := repo + ".gemspec"
  if searchForFile(files, file) {
    content, err := g.getFileContents(g, file, owner, repo)
    contentS := string(content)
    if err != nil {
    }
      patterns := []*regexp.Regexp{
        regexp.MustCompile(`\.description\s*=\s*("|'|%q\{|%Q\{)(.*?)("|'|\})`),
        regexp.MustCompile(`\.name\s*=\s*"(|')(.*?)("|')`),
        regexp.MustCompile(`\.summary\s*=\s*("|'|%q\{|%Q\{)(.*?)("|'|\})`),
      }
      var result []string
      for _, regex := range patterns {
        result = regex.FindStringSubmatch(contentS)
        if len(result) > 1 {
          out <- result[2]
        }
      }
    } else {
		log.WithError(err).WithFields(logrus.Fields{"service": "ExtractGemspec"}).Warnln("Couldn't get list of files for ", owner, "/", repo, ": ", err)    	
    }
  } else {
	log.WithFields(logrus.Fields{"service": "ExtractGemspec"}).Warnln("could not find .gemspec file")  	
  }
  close(out)
}

func (g *Github) ExtractPackageJson(ctx context.Context, owner string, repo string, files []string, out chan string) {
// func (g *Github) ExtractPackageJson(ctx context.Context, packageJsonChan chan<- *model.PackageJsonResult, owner string, repo string, files []string, out chan string) {
// func (g *Github) ExtractPackageJson(ctx context.Context, owner string, repo string, files []string) {
// func ExtractPackageJson(g *github.Client, owner string, repo string, files []string, out chan string) {
  file := "package.json"
  if g.searchForFile(files, file) {
    var pack pjson
    err := g.getJsonFileContents(g, file, owner, repo, &pack)
    if err == nil {
      out <- pack.Name
      out <- pack.Description
      out <- strings.Join(pack.Keywords, " ")
    }
  } else {
	log.WithFields(logrus.Fields{"service": "ExtractPackageJson"}).Warnln("could not find package.json file")  	
  }
  close(out)
}


// func (g *Github) getClient(token string) *github.Client {
func (g *Github) ExtractRepoInfo(owner string, repo string) (*github.Repository, error) {
// func (g *Github) ExtractRepoInfo(ctx context.Context, repoInfoChan chan<- *model.RepoInfoResult, token, user string, page, count int) {
// func (g *Github) ExtractRepoInfo(ctx context.Context, owner string, repo string) {
// func ExtractRepoInfo(client *github.Client, owner string, repo string) (*github.Repository, error) {
  info, response, err := g.Repositories.Get(owner, repo)
  g.checkResponse(response)
  if err != nil {
	log.WithError(err).WithFields(logrus.Fields{"service": "ExtractPackageJson"}).Warnln("Coulnd't get repository info ", owner, "/", repo, ": ", err)  	
    return info, err
  } else {
    return info, nil
  }
}
*/

//func (g *Github) GetReadme(ctx context.Context, owner string, repo string, files []string, out chan string) {
/*
func (g *Github) GithubExtractReadme(ctx context.Context, readmeChan chan<- *model.ReadmeResult, owner string, repo string) ([]string, error) {
// func (g *Github) GithubExtractReadme(ctx context.Context, starChan chan<- *model.StarResult, token string, owner string, repo string, files []string) {
// func GithubExtractReadme(g *github.Client, owner string, repo string, files []string, out chan string) {
  readme, err := g.getFileContents(g, "README.md", owner, repo)
  if err != nil {
	log.WithError(err).WithFields(logrus.Fields{"service": "GetReadme"}).Warnln("Couldn't get readme for ", owner, "/", repo, ": ", err)
    close(out)
    return
  }
  if readme == nil {
	log.WithFields(logrus.Fields{"service": "GetReadme"}).Warnln("Content of readme is nil ", owner, "/", repo)
    close(out)
    return
  }
  out <- string(readme)
  close(out)
}
*/

/*

// ref doc: https://developer.github.com/v3/repos/contents/#get-contents

{
  "type": "file",
  "encoding": "base64",
  "size": 5362,
  "name": "README.md",
  "path": "README.md",
  "content": "encoded content ...",
  "sha": "3d21ec53a331a6f037a91c368710b99387d012c1",
  "url": "https://api.github.com/repos/octokit/octokit.rb/contents/README.md",
  "git_url": "https://api.github.com/repos/octokit/octokit.rb/git/blobs/3d21ec53a331a6f037a91c368710b99387d012c1",
  "html_url": "https://github.com/octokit/octokit.rb/blob/master/README.md",
  "download_url": "https://raw.githubusercontent.com/octokit/octokit.rb/master/README.md",
  "_links": {
    "git": "https://api.github.com/repos/octokit/octokit.rb/git/blobs/3d21ec53a331a6f037a91c368710b99387d012c1",
    "self": "https://api.github.com/repos/octokit/octokit.rb/contents/README.md",
    "html": "https://github.com/octokit/octokit.rb/blob/master/README.md"
  }
}

// ref doc: 


*/
