package service

import (
	"context"
	"crypto/md5"
	"fmt"
	"strings"
	"time"
	//"regexp"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"net/http"
	"path"
	//log "github.com/sirupsen/logrus"
	"github.com/google/go-github/github"
	"github.com/hoop33/entrevista"
	"github.com/roscopecoltran/sniperkit-limo/model"
	// "github.com/gregjones/httpcache"
	// "github.com/patrickmn/go-cache"
	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
	// cregex "github.com/mingrammer/commonregex"
	// "gopkg.in/libgit2/git2go.v26"
	// "github.com/sourcegraph/go-vcs/vcs/git"
	// "github.com/sourcegraph/go-vcs/vcs/gitcmd"
	// "github.com/parnurzeal/gorequest"
	// tablib "github.com/agrison/go-tablib"
	// "github.com/davecgh/go-spew/spew"
	// jsoniter "github.com/json-iterator/go"
	// fuzz "github.com/google/gofuzz"
)

// https://github.com/hfurubotten/autograder/blob/master/game/entities/repo.go
// https://github.com/hfurubotten/autograder/blob/master/game/entities/tokens.go
// https://github.com/hfurubotten/autograder/blob/master/game/entities/entities.go
// https://github.com/hfurubotten/autograder

/*
type GitHubUser struct {
	Account 			string 		`json:"account"`
	// AccessToken      string 		`json:"accessToken"`
	Tokens            	[]string 	`json:"tokens"`
	Status            	string   	`json:"status"`
	NumOfStarred      	int
	IndicesOfStarrerd 	int
}
*/

// Github represents the Github service
// GitHub holds specific information that is used for GitHub integration.
// Github represents the Github service
type Github struct {
}

//type Github struct {
//IgnoreRepos 	[]string 							`yaml:"ignore_repos,omitempty" json:"ignore_repos,omitempty"` // A list of URLs that the bot can ignore.
/*
	Catalog  		map[string][]github.Repository 		`yaml:"-" json:"-"`
	// Github API v3 - responses
	User 			*github.User 						`yaml:"-" json:"-"`						// https://github.com/google/go-github/blob/master/github/users.go#L20-L68
	Repo 			*github.Repository 		  			`yaml:"-" json:"-"`						// https://github.com/google/go-github/blob/master/github/repos.go#L22-L117
	Starred 		*github.StarredRepository 			`yaml:"-" json:"-"`						// https://github.com/google/go-github/blob/master/github/activity_star.go#L13-L17
	Readme 			*github.RepositoryContent 			`yaml:"-" json:"-"`						// https://github.com/google/go-github/blob/master/github/repos_contents.go#L22-L38
	Related  		map[string][]github.Repository 		`yaml:"-" json:"-"` 					// https://github.com/google/go-github/blob/master/github/repos.go#L22-L117
	Language 		map[string]int 						`yaml:"-" json:"-"` 					// https://github.com/google/go-github/blob/master/github/repos.go#L445-L469
	// search
	SearchOptions  	*github.SearchOptions 				`yaml:"-" json:"-"`						// https://github.com/google/go-github/blob/master/github/search.go#L29-L49
	SearchRepo 		*github.RepositoriesSearchResult 	`yaml:"-" json:"-"`						// https://github.com/google/go-github/blob/master/github/search.go#L52-L56
	SearchCode 		*github.CodeSearchResult 			`yaml:"-" json:"-"`						// https://github.com/google/go-github/blob/master/github/search.go#L149-L154
	SearchUsers 	*github.UsersSearchResult 			`yaml:"-" json:"-"`						// https://github.com/google/go-github/blob/master/github/search.go#L114-L119
	SearchIssues 	*github.IssuesSearchResult 			`yaml:"-" json:"-"`						// https://github.com/google/go-github/blob/master/github/search.go#L99-L102
	SearchCommits 	*github.CommitsSearchResult 		`yaml:"-" json:"-"`						// https://github.com/google/go-github/blob/master/github/search.go#L67-L72
*/
//OwnerAccount 	string 								`yaml:"-" json:"-"`
//Type     		string 								`yaml:"-" json:"-"`
//PerPage  		int 								`yaml:"-" json:"-"`
//}

// Login + OAuth2
// https://github.com/Jimdo/repos/blob/master/github.go#L33-L41
// https://github.com/Jimdo/repos/blob/master/main.go#L31-L35

