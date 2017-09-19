package model

import (
	"errors"
	"fmt"
	"path"
	"strconv"
	"strings"
	"time"
	// "reflect"
	"github.com/blevesearch/bleve"
	"github.com/google/go-github/github"
	"github.com/jinzhu/gorm"
	//"github.com/qor/sorting"
	"github.com/xanzy/go-gitlab"
	// "github.com/amller/time/relative"
	jsoniter "github.com/json-iterator/go"
	// "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/skratchdot/open-golang/open"
	// el "github.com/src-d/enry"
	// rl "github.com/rai-project/linguist"
	// gl "github.com/generaltso/linguist"
	// jl "github.com/jhaynie/linguist"
	// tablib "github.com/agrison/go-tablib"
	//"github.com/davecgh/go-spew/spew"
	"github.com/olivere/nullable"
	// "github.com/k0kubun/pp"
	// "github.com/AlexSteele/deref"
	// fuzz "github.com/google/gofuzz"
	// "github.com/qor/qor"
	// "github.com/qor/admin"
	// es "gopkg.in/olivere/elastic.v5"
)

// https://github.com/rai-project/inle/blob/master/pkg/linguist/init.go
//
// PP_LINGUIST_URL=https://linguist:25032
// PP_LINGUIST_AUTH=1234
// result, err := jl.GetLanguageDetails(context.Background(), "test.js", []byte("var a = 1"))
// jl.AddExcludedFilename("foo.extension")
// jl.AddExcludedExtension(".extension")
// jl.linguist.NewMatcher("\\.somepath$")
// jl.AddExcludedRule(rule)
/*
files := []*linguist.File{
	linguist.NewFile("foo.properties", []byte("foo=1")),
	linguist.NewFile("foo.js", []byte("var foo=1")),
	linguist.NewFile("foo.jsx", []byte("var foo=1")),
}
results, err := linguist.GetLanguageDetailsMultiple(context.Background(), files)
*/

// https://github.com/redite/kleng/blob/master/core/gh.go
//
// https://github.com/alexchee/go-popular-repos/blob/master/cmd/go-popular-repos/main.go
// https://github.com/caarlos0/watchub
// https://github.com/timakin/octop/blob/master/client/sort_filter.go
// https://github.com/motemen/github-list-starred/blob/master/main.go
// https://github.com/kyokomi/gh-star-ranking/blob/master/github.go
// https://github.com/monmaru/ghstar/blob/master/ghstar.go
// https://github.com/glena/github-starred-catalog/blob/master/lib/ghclient.go
// https://github.com/kkeuning/gobservatory/blob/master/cmd/gobservatory/stars.go#L187
// https://github.com/Termina1/starlight/blob/master/handlers/repo_info_extractor.go
// https://github.com/Termina1/starlight/blob/master/star_extractor.go
// https://github.com/amoghe/polly/blob/master/frontman/github.go
// https://github.com/jeremyletang/amish/blob/master/main.go

