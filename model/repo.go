package model

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/google/go-github/github"
	"github.com/jinzhu/gorm"
	"github.com/skratchdot/open-golang/open"
	"github.com/xanzy/go-gitlab"
	// tablib "github.com/agrison/go-tablib"
    // "github.com/qor/qor"
    // "github.com/qor/admin"

)

// https://github.com/Termina1/repolight/blob/master/handlers/repo_info_extractor.go
// https://github.com/Termina1/repolight/blob/master/repo_extractor.go

// Repo represents a repored repository
type Repo struct {
	gorm.Model
	RemoteID    		string
	Name        		*string
	FullName    		*string
	Description 		*string
	Homepage    		*string
	URL         		*string
	Language    		*string
	Repogazers  		int
	ReporedAt   		time.Time
	ServiceID   		uint
	//UserName        		*string
	Tags        		[]Tag `gorm:"many2many:repo_tags;"`
	Topics      		[]Topic `gorm:"many2many:repo_topics;"`
	LanguagesDetected   []LanguageDetected `gorm:"many2many:repo_languages;"`
}

// RepoResult wraps a repo and an error
type RepoResult struct {
	Repo  *Repo
	Error error
}

// NewRepoFromGithub creates a Repo from a Github repo
func NewRepoFromGithub(timestamp *github.Timestamp, repo github.Repository) (*Repo, error) {
	// Require the GitHub ID
	if repo.ID == nil {
		return nil, errors.New("ID from GitHub is required")
	}

	// Set repogazers count to 0 if nil
	repogazersCount := 0
	if repo.RepogazersCount != nil {
		repogazersCount = *repo.RepogazersCount
	}

	reporedAt := time.Now()
	if timestamp != nil {
		reporedAt = timestamp.Time
	}

	return &Repo{
		RemoteID:    strconv.Itoa(*repo.ID),
		Name:        repo.Name,
		FullName:    repo.FullName,
		Description: repo.Description,
		Homepage:    repo.Homepage,
		URL:         repo.CloneURL,
		Language:    repo.Language,
		Repogazers:  repogazersCount,
		ReporedAt:   reporedAt,
		//Topics:   topics,
	}, nil
}

// NewRepoFromGitlab creates a Repo from a Gitlab repo
func NewRepoFromGitlab(repo gitlab.Project) (*Repo, error) {
	return &Repo{
		RemoteID:    strconv.Itoa(repo.ID),
		Name:        &repo.Name,
		FullName:    &repo.NameWithNamespace,
		Description: &repo.Description,
		Homepage:    &repo.WebURL,
		URL:         &repo.HTTPURLToRepo,
		Language:    nil,
		Repogazers:  repo.RepoCount,
		ReporedAt:   time.Now(), // OK, so this is a lie, but not in payload
	}, nil
}

// not ready yet !
func DumpRepoInfo(db *gorm.DB, repo *Repo, service *Service) (bool, error) {
	return false, db.Save(repo).Error
}

// CreateOrUpdateRepo creates or updates a repo and returns true if the repo was created (vs updated)
func CreateOrUpdateRepo(db *gorm.DB, repo *Repo, service *Service) (bool, error) {
	// Get existing by remote ID and service ID
	var existing Repo
	if db.Where("remote_id = ? AND service_id = ?", repo.RemoteID, service.ID).First(&existing).RecordNotFound() {
		repo.ServiceID = service.ID
		err := db.Create(repo).Error
		return err == nil, err
	}
	repo.ID = existing.ID
	repo.ServiceID = service.ID
	repo.CreatedAt = existing.CreatedAt
	return false, db.Save(repo).Error
}

// FindRepoByID finds a repo by ID
func FindRepoByID(db *gorm.DB, ID uint) (*Repo, error) {
	var repo Repo
	if db.First(&repo, ID).RecordNotFound() {
		return nil, fmt.Errorf("Repo '%d' not found", ID)
	}
	return &repo, db.Error
}