/*
func (g *Github) LoginWithToken(ctx context.Context) (string, error) {
//func NewGitHubFetcherWithToken(token string) *GitHubFetcher {
	tc := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	))
	client := github.NewClient(tc)

	return &GitHubFetcher{client}
}
*/

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
		log.WithError(err).WithFields(
			logrus.Fields{
				"prefix":      "svc-github",
				"src.file":    "svc-vcs-github.go",
				"method.name": "(g *Github) Login(...)",
				"method.prev": "interview.Run(...)",
			}).Warn("asking the login credentials via the cli app.")
		return "", err
	}
	return answers["token"].(string), nil
}

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

func NewBiblio(client *github.Client) *Biblio {
	biblio := &Biblio{
		Cache:  make(map[string]*RepositoryInfo),
		Client: client,
	}
	if client == nil {
		biblio.Client = NewGithubClient()
	}

	return biblio
}

func (b *Biblio) listRepositoriesByOrg(org string) ([]*github.Repository, error) {
	allRepositories := make([]*github.Repository, 0)
	for i := 1; ; i++ {
		repositories, _, err := b.Client.Repositories.ListByOrg(context.Background(), org,
			&github.RepositoryListByOrgOptions{
				ListOptions: github.ListOptions{
					Page:    i,
					PerPage: 100,
				},
			},
		)
		if err != nil {
			return nil, err
		}
		if len(repositories) == 0 {
			break
		}
		allRepositories = append(allRepositories, repositories...)
	}
	return allRepositories, nil
}

func (b *Biblio) countNewOpenIssues(org, repo string, lastSyncedIssue int) (int, int, error) {
	newLastSyncedIssue := lastSyncedIssue
	count := 0
	var once sync.Once

	for i := 1; ; i++ {
		var issues []*github.Issue
		issues, _, err := b.Client.Issues.ListByRepo(context.Background(), org, repo, &github.IssueListByRepoOptions{
			ListOptions: github.ListOptions{
				Page:    i,
				PerPage: 100,
			},
		})
		if err != nil {
			return 0, 0, err
		}
		if len(issues) == 0 {
			break
		}

		for _, issue := range issues {
			if *issue.Number <= lastSyncedIssue {
				return count, newLastSyncedIssue, nil
			} else {
				once.Do(func() {
					newLastSyncedIssue = *issue.Number
				})
				count++
			}
		}
	}
	return count, newLastSyncedIssue, nil
}

func (b *Biblio) GetRepositoriesInfo(org string, repositoris ...string) (map[string]*RepositoryInfo, error) {
	allRepositories, err := b.getRepositories(org, repositoris...)
	if err != nil {
		return nil, err
	}

	cachedOrganizationReposInfo := b.Cache
	newOrganizationReposInfoMap := make(map[string]*RepositoryInfo)

	for _, repo := range allRepositories {
		repoName := ""
		if repo.Name != nil {
			repoName = *repo.Name
		}

		if repoName == "" {
			continue
		}

		cachedRepoInfo := cachedOrganizationReposInfo[repoName]
		repoInfo := new(RepositoryInfo)

		// Track Issues
		lastSyncedIssue := 0
		if cachedRepoInfo != nil {
			lastSyncedIssue = cachedRepoInfo.LastSyncedIssue.IssueNumber
		}
		count, issueNumber, err := b.countNewOpenIssues(org, repoName, lastSyncedIssue)
		if err != nil {
			return nil, err
		}
		repoInfo.LastSyncedIssue.IssueNumber = issueNumber
		repoInfo.LastSyncedIssue.Count = count

		// Track Stargazers
		users, err := b.getStargazers(org, repoName)
		if err != nil {
			return nil, err
		}
		repoInfo.Stargazers = users

		// Track
		if repo.ForksCount != nil {
			repoInfo.ForksCount = *repo.ForksCount
		}

		newOrganizationReposInfoMap[repoName] = repoInfo
	}
	return newOrganizationReposInfoMap, nil
}

func (b *Biblio) InitializeCache(org string, repositories ...string) error {
	repositoriesInfo, err := b.GetRepositoriesInfo(org, repositories...)
	if err != nil {
		return err
	}
	b.Cache = repositoriesInfo
	return nil
}
*/

