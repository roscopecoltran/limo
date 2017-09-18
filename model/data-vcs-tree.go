package model

import (
	"github.com/jinzhu/gorm"
	"strings"
	//"github.com/qor/sorting"
	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
	// "github.com/xanzy/go-gitlab"
	// el "github.com/src-d/enry"
	// rl "github.com/rai-project/linguist"
	// gl "github.com/generaltso/linguist"
)

// https://github.com/Termina1/starlight/blob/master/handlers/filenames_extractor.go
// https://github.com/xanzy/go-gitlab/blob/master/examples/repository_files.go

type Tree struct {
	gorm.Model `json:"-" yaml:"-"`
	//sorting.SortingDESC
	//ServiceID   				uint 				`json:"service_id" yaml:"service_id"`
	RemoteURI string      `json:"remote_uri,omitempty" yaml:"remote_uri,omitempty"`
	SHA       *string     `json:"sha,omitempty" yaml:"sha,omitempty"`
	Count     int         `gorm:"-" json:"count,omitempty" yaml:"count,omitempty"`
	Entries   []TreeEntry `gorm:"many2many:star_tree_entries;" json:"tree,omitempty" yaml:"tree,omitempty"`
	Stars     []Star      `gorm:"many2many:star_trees;"`
}

// Tree represents a tree in the database
type TreeEntry struct {
	gorm.Model
	//sorting.SortingDESC
	Mode      *string `json:"mode,omitempty" yaml:"mode,omitempty"`
	SHA       *string `json:"sha,omitempty" yaml:"sha,omitempty"`
	Type      *string `json:"type,omitempty" yaml:"type,omitempty"`
	MimeType  string  `json:"mime_type,omitempty" yaml:"mime_type,omitempty"`
	Size      *int    `json:"size,omitempty" yaml:"size,omitempty"`
	Content   *string `json:"content,omitempty" yaml:"content,omitempty"`
	RemoteURL *string `json:"remote_url,omitempty" yaml:"remote_url,omitempty"`
	FilePath  *string `json:"file_path,omitempty" yaml:"file_path,omitempty"`
	Priority  int     `default:"1" json:"priority" yaml:"priority"` // If a changeset affects multiple segments, priority can describe the order of segments listed
	//Ast 						[]Ast 				`json:"content_ast,omitempty" yaml:"content_ast,omitempty"`
	//Topics 					[]Topic 			`gorm:"many2many:star_vcs_tree_file_patterns;" json:"topics" yaml:"topics"` 										// List of topics
	//Languages 					[]Language 			`gorm:"many2many:star_vcs_tree_languages;" json:"languages" yaml:"languages"` 										// List of topics
	//Stars     					[]Star 				`gorm:"many2many:star_vcs_tree_entries;"`
	// CheckFilePatterns 		[]Pattern 			`gorm:"many2many:star_vcs_tree_check_file_patterns;" json:"check_file_patterns" yaml:"check_file_patterns"` 							// List of regexps to specify which file to include in this Segment
	// CheckContentPatterns 	[]Pattern 			`gorm:"many2many:star_vcs_tree_check_content_patterns;" json:"check_content_patterns" yaml:"check_content_patterns"` 					// List of regexps to specify what patch content should be included in this Segment
	// IgnoreFilePatterns 		[]Pattern 			`gorm:"many2many:star_vcs_tree_ignore_file_patterns;" json:"ignore_file_patterns" yaml:"ignore_file_patterns"` 			// List of regexps to exclude files matched by FilePatterns regex
	// IgnoreContentPatterns 	[]Pattern 			`gorm:"many2many:star_vcs_tree_ignore_file_patterns;" json:"ignore_content_patterns" yaml:"ignore_content_patterns"` 	// List of regexps to exclude patch content matched by `ContentPatterns`
}

/*
	// Examples

	lang, safe := enry.GetLanguageByExtension("foo.go")
	fmt.Println(lang)
	// result: Go

	lang, safe := enry.GetLanguageByContent("foo.m", "<matlab-code>")
	fmt.Println(lang)
	// result: Matlab

	lang, safe := enry.GetLanguageByContent("bar.m", "<objective-c-code>")
	fmt.Println(lang)
	// result: Objective-C

	// all strategies together
	lang := enry.GetLanguage("foo.cpp", "<cpp-code>")

	// get a list of possible languages for a given file
	langs := enry.GetLanguages("foo.h",  "<cpp-code>")
	// result: []string{"C++", "C"}

	langs := enry.GetLanguagesByExtension("foo.asc", "<content>", nil)
	// result: []string{"AGS Script", "AsciiDoc", "Public Key"}

	langs := enry.GetLanguagesByFilename("Gemfile", "<content>", []string{})
	// result: []string{"Ruby"}

*/

