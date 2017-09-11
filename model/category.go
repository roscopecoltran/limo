package model

import (
	"errors"
	"fmt"
	"strings"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/ckaznocha/taggraph"
)

// https://github.com/cloudfoundry-incubator/cf-extensions/blob/master/bot/repos.go

// Category represents a category in the database
type Category struct {
	gorm.Model
	Name      		string `gorm:"type:varchar(64);not null;unique"`
	Count 			int    `gorm:"-"`

}

type CategoryResult struct {
	Category  	*Category
	Error 	error
}

// should provide a map[string]map[string]
func TestCategoriesGraph(query string) (taggraph.Tagger, error) {
	tagg.AddChildToTag("shirts", "clothes")
	tagg.AddChildToTag("pants", "clothes")
	tagg.AddChildToTag("dress clothes", "clothes")
	tagg.AddChildToTag("shirts", "dress clothes")
	tagg.AddChildToTag("shirts", "tops")
	tagg.AddChildToTag("tops", "casual")
	tagg.AddChildToTag("casual", "clothes")
	categories, ok := tagg.GetTag(query)
	if !ok {
		log.WithFields(logrus.Fields{"action": "PrintCategoriesGraph", "step": "GetTag", "model": "Category", "query": query}).Warnf("Tag category not found", query)
		return nil, fmt.Errorf("Tag category not found", query)
	}
	// iterate
	catTrees := categories.PathsToAllAncestorsAsString("->")
	for _, path := range catTrees {
		log.WithFields(logrus.Fields{"action": "PrintCategoriesGraph", "step": "PathsToAllAncestorsAsString", "model": "Category", "query": query, "path": path}).Info("New path discovered.")
	}
	return catTrees, nil

}

// FindCategories finds all categories
func FindCategories(db *gorm.DB) ([]Category, error) {
	var categories []Category
	db.Order("name").Find(&categories)
	return categories, db.Error
}

// FindCategoryByName finds a category by name
func FindCategoryByName(db *gorm.DB, name string) (*Category, error) {
	var category Category
	if db.Where("lower(name) = ?", strings.ToLower(name)).First(&category).RecordNotFound() {
		return nil, db.Error
	}
	return &category, db.Error
}

// FindOrCreateCategoryByName finds a category by name, creating if it doesn't exist
func FindOrCreateCategoryByName(db *gorm.DB, name string) (*Category, bool, error) {
	var category Category
	if db.Where("lower(name) = ?", strings.ToLower(name)).First(&category).RecordNotFound() {
		category.Name = name
		err := db.Create(&category).Error
		return &category, true, err
	}
	return &category, false, nil
}

// LoadStars loads the stars for a category
func (category *Category) LoadCategories(db *gorm.DB, match string) error {
	// Make sure category exists in database, or we will panic
	var existing Category
	if db.Where("id = ?", category.ID).First(&existing).RecordNotFound() {
		err := fmt.Errorf("Category '%d' not found", category.ID)
		log.WithError(err).WithFields(logrus.Fields{"section:": "model", "typology": "tag", "step": "LoadStars"}).Errorf("Category '%d' not found", category.ID)
		return err
	}
	if match != "" {
		var stars []Star
		db.Raw(`
			SELECT *
			FROM STARS S
			INNER JOIN STAR_CATEGORIES ST ON S.ID = ST.STAR_ID
			WHERE S.DELETED_AT IS NULL
			AND ST.CATEGORY_ID = ?
			AND LOWER(S.FULL_NAME) LIKE ?
			ORDER BY S.FULL_NAME`,
			category.ID,
			fmt.Sprintf("%%%s%%", strings.ToLower(match))).Scan(&stars)
		category.Stars = stars
		return db.Error
	}
	return db.Model(category).Association("Stars").Find(&category.Stars).Error
}

// Rename renames a category -- new name must not already exist
func (category *Category) Rename(db *gorm.DB, name string) error {
	// Can't rename to the same name
	if name == category.Name {
		return errors.New("You can't rename to the same name")
	}

	// If they're just changing case, allow. Otherwise, block the change
	if strings.ToLower(name) != strings.ToLower(category.Name) {
		existing, err := FindCategoryByName(db, name)
		if err != nil {
			log.WithError(err).WithFields(logrus.Fields{"section:": "model", "typology": "tag", "step": "Rename"}).Warnf("%#s", err)
			return err
		}
		if existing != nil {
			err := fmt.Errorf("Category '%s' already exists", existing.Name)
			log.WithError(err).WithFields(logrus.Fields{"section:": "model", "typology": "tag", "step": "Rename"}).Errorf("Category '%s' already exists", existing.Name)
			return err
		}
	}

	category.Name = name
	return db.Save(category).Error
}

// Delete deletes a category and disassociates it from any stars
func (category *Category) Delete(db *gorm.DB) error {
	if err := db.Model(category).Association("Stars").Clear().Error; err != nil {
		return err
	}
	return db.Delete(category).Error
}