/*

// gh notifications
// https://github.com/timakin/octop/blob/master/client/github.go

func downloadFile(tag string, outFileName string) {
	client := github.NewClient(nil)
	ctx := context.Background()
	opt := &github.RepositoryContentGetOptions{
		Ref: tag,
	}
	out, err := client.Repositories.DownloadContents(ctx, GITHUB_OWNER, GITHUB_REPO, "FINT-informasjonsmodell.xml", opt)
	if err != nil {
		fmt.Printf("Unable to download XMI file from GitHub: %s", err)
	}
	outFile, err := os.Create(outFileName)
	defer outFile.Close()
	_, err = io.Copy(outFile, out)
	if err != nil {
		fmt.Printf("Unable to write XMI file: %s", err)
	}
}

func getFilePath(tag string) string {
	homeDir, err := homedir.Dir()
	if err != nil {
		fmt.Println("Unable to get homedir.")
		os.Exit(2)
	}
	dir := fmt.Sprintf("%s/.fint-model/.cache", homeDir)
	err = os.MkdirAll(dir, 0777)

	if err != nil {
		fmt.Println("Unable to create .fint-model")
		os.Exit(2)
	}

	outFileName := fmt.Sprintf("%s/%s.xml", dir, tag)

	return outFileName
}

func toUtf8(fileName string) {
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error opening %s (%s)", fileName, err)
		os.Exit(2)
	}
	defer f.Close()

	r := charmap.Windows1252.NewDecoder().Reader(f)

	content, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(fileName, content, 0777)

	if err != nil {
		fmt.Println("\nio.Copy failed:", err)
	}

}

*/

func (g *Github) ParseFullName(fullName string) (owner string, name string, err error) {
	parts := strings.Split(fullName, "/")
	if len(parts) != 2 {
		err = fmt.Errorf("Invalid GitHub repository: %s", fullName)
		return
	}
	owner = parts[0]
	name = parts[1]
	return
}

func (g *Github) Stringify(message interface{}) string {
	var str string
	str = github.Stringify(message)
	str = strings.Replace(str, "\"", "", 2)
	return str
}

func (g *Github) GetRateLimit(ctx context.Context, token string) (int, int, error) {
	client := g.getClient(token)
	r, _, err := client.RateLimits(ctx)
	if err != nil {
		fmt.Errorf("Error getting core remaining: %v\n\n", err)
		return -1, -1, err
	}
	return r.Core.Remaining, r.Core.Limit, nil
}

// https://github.com/parnurzeal/gorequest
func (g *Github) GetLatestSHA(ctx context.Context, token string, user string, name string) (string, error) {
	client := g.getClient(token)
	// owner, name, err := g.ParseFullName(fullName)
	commits, _, err := client.Repositories.ListCommits(ctx, user, name, nil)
	if err != nil {
		return "", err
	}
	latestSHA := g.Stringify(commits[0].SHA)
	return latestSHA, nil
}

// return a list of Github Trees per SHA reqquested
func (g *Github) GetTrees(ctx context.Context, token string, user string, name string, shaList []string) ([]*github.Tree, error) {
	client := g.getClient(token)
	// owner, name, err := g.ParseFullName(fullName)
	if len(shaList) == 0 {
		shaList = []string{"master"}
	}
	var trees []*github.Tree
	for _, sha := range shaList {
		tree, _, errTree := client.Git.GetTree(ctx, user, name, sha, true)
		if errTree != nil {
			fmt.Printf("Git.GetTree returned error: %v", errTree)
			continue
		}
		trees = append(trees, tree)
	}
	return trees, nil
}

// returns only a list of filepaths per SHA requested
func (g *Github) GetTreesList(ctx context.Context, token string, user string, name string, shaList []string) (map[string][]string, error) {
	client := g.getClient(token)
	// owner, name, err := g.ParseFullName(fullName)
	if len(shaList) == 0 {
		shaList = []string{"master"}
	}
	trees := make(map[string][]string)
	for _, sha := range shaList {
		tree, _, errTree := client.Git.GetTree(ctx, user, name, sha, true)
		if errTree != nil {
			fmt.Printf("Git.GetTree returned error: %v", errTree)
			trees[sha] = make([]string, 0)
			// trees.Set(sha, make([]string, 0))
			continue
		}
		if len(tree.Entries) == 0 {
			continue
		}
		var paths []string
		for _, entry := range tree.Entries {
			paths = append(paths, *entry.Path) // parse paths for patterns
		}
		trees[sha] = paths
	}
	return trees, nil
}