// FindRepos finds all repos
func FindRepos(db *gorm.DB, match string) ([]Repo, error) {
	var repos []Repo
	if match != "" {
		db.Where("full_name LIKE ?",
			strings.ToLower(fmt.Sprintf("%%%s%%", match))).Order("full_name").Find(&repos)
	} else {
		db.Order("full_name").Find(&repos)
	}
	return repos, db.Error
}

// FindUntaggedRepos finds repos without any tags
func FindUntaggedRepos(db *gorm.DB, match string) ([]Repo, error) {
	var repos []Repo
	if match != "" {
		db.Raw(`
			SELECT *
			FROM REPOS S
			WHERE S.DELETED_AT IS NULL
			AND S.FULL_NAME LIKE ?
			AND S.ID NOT IN (
				SELECT REPO_ID
				FROM REPO_TAGS
			) ORDER BY S.FULL_NAME`,
			fmt.Sprintf("%%%s%%", strings.ToLower(match))).Scan(&repos)
	} else {
		db.Raw(`
			SELECT *
			FROM REPOS S
			WHERE S.DELETED_AT IS NULL
			AND S.ID NOT IN (
				SELECT REPO_ID
				FROM REPO_TAGS
			) ORDER BY S.FULL_NAME`).Scan(&repos)
	}
	return repos, db.Error
}

// FindReposByLanguageAndOrTag finds repos with the specified language and/or the specified tag
func FindReposByLanguageAndOrTag(db *gorm.DB, match string, language string, tagName string, union bool) ([]Repo, error) {
	operator := "AND"
	if union {
		operator = "OR"
	}

	var repos []Repo
	if match != "" {
		db.Raw(fmt.Sprintf(`
			SELECT * 
			FROM REPOS S, TAGS T 
			INNER JOIN REPO_TAGS ST ON S.ID = ST.REPO_ID 
			INNER JOIN TAGS ON ST.TAG_ID = T.ID 
			WHERE S.DELETED_AT IS NULL
			AND T.DELETED_AT IS NULL
			AND LOWER(S.FULL_NAME) LIKE ? 
			AND (LOWER(T.NAME) = ? 
			%s LOWER(S.LANGUAGE) = ?) 
			GROUP BY ST.REPO_ID 
			ORDER BY S.FULL_NAME`, operator),
			fmt.Sprintf("%%%s%%", strings.ToLower(match)),
			strings.ToLower(tagName),
			strings.ToLower(language)).Scan(&repos)
	} else {
		db.Raw(fmt.Sprintf(`
			SELECT * 
			FROM REPOS S, TAGS T 
			INNER JOIN REPO_TAGS ST ON S.ID = ST.REPO_ID 
			INNER JOIN TAGS ON ST.TAG_ID = T.ID 
			WHERE S.DELETED_AT IS NULL
			AND T.DELETED_AT IS NULL
			AND LOWER(T.NAME) = ? 
			%s LOWER(S.LANGUAGE) = ? 
			GROUP BY ST.REPO_ID 
			ORDER BY S.FULL_NAME`, operator),
			strings.ToLower(tagName),
			strings.ToLower(language)).Scan(&repos)
	}
	return repos, db.Error
}

// FindReposByLanguage finds repos with the specified language
func FindReposByLanguage(db *gorm.DB, match string, language string) ([]Repo, error) {
	var repos []Repo
	if match != "" {
		db.Where("full_name LIKE ? AND lower(language) = ?",
			strings.ToLower(fmt.Sprintf("%%%s%%", match)),
			strings.ToLower(language)).Order("full_name").Find(&repos)
	} else {
		db.Where("lower(language) = ?",
			strings.ToLower(language)).Order("full_name").Find(&repos)
	}
	return repos, db.Error
}