// http://jinzhu.me/gorm/models.html#model-definition
// Star represents a starred repository
type Star struct {
	gorm.Model `json:"-" yaml:"-"`
	//sorting.SortingDESC
	RemoteID     string    `gorm:"column:remote_id" json:"remote_id,omitempty" yaml:"remote_id,omitempty"`
	OwnerID      string    `gorm:"column:owner_id" json:"owner_id,omitempty" yaml:"owner_id,omitempty"`
	RemoteURI    string    `gorm:"column:remote_uri" json:"remote_uri,omitempty" yaml:"remote_uri,omitempty" gorm:"type:varchar(128);not null;"`
	OwnerLogin   *string   `gorm:"column:owner_login" json:"owner_login,omitempty" yaml:"owner_login,omitempty" gorm:"type:varchar(128);not null;"`
	Name         *string   `gorm:"column:name" json:"name,omitempty" yaml:"name,omitempty" gorm:"type:varchar(128);not null;"`
	FullName     *string   `gorm:"column:full_name" json:"full_name,omitempty" yaml:"full_name,omitempty" gorm:"type:varchar(128);not null;"`
	Description  *string   `gorm:"column:description" json:"description,omitempty" yaml:"description,omitempty"`
	Homepage     *string   `gorm:"column:home_page" json:"home_page,omitempty" yaml:"home_page,omitempty"`
	URL          *string   `gorm:"column:url" json:"url,omitempty" yaml:"url,omitempty"`
	Language     *string   `gorm:"column:language" json:"language,omitempty" yaml:"language,omitempty"`
	Avatar       *string   `gorm:"column:avatar" json:"avatar,omitempty" yaml:"avatar,omitempty"`
	HasWiki      *bool     `gorm:"column:has_wiki" json:"has_wiki,omitempty" yaml:"has_wiki,omitempty"`
	Stargazers   int       `gorm:"column:stargazers_count" json:"stargazers_count,omitempty" yaml:"stargazers_count,omitempty"`
	Watchers     int       `gorm:"column:watchers_count" json:"watchers_count,omitempty" yaml:"watchers_count,omitempty"`
	Forks        int       `gorm:"column:forks_count" json:"forks_count,omitempty" yaml:"forks_count,omitempty"`
	StarredAt    time.Time `gorm:"column:starred_at" json:"starred_at" yaml:"starred_at"`
	LastUpdate   time.Time `gorm:"column:last_update" json:"last_update" yaml:"last_update"`
	CreationDate time.Time `gorm:"column:creation_date" json:"creation_date" yaml:"creation_date"`
	PushedAt     time.Time `gorm:"column:pushed_at" json:"pushed_at" yaml:"pushed_at"`
	ServiceID    uint      `gorm:"column:service_id" json:"service_id" yaml:"service_id"`
	// Extra
	Readme         string `gorm:"column:readme" json:"readme,omitempty" yaml:"readme,omitempty"`
	TopicsList     string `gorm:"column:topics" json:"topics,omitempty" yaml:"topics,omitempty"`
	BranchesList   string `gorm:"column:branches" json:"branches,omitempty" yaml:"branches,omitempty"`                   // gorm:"many2many:star_branches;"`
	ReleasesList   string `gorm:"column:releases" json:"releases,omitempty" yaml:"releases,omitempty"`                   // gorm:"many2many:star_releases;"`
	ReleaseLatest  string `gorm:"column:release_latest" json:"release_latest,omitempty" yaml:"release_latest,omitempty"` // gorm:"many2many:star_release_latest;"`
	StargazersList string `gorm:"column:stargazers" json:"stargazers,omitempty" yaml:"stargazers,omitempty"`             // gorm:"many2many:star_stargazers;"`

	UserInfo    User        `gorm:"many2many:star_users_info;" json:"user_info,omitempty" yaml:"user_info,omitempty"`             // *github.User  gorm:"many2many:star_users;"`
	UserInfoVCS UserInfoVCS `gorm:"many2many:star_users_info_vcs;" json:"user_info_vcs,omitempty" yaml:"user_info_vcs,omitempty"` // *github.User  gorm:"many2many:star_users;"`

	Topics    []Topic    `gorm:"many2many:star_topics;" json:"topics,omitempty" yaml:"topics,omitempty"`
	Trees     []Tree     `gorm:"many2many:star_trees;" json:"trees,omitempty" yaml:"trees,omitempty"` // gorm:"many2many:star_trees;"`
	Languages []Language `gorm:"many2many:star_languages;" json:"languages,omitempty" yaml:"languages,omitempty"`
	Tags      []Tag      `gorm:"many2many:star_tags;" json:"tags,omitempty" yaml:"tags,omitempty" gorm:"many2many:star_tags;"`
	//Detections   		[]Detection 				`gorm:"many2many:star_detection;" json:"languages_detected,omitempty" yaml:"languages_detected,omitempty"` // gorm:"many2many:star_languages_detected;"`
}

// https://github.com/GrantSeltzer/go-baseball-savant/blob/master/bbsavant/read_file.go
// StarResult wraps a star and an error
type StarResult struct {
	Star  *Star
	Error error
	//Cache  		map[string]*RepositoryInfo
	ExtraInfo *GatewayBucket_GithubRepoExtraInfo
}