func (g *Github) GetStargazers(ctx context.Context, token string, user string, name string) ([]string, error) {
	client := g.getClient(token)
	// owner, name, err := g.ParseFullName(fullName)
	resList := make([]string, 0)
	for i := 1; ; i++ {
		res, _, err := client.Activity.ListStargazers(context.Background(), user, name,
			&github.ListOptions{
				Page:    i,
				PerPage: 100,
			},
		)
		if err != nil {
			return nil, err
		}
		if len(res) == 0 {
			break
		}
		for _, s := range res {
			resList = append(resList, *(s.User.Login))
		}
	}
	return resList, nil
}

func (g *Github) GetReleases(ctx context.Context, token string, user string, name string) ([]string, error) {
	client := g.getClient(token)
	// owner, name, err := g.ParseFullName(fullName)
	var rels []string
	res, _, err := client.Repositories.ListTags(ctx, user, name, &github.ListOptions{})
	if err != nil {
		fmt.Printf("Unable to get tag list from GitHub: %s", err)
		return rels, err
	}
	for _, rel := range res {
		rels = append(rels, rel.GetName())
	}
	return rels, nil
}

func (g *Github) GetReleaseLatest(ctx context.Context, token string, user string, name string) (string, error) {
	client := g.getClient(token)
	// owner, name, err := g.ParseFullName(fullName)
	res, _, err := client.Repositories.GetLatestRelease(ctx, user, name)
	if err != nil {
		fmt.Printf("Unable to get latest release from GitHub: %s", err)
		return "", err
	}
	return res.GetTagName(), nil
}

func (g *Github) GetBranches(ctx context.Context, token string, user string, name string) ([]string, error) {
	client := g.getClient(token)
	// owner, name, err := g.ParseFullName(fullName)
	res, _, err := client.Repositories.ListBranches(ctx, user, name, &github.ListOptions{})
	if err != nil {
		fmt.Printf("Unable to get branch list from GitHub: %s", err)
		return make([]string, 0), err
	}
	if len(res) == 0 {
		return make([]string, 0), err
	}
	var resList []string
	for _, b := range res {
		resList = append(resList, b.GetName())
	}
	return resList, nil
}

func (g *Github) GetUserInfo(ctx context.Context, token string, user string) (*github.User, error) {
	client := g.getClient(token)
	res, _, err := client.Users.Get(ctx, user)
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{
				"prefix":      "svc-github",
				"method.name": "(g *Github) GetUserInfo(...)",
				"method.prev": "client.Users.Get(...)",
				"var.token":   token,
				"var.owner":   user,
			}).Warn("error while getting the content of the readme.")
		return &github.User{}, err
	}
	return res, nil
}

func (g *Github) GetLanguages(ctx context.Context, token string, user string, name string) (map[string]int, error) {
	client := g.getClient(token)
	// owner, name, err := g.ParseFullName(fullName)
	res, _, err := client.Repositories.ListLanguages(ctx, user, name)
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{
				"prefix":      "svc-github",
				"method.name": "(g *Github) GetLanguages(...)",
				"method.prev": "client.Repositories.ListLanguages(...)",
				"var.token":   token,
				"var.owner":   user,
				"var.repo":    name,
			}).Warn("error while getting the languages of the repo.")
		return make(map[string]int), err
	}
	if len(res) == 0 {
		res = make(map[string]int)
	}
	return res, nil
}

/*
  date_list := cregex.Date(text)
  // ['Jan 9th 2012']
  time_list := cregex.Time(text)
  // ['5:00PM', '4:00']
  link_list := cregex.Links(text)
  // ['www.linkedin.com', 'harold.smith@gmail.com']
  phone_list := cregex.PhonesWithExts(text)
  // ['(519)-236-2723x341']
  email_list := cregex.Emails(text)
*/
func (g *Github) GetReadme(ctx context.Context, token string, user string, name string) (string, error) {
	client := g.getClient(token)
	// owner, name, err := g.ParseFullName(fullName)
	res, _, err := client.Repositories.GetReadme(ctx, user, name, nil)
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{
				"prefix":      "svc-github-readme",
				"method.name": "(g *Github) GetReadme(...)",
				"method.prev": "client.Repositories.GetReadme(...)",
				"var.token":   token,
				"var.owner":   user,
				"var.repo":    name,
			}).Warn("error while getting the content of the readme.")
		return "", err
	}
	var resContent string = ""
	var errContent error
	if resContent, errContent = res.GetContent(); err != nil {
		log.WithError(errContent).WithFields(
			logrus.Fields{"method.name": "NewStarFromGithub(...)",
				"method.prev": "extraInfo.Readme.GetContent(...)",
				"var.token":   token,
				"var.owner":   user,
				"var.repo":    name,
			}).Error("extracting error on readme informations with readme.GetContent")
		return "", err
	}
	return resContent, nil
}

