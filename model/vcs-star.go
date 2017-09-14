package model

import (
	// golang
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"path"
	"reflect"
	// search engine
	"github.com/blevesearch/bleve"
	// vcs api wrappers
	"github.com/google/go-github/github"
	"github.com/xanzy/go-gitlab"
	// database
	"github.com/jinzhu/gorm"
	// data processing
	jsoniter "github.com/json-iterator/go"
	// logs
	"github.com/sirupsen/logrus"
	// widgets
	"github.com/skratchdot/open-golang/open"
	// el "github.com/src-d/enry"
	// rl "github.com/rai-project/linguist"
	// gl "github.com/generaltso/linguist"
	// jl "github.com/jhaynie/linguist"

	// tablib "github.com/agrison/go-tablib"
	"github.com/davecgh/go-spew/spew"
	"github.com/k0kubun/pp"
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

// http://jinzhu.me/gorm/models.html#model-definition
// Star represents a starred repository
type Star struct {
	gorm.Model
	RemoteID    		string
	OwnerID        		string 		
	RemoteURI        	string 			`gorm:"type:varchar(128);not null;"`
	OwnerLogin        	*string 		`gorm:"type:varchar(128);not null;"`
	Name        		*string 		`gorm:"type:varchar(128);not null;"`
	FullName    		*string 		`gorm:"type:varchar(128);not null;"`
	Description 		*string
	Homepage    		*string
	//RemoteURI         	string
	URL         		*string
	Language    		*string
	Avatar				*string
	HasWiki 			*bool
	// Readme       		string 		`json:"readme"`
	// SHA  			*string
	// ForkedFromProject *string
	// Snippets 		*bool
	// Topics    			*[]string
	Stargazers  		int
	Watchers  			int
	Forks 				int

	StarredAt   		time.Time
	LastUpdate   		time.Time
	CreationData   		time.Time
	PushedAt 			time.Time

	ServiceID   		uint 				// `gorm:"index:idx_name_code"`
	ReadmeDoc        	string

	//Admins    		map[string]interface{}

	// Extra
	TopicsList 			string
	Topics      		[]Topic 					`gorm:"many2many:star_topics;"`
	User      			*github.User 				`gorm:"many2many:star_users;"`
	Tree      			[]GatewayBucket_GithubTree 	`gorm:"many2many:star_trees;"`
	Tags        		[]Tag 						`gorm:"many2many:star_tags;"`
	Languages      		map[string]int 				`gorm:"many2many:star_languages;"`
	Readme      		*github.RepositoryContent 	`gorm:"many2many:star_readmes;"`
	LanguagesDetected   []LanguageDetected 			`gorm:"many2many:star_languages_detected;"`
}

// https://github.com/GrantSeltzer/go-baseball-savant/blob/master/bbsavant/read_file.go
// StarResult wraps a star and an error
type StarResult struct {
	Star  			*Star
	ExtraInfo 		*GatewayBucket_GithubRepoExtraInfo
	//Cache  		map[string]*RepositoryInfo
	Error 			error
}

// https://github.com/google/go-github/blob/master/github/repos.go#L21-L117
// NewStarFromGithub creates a Star from a Github star
func NewStarFromGithub(timestamp *github.Timestamp, star github.Repository, extraInfo GatewayBucket_GithubRepoExtraInfo) (*Star, error) {
	// Require the GitHub ID
	if star.ID == nil {
		log.WithFields(
			logrus.Fields{"model": "NewStarFromGithub"}).Warn("ID from GitHub is required")
		return nil, errors.New("ID from GitHub is required")
	}

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

	// Set stargazers count to 0 if nil
	forksCount := 0
	if star.ForksCount != nil {
		forksCount = *star.ForksCount
	}

	starredAt := time.Now()
	if timestamp != nil {
		starredAt = timestamp.Time
	}

	var shortForm 	= "2006-01-02 15:04:05 -0700 UTC"
	createdAt, _ 	:= time.Parse(shortForm, fmt.Sprintf("%s", *star.CreatedAt))
	updatedAt, _ 	:= time.Parse(shortForm, fmt.Sprintf("%s", *star.UpdatedAt))
	pushedAt, _ 	:= time.Parse(shortForm, fmt.Sprintf("%s", *star.PushedAt))

	starUri 		:= path.Join("github.com", fmt.Sprintf("%s", *star.Owner.Login), fmt.Sprintf("%s", *star.Name))

	log.WithFields(
		logrus.Fields{	"starUri": 		starUri, 
						"createdAt": 	createdAt, 
						"updatedAt": 	updatedAt}).Info("")	

	pp.Println(star)

	// https://stackoverflow.com/questions/18926303/iterate-through-a-struct-in-go
    valExtraInfo 	:= reflect.ValueOf(extraInfo)
    valuesExtraInfo := make([]interface{}, valExtraInfo.NumField())

	//for _, info := range valExtraInfo.NumField() {
	for i := 0; i < valExtraInfo.NumField(); i++ {
		valuesExtraInfo[i] 	= valExtraInfo.Field(i).Interface()
		pp.Println(valuesExtraInfo[i])
		// topics = append(topics, Topic{Name: t})
	}

	fmt.Println(valuesExtraInfo)

	var readmeContent string
	var err error
	if *extraInfo.Readme.Content != "" {
		if readmeContent, err = extraInfo.Readme.GetContent(); err != nil {
			log.WithError(err).WithFields(
				logrus.Fields{	"action": 		"NewReadmeFromGithub", 
								"step": 		"readmeGetContent", 
								"userId": 		*star.Owner.ID, 
								"remoteId": 	*star.ID, 
								"remoteUri": 	starUri}).Warn("extracting error on readme informations with readme.GetContent")
			return nil, err
		}
	}

	var topics []Topic
	if len(star.Topics) > 0 {
		for _, t := range star.Topics {
			topics = append(topics, Topic{Name: t})
		}
		log.WithFields(
			logrus.Fields{	"starUri": 	starUri, 
							"topics": 	strings.Join(star.Topics, ",")}).Warn("")	
	}

	starMetaInfo :=	&Star{
		RemoteID:    	strconv.Itoa(*star.ID),
		OwnerID:	 	strconv.Itoa(*star.Owner.ID),
		RemoteURI:	 	starUri,
		OwnerLogin:  	star.Owner.Login,
		Name:        	star.Name,
		FullName:    	star.FullName,
		Description: 	star.Description,
		Homepage:    	star.Homepage,
		URL:         	star.CloneURL,
		Language:    	star.Language,
		Avatar: 	 	star.Owner.AvatarURL,
		HasWiki: 	 	star.HasWiki,
		Stargazers:  	stargazersCount,
		Watchers:  	 	watchersCount,
		Forks:  	 	forksCount,

		StarredAt:    	starredAt,
		LastUpdate:   	updatedAt,
		CreationData: 	createdAt,
		PushedAt: 	  	pushedAt,

		// Extra
		User: 	 		extraInfo.User,
		Readme: 	 	extraInfo.Readme,
		ReadmeDoc: 	 	readmeContent,
		Languages: 	 	extraInfo.Languages,
		Tree: 			extraInfo.Trees,
		Topics:    	 	topics,
		TopicsList:  	strings.Join(star.Topics, ","),

	}

	spew.Dump(extraInfo)

	dumpPrefixPath  := path.Join("cache", "vcs", "github.com", *star.FullName)

	dumpRepoInfo, err := jsoniter.Marshal(star)
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"action": 			"NewReadmeFromGithub", 
							"step": 			"JsoniterMarshalRepoInfo", 
							"dumpPrefixPath": 	dumpPrefixPath}).Warn("dump error on readme informations received with jsoniter")
	} else {
		if err := NewDump([]byte(fmt.Sprintf("[%s]\n", dumpRepoInfo)), dumpPrefixPath, "repository", []string{"json", "yaml"}); err != nil {
			return nil, errors.New("Could not dump the data to file for all export format.")
		}
	}

	dumpStar, err := jsoniter.Marshal(starMetaInfo)
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"service": "GetStars", 
							"dumpPrefixPath": dumpPrefixPath}).Warn("dump error on repo starred with jsoniter")
	} else {
		dumpPrefixPath  := path.Join("cache", "vcs", "github.com", *star.FullName)
		if err := NewDump([]byte(fmt.Sprintf("[%s]\n", dumpStar)), dumpPrefixPath, "meta", []string{"json", "yaml"}); err != nil {
			return nil, errors.New("Could not dump the data to file for all export format.")
		}
	}

	return starMetaInfo, nil
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
		RemoteID:    		strconv.Itoa(star.ID),
		Name:        		&star.Name,
		FullName:    		&star.NameWithNamespace,
		Description: 		&star.Description,
		Homepage:    		&star.WebURL,
		URL:         		&star.HTTPURLToRepo,
		Language:    		nil,
		// Topics:      	   nil,
		// LanguagesDetected:  nil,
		Avatar: 	 		&star.AvatarURL,
		HasWiki:			&star.WikiEnabled,
		Stargazers:  		star.StarCount,
		Forks:  	 		star.ForksCount,
		//ForkedFromProject:  star.ForkedFromProject.PathWithNamespace,
		// Snippets: repo.SnippetsEnabled,
		StarredAt:   		time.Now(), // OK, so this is a lie, but not in payload
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



