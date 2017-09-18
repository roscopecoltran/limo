package model

import (
	"errors"
	"fmt"
	// //"log"
	"path"
	// "os"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	// "github.com/qor/qor"
	// "github.com/qor/admin"
	// "github.com/xanzy/go-gitlab"
	"github.com/google/go-github/github"
	"github.com/jinzhu/gorm"
	//"github.com/qor/sorting"
	// "github.com/davecgh/go-spew/spew"
	jsoniter "github.com/json-iterator/go"
)

// https://github.com/grimmer0125/search-github-starred/blob/develop/githubAPI.go
// https://github.com/google/go-github/blob/master/github/examples_test.go
// https://github.com/google/go-github/blob/master/github/search_test.go
// https://github.com/google/go-github/blob/master/github/activity_star_test.go
// https://github.com/google/go-github/blob/master/github/search.go

// http://jinzhu.me/gorm/models.html#model-definition
// Readme  represents a readme in the database
type Readme struct {
	gorm.Model `json:"-" yaml:"-"`
	//sorting.SortingDESC
	UserID      string  `gorm:"column:user_id" json:"open_id" yaml:"user_id"`
	RemoteID    string  `gorm:"column:remote_id" json:"remote_id" yaml:"open_id"`
	RemoteURI   string  `gorm:"column:remote_uri" json:"remote_uri" yaml:"remote_uri"`
	Name        *string `gorm:"column:name" json:"name" yaml:"name"`
	Path        *string `gorm:"column:path" json:"path" yaml:"path"`
	Content     *string `gorm:"column:content" json:"content" yaml:"content"`
	Decoded     string  `gorm:"column:decoded" json:"decoded" yaml:"decoded"`
	SHA         *string `gorm:"column:sha" json:"sha" yaml:"sha"` // boltholdIndex:"SHA"
	URL         *string `gorm:"column:url" json:"url" yaml:"url"`
	DownloadURL *string `gorm:"column:download_url" json:"download_url" yaml:"download_url"`
	Language    string  `gorm:"column:language" json:"language" yaml:"language"`
	Type        *string `gorm:"column:type" json:"type" yaml:"type"`
	Encoding    *string `gorm:"column:encoding" json:"encoding" yaml:"encoding"`
	Size        *int    `gorm:"column:size" json:"size" yaml:"size"`
}

//type Token struct {
//	name string
//}

type ReadmeResult struct {
	Readme *Readme
	Error  error
}

// iAreaId, ok := val.(int) // Alt. non panicking version
// i, errInt := strconv.ParseInt(v.(string), 10, 64)

// TokenizeContent("I'm from Iceland and I make goat cheese. How about you? Do you work?")
// => ["iceland", "goat", "cheese"]

// ref: https://github.com/rawlingsj/gostats/blob/f8a768818f4689071d181543c92cd1beb6f734b4/vendor/github.com/google/go-github/github/examples_test.go#L29-L45
func NewReadmeFromGithub(readme github.RepositoryContent, remoteId int, userId int, remoteUri string) (*Readme, error) {
	// TestTopicsGraph()
	decoded, err := readme.GetContent()
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"action": "NewReadmeFromGithub", "step": "GetContent", "userId": userId, "remoteId": remoteId, "remoteUri": remoteUri}).Warn("extracting error on readme informations with readme.GetContent")
		return nil, err
	}
	readmeMetaInfo := &Readme{
		RemoteID:    strconv.Itoa(remoteId),
		UserID:      strconv.Itoa(userId),
		RemoteURI:   remoteUri,
		Name:        readme.Name,
		Path:        readme.Path,
		Content:     readme.Content,
		Decoded:     decoded,
		SHA:         readme.SHA,
		URL:         readme.URL,
		DownloadURL: readme.DownloadURL,
		Type:        readme.Type,
		Encoding:    readme.Encoding,
		Size:        readme.Size,
	}
	if decoded != "" {
		dumpReadme, err := jsoniter.Marshal(readmeMetaInfo)
		if err != nil {
			log.WithError(err).WithFields(logrus.Fields{"action": "NewReadmeFromGithub", "step": "JsoniterMarshalReadme", "userId": userId, "remoteId": remoteId, "remoteUri": remoteUri}).Warn("dump error on readme informations received with jsoniter")
		} else {
			// add a common prefix cache
			dumpPrefixPath := path.Join("cache", "vcs", remoteUri)
			if err := NewDump([]byte(fmt.Sprintf("[%s]\n", dumpReadme)), dumpPrefixPath, "readme", []string{"json", "yaml"}); err != nil {
				log.WithError(err).WithFields(logrus.Fields{"action": "NewReadmeFromGithub", "step": "NewDumpReadme", "readme": readmeMetaInfo}).Warn("could not dump Readme to specified format")
				return nil, err // errors.New("Could not get the readme file content.")
			}
		}
	}
	return readmeMetaInfo, nil
}

// FindReadmes finds all readmes
func FindReadmes(db *gorm.DB) ([]Readme, error) {
	var readmes []Readme
	db.Order("name").Find(&readmes)
	return readmes, db.Error
}

// FindReadme ByName finds a readme by name
func FindReadmeByName(db *gorm.DB, name string) (*Readme, error) {
	var readme Readme
	if db.Where("lower(name) = ?", strings.ToLower(name)).First(&readme).RecordNotFound() {
		return nil, db.Error
	}
	return &readme, db.Error
}

// FindOrCreateReadme ByName finds a readme by name, creating if it doesn't exist
func FindOrCreateReadmeByName(db *gorm.DB, name string) (*Readme, bool, error) {
	var readme Readme
	if db.Where("lower(name) = ?", strings.ToLower(name)).First(&readme).RecordNotFound() {
		*readme.Name = name
		err := db.Create(&readme).Error
		return &readme, true, err
	}
	return &readme, false, nil
}

// Rename renames a readme -- new name must not already exist
func (readme *Readme) Rename(db *gorm.DB, name string) error {
	// Can't rename to the same name
	if name == *readme.Name {
		// err := errors.New("You can't rename to the same name")
		// log.WithError(err).WithFields(logrus.Fields{"action": "Rename", "model": "Readme"}).Warn("You can't rename to the same name")
		return errors.New("You can't rename to the same name")
	}
	// If they're just changing case, allow. Otherwise, block the change
	if strings.ToLower(name) != strings.ToLower(*readme.Name) {
		existing, err := FindReadmeByName(db, name)
		if err != nil {
			return err
		}
		if existing != nil {
			return fmt.Errorf("Readme  '%s' already exists", existing.Name)
		}
	}
	*readme.Name = name
	return db.Save(readme).Error
}

// Delete deletes a readme and disassociates it from any stars
func (readme *Readme) Delete(db *gorm.DB) error {
	if err := db.Model(readme).Association("Readmes").Clear().Error; err != nil {
		return err
	}
	return db.Delete(readme).Error
}