// GetStars returns the stars for the specified user (empty string for authenticated user)
func (g *Github) GetStars(ctx context.Context, starChan chan<- *model.StarResult, token string, user string, isAugmented bool) {
	log.WithFields(logrus.Fields{
		"prefix":          "svc-github-stars",
		"method.name":     "(g *Github) GetStars(...)",
		"method.next":     "g.getClient(...)",
		"var.token":       token,
		"var.isAugmented": isAugmented,
	}).Info("GetStars returns the stars for the specified user (empty string for authenticated user).")
	client := g.getClient(token) // Important: topics are requiring a change in the header sent by go-github, please append "application/vnd.github.mercy-preview+json"
	currentPage := 1             // The first response will give us the correct value for the last page
	lastPage := 1
	currentStar := 1
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
			log.WithError(err).WithFields(
				logrus.Fields{
					"prefix":          "svc-github-stars",
					"method.name":     "(g *Github) GetStars(...)",
					"method.prev":     "client.Activity.ListStarred(...)",
					"next":            "starChan <- &model.StarResult{}",
					"var.token":       token,
					"var.currentPage": currentPage,
					"var.lastPage":    lastPage,
				}).Warn("error while processing a list of starred repositories, let's forward it to the channel.")
			starChan <- &model.StarResult{
				Error:     err,
				Star:      nil,
				ExtraInfo: nil,
			}
		} else {
			lastPage = response.LastPage // Set last page only if we didn't get an error
			for _, repo := range repos { // Create a Star for each repository and put it on the channel
				var ownerName, repoUri, repoName string
				if *repo.Repository.Owner.Login != "" {
					ownerName = fmt.Sprintf("%s", *repo.Repository.Owner.Login)
				}
				if *repo.Repository.Name != "" {
					repoName = fmt.Sprintf("%s", *repo.Repository.Name) // string(*repo.Repository.Name)
				}
				repoUri = path.Join("github.com", ownerName, repoName)
				extraInfo := &model.GatewayBucket_GithubRepoExtraInfo{}

				// https://github.com/rafaeldias/async
				// Concurrent, Parallel, Waterfall

				if isAugmented {

					// https://github.com/kamildrazkiewicz/go-flow/blob/master/example/main.go
					// repository languages details returned by Github API v3
					languageInfo, err := g.GetLanguages(ctx, token, ownerName, repoName)
					if err != nil {
						log.WithError(err).WithFields(
							logrus.Fields{
								"prefix":            "svc-github-lang",
								"method.name":       "(g *Github) GetStars(...)",
								"method.prev":       "g.GetLanguages(...)",
								"var.repoUri":       repoUri,
								"var.star_owner":    ownerName,
								"var.star_owner_id": *repo.Repository.Owner.ID,
								"var.star_name":     repoName,
								"var.star_id":       *repo.Repository.ID,
							}).Warn("error while getting the readme content.")
						languageInfo = make(map[string]int)
					}
					// repository owner details
					userInfo, err := g.GetUserInfo(ctx, token, ownerName)
					if err != nil {
						log.WithError(err).WithFields(
							logrus.Fields{
								"prefix":            "svc-github-user",
								"method.name":       "(g *Github) GetStars(...)",
								"method.prev":       "g.GetUserInfo(...)",
								"var.star_owner":    ownerName,
								"var.star_owner_id": *repo.Repository.Owner.ID,
							}).Warn("error while getting additional info about the repository's owner.")
						userInfo = &github.User{}
					}

					readmeInfo, err := g.GetReadme(ctx, token, ownerName, repoName)
					if err != nil {
						log.WithError(err).WithFields(
							logrus.Fields{
								"prefix":            "svc-github-readme",
								"method.name":       "(g *Github) GetStars(...)",
								"method.prev":       "g.GetReadme(...)",
								"var.star_owner":    ownerName,
								"var.star_owner_id": *repo.Repository.Owner.ID,
								"var.star_name":     repoName,
								"var.star_id":       *repo.Repository.ID,
							}).Warn("error while getting the additional readme related info.")
					}

					branches, err := g.GetBranches(ctx, token, ownerName, repoName)
					if err != nil {
						fmt.Printf("g.GetBranches returned error: %v", err)
					}

					// https://github.com/pengwynn/flint/blob/master/flint/remote_project.go
					trees, err := g.GetTrees(ctx, token, ownerName, repoName, branches)
					if err != nil {
						fmt.Printf("g.GetTrees returned error: %v", err)
					}

					releases, err := g.GetReleases(ctx, token, ownerName, repoName)
					if err != nil {
						fmt.Printf("g.GetReleases returned error: %v", err)
					}

					latestRelease, err := g.GetReleaseLatest(ctx, token, ownerName, repoName)
					if err != nil {
						fmt.Printf("g.GetReleaseLatest returned error: %v", err)
					}

					stargazers, err := g.GetStargazers(ctx, token, ownerName, repoName)
					if err != nil {
						fmt.Printf("g.GetStargazers returned error: %v", err)
					}

					rateRemaining, rateLimit, err := g.GetRateLimit(ctx, token)
					if err != nil {
						fmt.Printf("g.GetRateLimit returned error: %v", err)
					}

					log.WithFields(logrus.Fields{
						"prefix":            "svc-github",
						"parent":            "(g *Github) GetStars(...)",
						"method.name":       "g.GetRateLimit(...)",
						"var.rateLimit":     rateLimit,
						"var.rateRemaining": rateRemaining,
					}).Info("get the api ratelimit...")

					extraInfo = &model.GatewayBucket_GithubRepoExtraInfo{
						UserInfo:      userInfo,
						Readme:        readmeInfo,
						Languages:     languageInfo,
						Branches:      branches,
						Releases:      releases,
						ReleaseLatest: latestRelease,
						Trees:         trees,
						Stargazers:    stargazers,
					}

				}

				star, err := model.NewStarFromGithub(repo.StarredAt, *repo.Repository, *extraInfo) // channels (default: 20)
				starChan <- &model.StarResult{
					Error:     err,
					Star:      star,
					ExtraInfo: extraInfo,
				}

				currentStar++
				log.WithFields(logrus.Fields{
					"prefix":          "svc-github",
					"parent":          "(g *Github) GetStars(...)",
					"method.name":     "starChan <- &model.StarResult{...}",
					"var.currentStar": currentStar,
					"var.currentPage": currentPage,
				}).Info("fetching content from new pages...")
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
	fetchedItemCount := 1
	for currentPage <= lastPage {
		events, _, err := client.Activity.ListEventsReceivedByUser(ctx, user, false, &github.ListOptions{
			Page: currentPage,
		})
		if err != nil {
			log.WithError(err).WithFields(
				logrus.Fields{
					"prefix":           "svc-github",
					"parent":           "(g *Github) GetEvents(...)",
					"method.name":      "client.Activity.ListEventsReceivedByUser(...)",
					"next":             "eventChan <- &model.EventResult{...}",
					"currentPage":      currentPage,
					"lastPage":         lastPage,
					"fetchedItemCount": fetchedItemCount,
				}).Warn("error while fetching additional events data/page info.")
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
				fetchedItemCount++
			}
		}
		log.WithFields(logrus.Fields{
			"prefix":           "svc-github",
			"parent":           "(g *Github) GetEvents(...)",
			"method.name":      "eventChan <- &model.EventResult{...}",
			"fetchedItemCount": fetchedItemCount,
			"currentPage":      currentPage,
		}).Info("fetching content from new pages...")
		currentPage++
	}
	close(eventChan)
}