// https://github.com/google/go-github/blob/master/github/repos.go#L21-L117
// NewStarFromGithub creates a Star from a Github repo
func NewStarFromGithub(timestamp *github.Timestamp, repo github.Repository, extraInfo GatewayBucket_GithubRepoExtraInfo) (*Star, error) {

	// Require the GitHub ID
	if repo.ID == nil {
		errMsg := errors.New("Repository ID from GitHub is required")
		log.WithError(errMsg).WithFields(
			logrus.Fields{
				"prefix":      "vcs-github-new-star",
				"method.name": "NewStarFromGithub(...)",
				"src.file":    "model/vcs-star.go)",
			}).Error("missing identifier...")
		return nil, errMsg
	}

	// Validate/filter repo data
	stargazersCount := nullable.IntWithDefault(repo.StargazersCount, 0) // Set 'stargazers' count to 0 if nil
	watchersCount := nullable.IntWithDefault(repo.WatchersCount, 0)     // Set 'watchers' count to 0 if nil
	forksCount := nullable.IntWithDefault(repo.ForksCount, 0)           // Set 'forks' count to 0 if nil

	starredAt := time.Now()
	if timestamp != nil {
		starredAt = timestamp.Time
	}

	ctime, _ := time.Parse(defaultDateShort, fmt.Sprintf("%s", repo.CreatedAt))

	// time.Time
	createdAt := repo.CreatedAt.Time
	updatedAt := repo.UpdatedAt.Time
	pushedAt := repo.PushedAt.Time

	// createdAt, _ := time.Parse(defaultDateShort, repo.CreatedAt.String()) // convert 'created_at' date to short format "2017-01-02 15:04:05 -0700 UTC"
	// updatedAt, _ := time.Parse(defaultDateShort, repo.UpdatedAt.String()) // convert 'updated_at' date to short format "2017-01-02 15:04:05 -0700 UTC"
	// pushedAt, _ := time.Parse(defaultDateShort, repo.PushedAt.String())   // convert 'pushed_at' date to short format "2017-01-02 15:04:05 -0700 UTC"

	remoteURI := path.Join(Default_VCS_Github_Domain, // register the remote URI in the database (without any protocol prefix, eg 'https://')
		fmt.Sprintf("%s", *repo.Owner.Login),
		fmt.Sprintf("%s", *repo.Name))

	// rt := relative.Convert(time.Now())

	log.WithFields(
		logrus.Fields{
			"src.file":                    "model/vcs-star.go",
			"prefix":                      "vcs-github-new-star",
			"method.name":                 "NewStarFromGithub(...)",
			"var.remoteURI":               remoteURI,
			"var.repo.createdAt":          createdAt,
			"var.ctime":                   ctime,
			"var.repo.CreatedAt.String()": repo.CreatedAt.String(),
			"var.repo.CreatedAtShort":     createdAt,
		}).Info("checking repo timeline information.")

	// 2016-04-04 15:19:46 +0000 UTC

	/*
			// https://stackoverflow.com/questions/18926303/iterate-through-a-struct-in-go
		    valExtraInfo 	:= reflect.ValueOf(extraInfo)
		    valuesExtraInfo := make([]interface{}, valExtraInfo.NumField())
			for i := 0; i < valExtraInfo.NumField(); i++ {
				valuesExtraInfo[i] 	= valExtraInfo.Field(i).Interface()
				pp.Println(valuesExtraInfo[i])
				// topics = append(topics, Topic{Name: t})
			}
	*/

	var topics []Topic
	if len(repo.Topics) > 0 {
		for _, t := range repo.Topics {
			topics = append(topics, Topic{Name: t})
		}
		log.WithFields(
			logrus.Fields{"method.name": "NewStarFromGithub(...)",
				"method.prev":   "extraInfo.Readme.GetContent(...)",
				"var.remoteURI": remoteURI,
				"var.topics":    strings.Join(repo.Topics, ","),
			}).Info("")
	}

	var languagesList []Language
	if len(extraInfo.Languages) > 0 {
		for langName, byteCode := range extraInfo.Languages {
			languagesList = append(languagesList, Language{Name: langName, ByteCode: byteCode})
		}
	}

	var treeBuckets []Tree
	if len(extraInfo.Trees) > 0 {
		for _, t := range extraInfo.Trees {
			var treeBucketEntries []TreeEntry
			for _, e := range t.Entries {
				treeEntry := TreeEntry{
					SHA:       e.SHA,
					FilePath:  e.Path,
					RemoteURL: e.URL,
					Type:      e.Type,
					Mode:      e.Mode,
					Size:      e.Size,
					Content:   e.Content,
				}
				treeBucketEntries = append(treeBucketEntries, treeEntry)
			}
			treeBucket := Tree{
				RemoteURI: remoteURI,
				SHA:       t.SHA,
				Entries:   treeBucketEntries,
			}
			treeBuckets = append(treeBuckets, treeBucket)
		}
	}

	userAccountEmail := nullable.StringWithDefault(extraInfo.UserInfo.Email, "hidden@github.com") // Set 'public_repos' count to 0 if nil
	collaboratorsCount := nullable.IntWithDefault(extraInfo.UserInfo.Collaborators, 0)            // Set 'collaborators' count to 0 if nil
	totalPrivateReposCount := nullable.IntWithDefault(extraInfo.UserInfo.TotalPrivateRepos, 0)    // Set 'collaborators' count to 0 if nil
	ownedPrivateReposCount := nullable.IntWithDefault(extraInfo.UserInfo.OwnedPrivateRepos, 0)    // Set 'collaborators' count to 0 if nil
	privateGistsCount := nullable.IntWithDefault(extraInfo.UserInfo.PrivateGists, 0)              // Set 'collaborators' count to 0 if nil

	publicReposCount := nullable.IntWithDefault(extraInfo.UserInfo.PublicRepos, 0) // Set 'collaborators' count to 0 if nil
	publicGistsCount := nullable.IntWithDefault(extraInfo.UserInfo.PublicGists, 0) // Set 'collaborators' count to 0 if nil
	followersCount := nullable.IntWithDefault(extraInfo.UserInfo.Followers, 0)     // Set 'collaborators' count to 0 if nil
	followingCount := nullable.IntWithDefault(extraInfo.UserInfo.Following, 0)     // Set 'collaborators' count to 0 if nil

	userGithubMetaInfo := &UserInfoVCS{
		UserID:            strconv.Itoa(*repo.Owner.ID), //strconv.Itoa(extraInfo.UserInfo.ID),
		Login:             extraInfo.UserInfo.Login,
		LoginEmail:        userAccountEmail,
		PublicRepos:       publicReposCount,
		PublicGists:       publicGistsCount,
		Followers:         followersCount,
		Following:         followingCount,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
		AccountType:       extraInfo.UserInfo.Type,
		SiteAdmin:         extraInfo.UserInfo.SiteAdmin,
		TotalPrivateRepos: totalPrivateReposCount,
		OwnedPrivateRepos: ownedPrivateReposCount,
		PrivateGists:      privateGistsCount,
		DiskUsage:         extraInfo.UserInfo.DiskUsage,
		Collaborators:     collaboratorsCount,
		ProgLanguages:     languagesList,
		AvatarURL:         extraInfo.UserInfo.AvatarURL,
	}

	userMetaInfo := &User{
		HTMLURL:    extraInfo.UserInfo.HTMLURL,
		GravatarID: extraInfo.UserInfo.GravatarID,
		Name:       extraInfo.UserInfo.Name,
		Company:    extraInfo.UserInfo.Company,
		Blog:       extraInfo.UserInfo.Blog,
		Location:   extraInfo.UserInfo.Location,
		Email:      userAccountEmail,
		Hireable:   extraInfo.UserInfo.Hireable,
		Bio:        extraInfo.UserInfo.Bio,
		Vcs:        *userGithubMetaInfo,
	}

	repoMetaInfo := &Star{
		RemoteID:       strconv.Itoa(*repo.ID),
		OwnerID:        strconv.Itoa(*repo.Owner.ID),
		RemoteURI:      remoteURI,
		OwnerLogin:     repo.Owner.Login,
		Name:           repo.Name,
		FullName:       repo.FullName,
		Description:    repo.Description,
		Homepage:       repo.Homepage,
		URL:            repo.CloneURL,
		Language:       repo.Language,
		Avatar:         repo.Owner.AvatarURL,
		HasWiki:        repo.HasWiki,
		Stargazers:     stargazersCount,
		Watchers:       watchersCount,
		Forks:          forksCount,
		StarredAt:      starredAt,
		LastUpdate:     updatedAt,
		CreationDate:   createdAt,
		PushedAt:       pushedAt,
		StargazersList: strings.Join(extraInfo.Stargazers, ","),
		BranchesList:   strings.Join(extraInfo.Branches, ","),
		ReleasesList:   strings.Join(extraInfo.Releases, ","),
		ReleaseLatest:  extraInfo.ReleaseLatest,
		Readme:         extraInfo.Readme,
		Trees:          treeBuckets,
		Languages:      languagesList,
		Topics:         topics,
		TopicsList:     strings.Join(repo.Topics, ","),
		UserInfo:       *userMetaInfo,
		UserInfoVCS:    *userGithubMetaInfo,
	}

	//pp.Print(repo)
	//pp.Print(repoMetaInfo)
	dumpPrefixPath := path.Join("cache", "api", "github.com", *repo.FullName)

	dumpRepoInfo, err := jsoniter.Marshal(repo)
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{"method.name": "NewStarFromGithub(...)",
				"method.prev":        "jsoniter.Marshal(repo)",
				"var.dumpPrefixPath": dumpPrefixPath,
				"var.repo":           repo,
			}).Warn("dump error on readme informations received with jsoniter")
	} else {
		if err := NewDump([]byte(fmt.Sprintf("[%s]\n", dumpRepoInfo)), dumpPrefixPath, "repository", []string{"json", "yaml"}); err != nil {
			return nil, errors.New("Could not dump the data to file for all export format.")
		}
	}

	dumpStar, err := jsoniter.Marshal(repoMetaInfo)
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{
				"method.name":        "NewStarFromGithub(...)",
				"method.prev":        "jsoniter.Marshal(repoMetaInfo)",
				"var.dumpPrefixPath": dumpPrefixPath,
				"var.dumpStar":       dumpStar,
			}).Warn("dump error on repo starred with jsoniter")
	} else {
		dumpPrefixPath := path.Join("cache", "api", "github.com", *repo.FullName)
		if err := NewDump([]byte(fmt.Sprintf("[%s]\n", dumpStar)), dumpPrefixPath, "meta", []string{"json", "yaml"}); err != nil {
			return nil, errors.New("Could not dump the data to file for all export format.")
		}
	}

	dumpUser, err := jsoniter.Marshal(userMetaInfo)
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{
				"method.name":        "NewStarFromGithub(...)",
				"method.prev":        "jsoniter.Marshal(userMetaInfo)",
				"var.dumpPrefixPath": dumpPrefixPath,
				"var.dumpUser":       dumpUser,
			}).Warn("dump error on repo starred with jsoniter")
	} else {
		dumpPrefixPath := path.Join("cache", "api", "github.com", *repo.Owner.Login)
		if err := NewDump([]byte(fmt.Sprintf("[%s]\n", dumpUser)), dumpPrefixPath, "user", []string{"json", "yaml"}); err != nil {
			return nil, errors.New("Could not dump the data to file for all export format.")
		}
	}

	return repoMetaInfo, nil

}