// FuzzyFindReposByName finds repos with approximate matching for full name and name
func FuzzyFindReposByName(db *gorm.DB, name string) ([]Repo, error) {
	// Try each of these, and as soon as we hit, return
	// 1. Exact match full name
	// 2. Exact match name
	// 3. Case-insensitive full name
	// 4. Case-insensitive name
	// 5. Case-insensitive like full name
	// 6. Case-insensitive like name
	var repos []Repo
	db.Where("full_name = ?", name).Order("full_name").Find(&repos)
	if len(repos) == 0 {
		db.Where("name = ?", name).Order("full_name").Find(&repos)
	}
	if len(repos) == 0 {
		db.Where("lower(full_name) = ?", strings.ToLower(name)).Order("full_name").Find(&repos)
	}
	if len(repos) == 0 {
		db.Where("lower(name) = ?", strings.ToLower(name)).Order("full_name").Find(&repos)
	}
	if len(repos) == 0 {
		db.Where("full_name LIKE ?", strings.ToLower(fmt.Sprintf("%%%s%%", name))).Order("full_name").Find(&repos)
	}
	if len(repos) == 0 {
		db.Where("name LIKE ?", strings.ToLower(fmt.Sprintf("%%%s%%", name))).Order("full_name").Find(&repos)
	}
	return repos, db.Error
}

// FindLanguages finds all languages
func FindLanguages(db *gorm.DB) ([]string, error) {
	var languages []string
	db.Table("repos").Order("language").Pluck("distinct(language)", &languages)
	return languages, db.Error
}

// AddTag adds a tag to a repo
func (repo *Repo) AddTag(db *gorm.DB, tag *Tag) error {
	repo.Tags = append(repo.Tags, *tag)
	return db.Save(repo).Error
}

// LoadTags loads the tags for a repo
func (repo *Repo) LoadTags(db *gorm.DB) error {
	// Make sure repo exists in database, or we will panic
	var existing Repo
	if db.Where("id = ?", repo.ID).First(&existing).RecordNotFound() {
		return fmt.Errorf("Repo '%d' not found", repo.ID)
	}
	return db.Model(repo).Association("Tags").Find(&repo.Tags).Error
}

// LoadTags loads the tags for a repo
func (repo *Repo) LoadTopics(db *gorm.DB) error {
	// Make sure repo exists in database, or we will panic
	var existing Repo
	if db.Where("id = ?", repo.ID).First(&existing).RecordNotFound() {
		return fmt.Errorf("Repo '%d' not found", repo.ID)
	}
	return db.Model(repo).Association("Topics").Find(&repo.Topics).Error
}

// RemoveAllTags removes all tags for a repo
func (repo *Repo) RemoveAllTags(db *gorm.DB) error {
	return db.Model(repo).Association("Tags").Clear().Error
}

// RemoveTag removes a tag from a repo
func (repo *Repo) RemoveTag(db *gorm.DB, tag *Tag) error {
	return db.Model(repo).Association("Tags").Delete(tag).Error
}

// HasTag returns whether a repo has a tag. Note that you must call LoadTags first -- no reason to incur a database call each time
func (repo *Repo) HasTag(tag *Tag) bool {
	if len(repo.Tags) > 0 {
		for _, t := range repo.Tags {
			if t.Name == tag.Name {
				return true
			}
		}
	}
	return false
}

// Index adds the repo to the index
func (repo *Repo) Index(index bleve.Index, db *gorm.DB) error {
	if err := repo.LoadTags(db); err != nil {
		return err
	}
	return index.Index(fmt.Sprintf("%d", repo.ID), repo)
}

// OpenInBrowser opens the repo in the browser
func (repo *Repo) OpenInBrowser(preferHomepage bool) error {
	var URL string
	if preferHomepage && repo.Homepage != nil && *repo.Homepage != "" {
		URL = *repo.Homepage
	} else if repo.URL != nil && *repo.URL != "" {
		URL = *repo.URL
	} else {
		if repo.Name != nil {
			return fmt.Errorf("No URL for repo '%s'", *repo.Name)
		}
		return errors.New("No URL for repo")
	}
	return open.Repot(URL)
}
