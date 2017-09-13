package model

import (
	"errors"
	"fmt"
	//"log"
	"strings"
	"github.com/jinzhu/gorm"
	// "github.com/sirupsen/logrus"

	// el "github.com/src-d/enry"
	// rl "github.com/rai-project/linguist"
	// gl "github.com/generaltso/linguist"

)

// https://github.com/Termina1/starlight/blob/master/handlers/filenames_extractor.go
// https://github.com/xanzy/go-gitlab/blob/master/examples/repository_files.go

// Tree represents a tree in the database
type Tree struct {
	gorm.Model
	Name      		string
	TreeCount 		int    `gorm:"-"`
	Stars     		[]Star `gorm:"many2many:star_trees;"`
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
func FindTreeByName(db *gorm.DB, name string) (*Tree, error) {
	var tree Tree
	if db.Where("lower(name) = ?", strings.ToLower(name)).First(&tree).RecordNotFound() {
		return nil, db.Error
	}
	return &tree, db.Error
}

// FindOrCreateTreeByName finds a tree by name, creating if it doesn't exist
func FindOrCreateTreeByName(db *gorm.DB, name string) (*Tree, bool, error) {
	var tree Tree
	if db.Where("lower(name) = ?", strings.ToLower(name)).First(&tree).RecordNotFound() {
		tree.Name = name
		err := db.Create(&tree).Error
		return &tree, true, err
	}
	return &tree, false, nil
}

// Rename renames a tree -- new name must not already exist
func (tree *Tree) Rename(db *gorm.DB, name string) error {
	// Can't rename to the same name
	if name == tree.Name {
		return errors.New("You can't rename to the same name")
	}

	// If they're just changing case, allow. Otherwise, block the change
	if strings.ToLower(name) != strings.ToLower(tree.Name) {
		existing, err := FindTreeByName(db, name)
		if err != nil {
			return err
		}
		if existing != nil {
			return fmt.Errorf("Tree '%s' already exists", existing.Name)
		}
	}

	tree.Name = name
	return db.Save(tree).Error
}

// Delete deletes a tree and disassociates it from any stars
func (tree *Tree) Delete(db *gorm.DB) error {
	if err := db.Model(tree).Association("Stars").Clear().Error; err != nil {
		return err
	}
	return db.Delete(tree).Error
}