// ref. https://github.com/xanzy/go-gitlab/blob/master/projects.go#L33-L175
// NewStarFromGitlab creates a Star from a Gitlab star
func NewStarFromGitlab(star gitlab.Project) (*Star, error) {
	/*
		// Set stargazers count to 0 if nil
		stargazersCount := 0
		if star.StargazersCount != nil {
			stargazersCount = *star.StargazersCount
		}

		// Set stargazers count to 0 if nil
		watchersCount := 0
		if star.WatchersCount != nil {
			watchersCount = *star.WatchersCount
		}
	*/
	return &Star{
		RemoteID:    strconv.Itoa(star.ID),
		Name:        &star.Name,
		FullName:    &star.NameWithNamespace,
		Description: &star.Description,
		Homepage:    &star.WebURL,
		URL:         &star.HTTPURLToRepo,
		Language:    nil,
		// Topics:      	   nil,
		// LanguagesDetected:  nil,
		Avatar:     &star.AvatarURL,
		HasWiki:    &star.WikiEnabled,
		Stargazers: star.StarCount,
		Forks:      star.ForksCount,
		//ForkedFromProject:  star.ForkedFromProject.PathWithNamespace,
		// Snippets: repo.SnippetsEnabled,
		StarredAt: time.Now(), // OK, so this is a lie, but not in payload
	}, nil

}

