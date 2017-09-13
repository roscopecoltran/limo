package model

import (
	// golang
	// "errors"
	"io/ioutil"
	"os"
	"strings"
	"time"
	// "crypto/md5"
	// "fmt"
	// database
	"github.com/jinzhu/gorm"
	// wiki - markdowns
	"github.com/mschoch/blackfriday-text"
	"github.com/russross/blackfriday"
	// vcs
	// "github.com/google/go-github/github"	
	// notification fs
	// "github.com/fsnotify/fsnotify"
	// search
	// "github.com/blevesearch/bleve"
	// logs
	"github.com/sirupsen/logrus"
)

type WikiPage struct {
	gorm.Model
	Name               string    `yaml:"name" json:"name"`
	Body               string    `yaml:"body" json:"body"`
	ModifiedBy         string    `yaml:"modified_by" json:"modified_by"`
	ModifiedByName     string    `yaml:"modified_by_name" json:"modified_by_name"`
	ModifiedByEmail    string    `yaml:"modified_by_email" json:"modified_by_email"`
	ModifiedByGravatar string    `yaml:"modified_by_gravatar" json:"modified_by_gravatar"`
	Modified           time.Time `yaml:"modified" json:"modified"`
}

type WikiResult struct {
	WikiPage  	*WikiPage
	Error 		error
}

func (w *WikiPage) Type() string {
	return "wiki"
}

func NewWikiFromFile(path string) (*WikiPage, error) {
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"section:": "model", "typology": "wiki", "step": "NewWikiFromFile"}).Warnf("%#s", err)
		return nil, err
	}
	cleanedUpBytes := cleanupMarkdown(fileBytes)
	name := path
	lastSlash := strings.LastIndex(path, string(os.PathSeparator))
	if lastSlash > 0 {
		name = name[lastSlash+1:]
	}
	if strings.HasSuffix(name, ".md") {
		name = name[0 : len(name)-len(".md")]
	}
	rv := WikiPage{
		Name: name,
		Body: string(cleanedUpBytes),
	}
	return &rv, nil
}

func cleanupMarkdown(input []byte) []byte {
	extensions := 0
	renderer := blackfridaytext.TextRenderer()
	output := blackfriday.Markdown(input, renderer, extensions)
	return output
}

/*
func OpenGitRepo(path string) *github.Repository {
	repo, err := github.OpenRepository(path)
	if err != nil {
		log.Fatal(err)
	}
	return repo
}
*/

/*
func DoGitStuff(repo *github.Repository, path string, wiki *WikiPage) {

	// lookup head
	head, err := repo.Head()
	if err != nil {
		log.Print(err)
	} else {
		// lookup commit object
		headOid := head.Target()
		commit, err := repo.LookupCommit(headOid)
		if err != nil {
			log.Print(err)
		}

		// start diffing backwards
		diffCommit, err := recursiveDiffLookingForFile(repo, commit, path)
		if err != nil {
			log.Print(err)
		} else if diffCommit != nil {
			author := diffCommit.Author()
			wiki.ModifiedByName = author.Name
			wiki.ModifiedByEmail = author.Email
			wiki.Modified = author.When
			if wiki.ModifiedByEmail != "" {
				wiki.ModifiedByGravatar = gravatarHashFromEmail(wiki.ModifiedByEmail)
				log.Printf("gravatar hash is: %s", wiki.ModifiedByGravatar)
			}
		} else {
			log.Printf("unable to find commit where file changed")
		}
	}
}
*/


/*
func recursiveDiffLookingForFile(repo *github.Repository, commit *github.Commit, path string) (*github.Commit, error) {
	log.Printf("checking commit %s", commit.Id())
	// if there is a parent, diff against it
	// totally not going to think about branches
	if commit.ParentCount() > 0 {
		parent := commit.Parent(0)

		found := false
		dcb := func(dd github.DiffDelta, x float64) (github.DiffForEachHunkCallback, error) {
			if dd.NewFile.Path == path {
				found = true
			} else if dd.OldFile.Path == path {
				found = true
			}
			return nil, nil
		}

		parentTree, err := parent.Tree()
		if err != nil {
			return nil, err
		}
		commitTree, err := commit.Tree()
		if err != nil {
			return nil, err
		}
		diffOptions, err := github.DefaultDiffOptions()
		if err != nil {
			return nil, err
		}
		diff, err := repo.DiffTreeToTree(parentTree, commitTree, &diffOptions)
		if err != nil {
			return nil, err
		} else {
			diff.ForEach(dcb, github.DiffDetailFiles)
			if found {
				return commit, nil
			} else {
				return recursiveDiffLookingForFile(repo, parent, path)
			}
		}
	} else {
		// if there is no parent check to see if this file
		// was in the commit, if so, this is its
		commitTree, err := commit.Tree()
		if err != nil {
			return nil, err
		}
		treeEntry := commitTree.EntryByName(path)
		if treeEntry != nil {
			return commit, nil
		}
		return nil, nil
	}
}
*/