func ProcessingFilesContent(githubTrees []github.Tree, processor string, forceLocalRepository bool) []Language {
	var languages []Language // LanguageDetection, LanguageDetected
	if len(githubTrees) > 0 {
		for _, t := range githubTrees {
			// var treeEntries []TreeEntry
			for _, e := range t.Entries {
				// processors: enry, rai-linguist, gso-linguist
				// do something...
				// refs:
				//  - gso-linguist: 	https://github.com/generaltso/linguist/tree/master/cmd/l/main.go
				//  - gso-linguist: 	https://github.com/generaltso/linguist/blob/master/cmd/tokenizer_test/main.go
				//  - rai-linguist: 	https://github.com/rai-project/linguist/blob/master/linguist_test.go
				//  - enry: 			https://github.com/src-d/enry/blob/master/common_test.go
				log.WithFields(logrus.Fields{
					"method.name": "DetectFilesLanguages(...)",
					"var.e.Type":  e.Type,
					"var.e.Path":  e.Path,
				}).Info("processing new detection...")
			}
		}
	}
	return languages
}

/*
func MapTrees(githubTrees []github.Tree, remoteURI string) ([]Tree) {
	var trees []Tree
	if len(githubTrees) > 0 {
		for _, t := range githubTrees {
			var treeEntries []TreeEntry
			for _, e := range tree.Entries {
				treeEntry := &TreeEntry{
					Mode: 		e.Mode,
					SHA: 		e.SHA,
					Type: 		e.Type,
					Size: 		e.Size,
					FilePath: 	e.Path,
					RemoteURL: 	e.URL,
					Content: 	e.Content,
				}
				treeEntries = append(treeEntries, treeEntry)
			}
			treeList := &Tree{
				RemoteURI: 	remoteURI,
				SHA: 		t.SHA,
				Entries: 	treeEntries,
			}
			trees = append(trees, treeList)
		}
	}
	return trees
}*/

// FindOrCreateTreeByName finds a tree by name, creating if it doesn't exist
func CheckTreeEntry(db *gorm.DB, uri string, sha *string) (*Tree, bool, error) {
	var tree Tree
	if db.Where("lower(uri) = ? AND lower(sha)", strings.ToLower(uri), strings.ToLower(*sha)).First(&tree).RecordNotFound() {
		tree.RemoteURI = uri
		tree.SHA = sha
		err := db.Create(&tree).Error
		if err != nil {
			log.WithError(err).WithFields(
				logrus.Fields{
					"method.name": "CheckTreeEntry(...)",
				}).Error("processing new detection...")
		}
		return &tree, true, err
	}
	return &tree, false, nil
}

/*
func (dbs *DatabaseDrivers) importTrees(githubTrees []github.Tree) ([]*Tree, error) {
	var trees []Tree
	if len(githubTrees) > 0 {
		for _, tree := range githubTrees {
			var treeEntries []TreeEntry
			for _, entry := range tree.Entries {
				treeEntries = append(treeEntries, &TreeEntry{
					Mode: 		entry.Mode,
					SHA: 		entry.SHA,
					Type: 		entry.Type,
					Size: 		entry.Size,
					Content: 	entry.Content,
					RemoteURL: 	entry.URL,
					Content: 	entry.Content,
				})
			}
			treeList := &Tree{
				SHA: 		*tree.SHA,
				Entries: 	treeEntries,
			}
			trees = append(trees, treeList)
		}
	}
	return trees
}
*/

/*

// https://github.com/asciimoo/chiefr/blob/master/chiefr.go

// Describe a project segment and its members and resources
// ProjectSegment can be any logical piece of a project
type ProjectSegment struct {
	// Name of the segment
	Name string `ini:"-"`
	// Repository to submit patches
	Repository string
	// URL of the chat service
	Chat string
	// URL of the mailing list
	MailList string
	// URL of the issue tracker
	IssueTracker string
	// Comma separated list of project members who are responsible for this Segment
	Chiefs []string
	// Comma separated list of project members who are responsible only for code reviews in this Segment
	Reviewers []string
	// List of regexps to specify which file to include in this Segment
	FilePatterns []string
	// List of regexps to specify what patch content should be included in this Segment
	ContentPatterns []string
	// List of regexps to exclude files matched by FilePatterns regex
	FileExcludePatterns []string
	// List of regexps to exclude patch content matched by `ContentPatterns`
	ContentExcludePatterns []string
	// If a changeset affects multiple segments, priority can describe the order of segments listed
	Priority int
	// Comma separated list of segment's topics
	Topics []string
}

type ProjectSegments map[string]*ProjectSegment

type Config struct {
	Segments ProjectSegments
}

*/

// FindTrees finds all trees
func FindTrees(db *gorm.DB) ([]Tree, error) {
	var trees []Tree
	db.Order("name").Find(&trees)
	return trees, db.Error
}

// FindTreeByName finds a tree by name
func FindTreeByURI(db *gorm.DB, RemoteURI string) (*Tree, error) {
	var tree Tree
	if db.Where("lower(RemoteURI) = ?", strings.ToLower(RemoteURI)).First(&tree).RecordNotFound() {
		return nil, db.Error
	}
	return &tree, db.Error
}

// Delete deletes a tree and disassociates it from any stars
func (tree *Tree) Delete(db *gorm.DB) error {
	if err := db.Model(tree).Association("Stars").Clear().Error; err != nil {
		return err
	}
	return db.Delete(tree).Error
}