// not ready yet !
func DumpStarInfo(db *gorm.DB, star *Star, service *Service) (bool, error) {
	return false, db.Save(star).Error
}

// CreateOrUpdateStar creates or updates a star and returns true if the star was created (vs updated)
func CreateOrUpdateStar(db *gorm.DB, star *Star, service *Service) (bool, error) {
	// Get existing by remote ID and service ID
	var existing Star
	if db.Where("remote_id = ? AND service_id = ?", star.RemoteID, service.ID).First(&existing).RecordNotFound() {
		star.ServiceID = service.ID
		err := db.Create(star).Error
		return err == nil, err
	}
	star.ID = existing.ID
	star.ServiceID = service.ID
	star.CreatedAt = existing.CreatedAt
	return false, db.Save(star).Error
}

// FindStarByID finds a star by ID
func FindStarByID(db *gorm.DB, ID uint) (*Star, error) {
	var star Star
	if db.First(&star, ID).RecordNotFound() {
		return nil, fmt.Errorf("Star '%d' not found", ID)
	}
	return &star, db.Error
}

// FindStars finds all stars
func FindStars(db *gorm.DB, match string) ([]Star, error) {
	var stars []Star
	if match != "" {
		db.Where("full_name LIKE ?",
			strings.ToLower(fmt.Sprintf("%%%s%%", match))).Order("full_name").Find(&stars)
	} else {
		db.Order("full_name").Find(&stars)
	}
	return stars, db.Error
}