// GetTrending returns the trending repositories
func (g *Github) GetTrending(ctx context.Context, trendingChan chan<- *model.StarResult, token string, language string, verbose bool) {
	client := g.getClient(token)
	log.WithFields(logrus.Fields{
		"prefix":      "svc-github",
		"parent":      "(g *Github) GetTrending(...)",
		"method.name": "g.getClient(...)",
		"next":        "g.getDateSearchString(...)",
		"token":       token,
		"language":    language,
		"verbose":     verbose,
	}).Info("returning the trending repositories...")
	// TODO perhaps allow them to specify multiple pages?
	// Might be overkill -- first page probably plenty
	// TODO Make this more configurable. Sort by stars, forks, default.
	// Search by number of stars, pushed, created, or whatever.
	// Lots of possibilities.
	q := g.getDateSearchString()
	if language != "" {
		q = fmt.Sprintf("language:%s %s", language, q)
		log.WithFields(logrus.Fields{
			"prefix":      "svc-github",
			"parent":      "(g *Github) GetTrending(...)",
			"method.name": "g.getDateSearchString(...)",
		}).Info("language is not empty...")
	}
	if verbose {
		// fmt.Println("q =", q)
		log.WithFields(logrus.Fields{
			"prefix":      "svc-github",
			"parent":      "(g *Github) GetTrending(...)",
			"method.name": "g.getDateSearchString(...)",
			"q":           q,
		}).Info("verbose mode")
	}
	fetchedItemCount := 1
	result, _, err := client.Search.Repositories(ctx, q, &github.SearchOptions{
		Sort:  "stars",
		Order: "desc",
	})
	// If we got an error, put it on the channel
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{
				"prefix":      "svc-github",
				"parent":      "(g *Github) GetTrending(...)",
				"method.name": "client.Search.Repositories(...)",
				"next":        "trendingChan <- &model.StarResult{...}",
			}).Warn("error while fetching additional trending infos.")
		trendingChan <- &model.StarResult{
			Error: err,
			Star:  nil,
		}
	} else {
		// Create a Star for each repository and put it on the channel
		for _, repo := range result.Repositories {
			star, err := model.NewStarFromGithub(nil, repo, model.GatewayBucket_GithubRepoExtraInfo{}) // add extra info from trending
			if err != nil {
				log.WithError(err).WithFields(
					logrus.Fields{
						"prefix":            "svc-github",
						"parent":            "(g *Github) GetTrending(...)",
						"method.name":       "model.NewStarFromGithub(...)",
						"next":              "trendingChan <- &model.StarResult{...}",
						"fetchedEventCount": fetchedItemCount,
					}).Warn("error while trying to register a new star from trendings.")
			}
			trendingChan <- &model.StarResult{
				Error: err,
				Star:  star,
			}
			fetchedItemCount++
		}
	}
	close(trendingChan)
}

