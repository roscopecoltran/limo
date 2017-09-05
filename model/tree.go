package model

import (
	"errors"
	"fmt"
	"log"
	"strings"
    // "github.com/qor/qor"
    // "github.com/qor/admin"
	"github.com/jinzhu/gorm"
)

// https://github.com/Termina1/starlight/blob/master/handlers/filenames_extractor.go
// 

// Tree represents a tree in the database
type Tree struct {
	gorm.Model
	Name      		string
	TreeCount 		int    `gorm:"-"`
	StarCount 		int    `gorm:"-"`
	Stars     		[]Star `gorm:"many2many:star_trees;"`
}

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

// FindTreesWithStarCount finds all trees and gets their count of stars
func FindTreesWithStarCount(db *gorm.DB) ([]Tree, error) {
	var trees []Tree

	// Create resources from GORM-backend model
	// Admin.AddResource(&Tree{})

	rows, err := db.Raw(`
		SELECT T.NAME, COUNT(ST.TREE_ID) AS STARCOUNT
		FROM TREES T
		LEFT JOIN STAR_TREES ST ON T.ID = ST.TREE_ID
		WHERE T.DELETED_AT IS NULL
		GROUP BY T.ID
		ORDER BY T.NAME`).Rows()

	if err != nil {
		return trees, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	for rows.Next() {
		var tree Tree
		if err = rows.Scan(&tree.Name, &tree.StarCount); err != nil {
			return trees, err
		}
		trees = append(trees, tree)
	}
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

// LoadStars loads the stars for a tree
func (tree *Tree) LoadStars(db *gorm.DB, match string) error {
	// Make sure tree exists in database, or we will panic
	var existing Tree
	if db.Where("id = ?", tree.ID).First(&existing).RecordNotFound() {
		return fmt.Errorf("Tree '%d' not found", tree.ID)
	}

	if match != "" {
		var stars []Star
		db.Raw(`
			SELECT *
			FROM TREES S
			INNER JOIN STAR_TREES ST ON S.ID = ST.STAR_ID
			WHERE S.DELETED_AT IS NULL
			AND ST.TREE_ID = ?
			AND LOWER(S.FULL_NAME) LIKE ?
			ORDER BY S.FULL_NAME`,
			tree.ID,
			fmt.Sprintf("%%%s%%", strings.ToLower(match))).Scan(&stars)
		tree.Stars = stars
		return db.Error
	}
	return db.Model(tree).Association("Stars").Find(&tree.Stars).Error
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