// FindUntaggedStars finds stars without any tags
func FindUntaggedStars(db *gorm.DB, match string) ([]Star, error) {
	var stars []Star
	if match != "" {
		db.Raw(`
			SELECT *
			FROM STARS S
			WHERE S.DELETED_AT IS NULL
			AND S.FULL_NAME LIKE ?
			AND S.ID NOT IN (
				SELECT STAR_ID
				FROM STAR_TAGS
			) ORDER BY S.FULL_NAME`,
			fmt.Sprintf("%%%s%%", strings.ToLower(match))).Scan(&stars)
	} else {
		db.Raw(`
			SELECT *
			FROM STARS S
			WHERE S.DELETED_AT IS NULL
			AND S.ID NOT IN (
				SELECT STAR_ID
				FROM STAR_TAGS
			) ORDER BY S.FULL_NAME`).Scan(&stars)
	}
	return stars, db.Error
}

// FindStarsByLanguageAndOrTag finds stars with the specified language and/or the specified tag
func FindStarsByLanguageAndOrTag(db *gorm.DB, match string, language string, tagName string, union bool) ([]Star, error) {
	operator := "AND"
	if union {
		operator = "OR"
	}

	var stars []Star
	if match != "" {
		db.Raw(fmt.Sprintf(`
			SELECT * 
			FROM STARS S, TAGS T 
			INNER JOIN STAR_TAGS ST ON S.ID = ST.STAR_ID 
			INNER JOIN TAGS ON ST.TAG_ID = T.ID 
			WHERE S.DELETED_AT IS NULL
			AND T.DELETED_AT IS NULL
			AND LOWER(S.FULL_NAME) LIKE ? 
			AND (LOWER(T.NAME) = ? 
			%s LOWER(S.LANGUAGE) = ?) 
			GROUP BY ST.STAR_ID 
			ORDER BY S.FULL_NAME`, operator),
			fmt.Sprintf("%%%s%%", strings.ToLower(match)),
			strings.ToLower(tagName),
			strings.ToLower(language)).Scan(&stars)
	} else {
		db.Raw(fmt.Sprintf(`
			SELECT * 
			FROM STARS S, TAGS T 
			INNER JOIN STAR_TAGS ST ON S.ID = ST.STAR_ID 
			INNER JOIN TAGS ON ST.TAG_ID = T.ID 
			WHERE S.DELETED_AT IS NULL
			AND T.DELETED_AT IS NULL
			AND LOWER(T.NAME) = ? 
			%s LOWER(S.LANGUAGE) = ? 
			GROUP BY ST.STAR_ID 
			ORDER BY S.FULL_NAME`, operator),
			strings.ToLower(tagName),
			strings.ToLower(language)).Scan(&stars)
	}
	return stars, db.Error
}