// get CreatedAt from repo
// func (g *Github) getCreatedAtFromRepo(owner string, repo string) (createdAt time.Time, err error) {
func getCreatedAtFromRepo(ctx context.Context, client *github.Client, owner string, name string) (createdAt time.Time, err error) {
	repoInfo, _, err := client.Repositories.Get(ctx, owner, name)
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{
				"prefix":    "svc-github",
				"parent":    "getCreatedAtFromRepo(...)",
				"repoOwner": owner,
				"repoName":  name,
				// "repoInfo": 			repoInfo,
			}).Warnln("converting repository creation date...")
		return
	}
	var shortForm = "2006-01-02 15:04:05 -0700 UTC"
	ctime, _ := time.Parse(shortForm, fmt.Sprintf("%s", repoInfo.CreatedAt))
	return ctime, nil
}

func (g *Github) getDateSearchString() string {
	// TODO make this configurable
	// Default should be in configuration file
	// and should be able to override from command line
	// TODO should be able to specify whether "created" or "pushed"
	date := time.Now().Add(-7 * (24 * time.Hour))
	dateStr := fmt.Sprintf("created:>%s", date.Format("2006-01-02"))
	log.WithFields(logrus.Fields{
		"prefix":  "svc-github",
		"parent":  "getDateSearchString()",
		"date":    date,
		"dateStr": dateStr,
	}).Info("convert date to search str.")
	return dateStr
}

// https://github.com/timakin/octop/blob/master/client/interface.go
func (g *Github) getClient(token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	// https://github.com/aerokite/go-github-watcher
	// https://github.com/hairyhenderson/github-sync-labels-milestones/blob/master/sync/client.go
	//
	// https://github.com/tmthrgd/jekyll-history-service/blob/local/github-client.go
	// client := github.NewClient(transport.Client())
	//
	return github.NewClient(tc)
}

// https://github.com/HailoOSS/build-service/blob/master/githubrepo.go#L21
//
// NewGitHubClient - creates a client, authenticated by OAuth2 via a static token
func (g *Github) getClientCached(token string, cachePath string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	c := diskcache.New(cachePath)
	t := httpcache.NewTransport(c)
	hc := &http.Client{
		Transport: &oauth2.Transport{
			Base:   t,
			Source: ts,
		},
	}
	return github.NewClient(hc)
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