// FindStarsByLanguage finds stars with the specified language
func FindStarsByLanguage(db *gorm.DB, match string, language string) ([]Star, error) {
	var stars []Star
	if match != "" {
		db.Where("full_name LIKE ? AND lower(language) = ?",
			strings.ToLower(fmt.Sprintf("%%%s%%", match)),
			strings.ToLower(language)).Order("full_name").Find(&stars)
	} else {
		db.Where("lower(language) = ?",
			strings.ToLower(language)).Order("full_name").Find(&stars)
	}
	return stars, db.Error
}

// FuzzyFindStarsByName finds stars with approximate matching for full name and name
func FuzzyFindStarsByName(db *gorm.DB, name string) ([]Star, error) {
	// Try each of these, and as soon as we hit, return
	// 1. Exact match full name
	// 2. Exact match name
	// 3. Case-insensitive full name
	// 4. Case-insensitive name
	// 5. Case-insensitive like full name
	// 6. Case-insensitive like name
	var stars []Star
	db.Where("full_name = ?", name).Order("full_name").Find(&stars)
	if len(stars) == 0 {
		db.Where("name = ?", name).Order("full_name").Find(&stars)
	}
	if len(stars) == 0 {
		db.Where("lower(full_name) = ?", strings.ToLower(name)).Order("full_name").Find(&stars)
	}
	if len(stars) == 0 {
		db.Where("lower(name) = ?", strings.ToLower(name)).Order("full_name").Find(&stars)
	}
	if len(stars) == 0 {
		db.Where("full_name LIKE ?", strings.ToLower(fmt.Sprintf("%%%s%%", name))).Order("full_name").Find(&stars)
	}
	if len(stars) == 0 {
		db.Where("name LIKE ?", strings.ToLower(fmt.Sprintf("%%%s%%", name))).Order("full_name").Find(&stars)
	}
	return stars, db.Error
}

// FindLanguages finds all languages
func FindLanguages(db *gorm.DB) ([]string, error) {
	var languages []string
	db.Table("stars").Order("language").Pluck("distinct(language)", &languages)
	return languages, db.Error
}

// AddTag adds a tag to a star
func (star *Star) AddTag(db *gorm.DB, tag *Tag) error {
	star.Tags = append(star.Tags, *tag)
	return db.Save(star).Error
}

// LoadTags loads the tags for a star
func (star *Star) LoadTags(db *gorm.DB) error {
	// Make sure star exists in database, or we will panic
	var existing Star
	if db.Where("id = ?", star.ID).First(&existing).RecordNotFound() {
		return fmt.Errorf("Star '%d' not found", star.ID)
	}
	return db.Model(star).Association("Tags").Find(&star.Tags).Error
}

// LoadTags loads the tags for a star
func (star *Star) LoadTopics(db *gorm.DB) error {
	// Make sure star exists in database, or we will panic
	var existing Star
	if db.Where("id = ?", star.ID).First(&existing).RecordNotFound() {
		return fmt.Errorf("Star '%d' not found", star.ID)
	}
	return db.Model(star).Association("Topics").Find(&star.Topics).Error
}

// RemoveAllTags removes all tags for a star
func (star *Star) RemoveAllTags(db *gorm.DB) error {
	return db.Model(star).Association("Tags").Clear().Error
}

// RemoveTag removes a tag from a star
func (star *Star) RemoveTag(db *gorm.DB, tag *Tag) error {
	return db.Model(star).Association("Tags").Delete(tag).Error
}

// HasTag returns whether a star has a tag. Note that you must call LoadTags first -- no reason to incur a database call each time
func (star *Star) HasTag(tag *Tag) bool {
	if len(star.Tags) > 0 {
		for _, t := range star.Tags {
			if t.Name == tag.Name {
				return true
			}
		}
	}
	return false
}

// Index adds the star to the index
func (star *Star) Index(index bleve.Index, db *gorm.DB) error {
	if err := star.LoadTags(db); err != nil {
		return err
	}
	return index.Index(fmt.Sprintf("%d", star.ID), star)
}

// OpenInBrowser opens the star in the browser
func (star *Star) OpenInBrowser(preferHomepage bool) error {
	var URL string
	if preferHomepage && star.Homepage != nil && *star.Homepage != "" {
		URL = *star.Homepage
	} else if star.URL != nil && *star.URL != "" {
		URL = *star.URL
	} else {
		if star.Name != nil {
			return fmt.Errorf("No URL for star '%s'", *star.Name)
		}
		return errors.New("No URL for star")
	}
	return open.Start(URL)
}
