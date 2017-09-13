package model

import (
	//"errors"
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
	entities, ok := tagg.GetTag(query)
	if !ok {
		log.WithFields(logrus.Fields{"action": "PrintCategoriesGraph", "step": "GetTag", "model": "Category", "query": query}).Warnf("Tag category not found", query)
		return nil, fmt.Errorf("Tag category not found", query)
	}
	// iterate
	catTrees := entities.PathsToAllAncestorsAsString("->")
	for _, path := range catTrees {
		log.WithFields(logrus.Fields{"action": "PrintCategoriesGraph", "step": "PathsToAllAncestorsAsString", "model": "Category", "query": query, "path": path}).Info("New path discovered.")
	}
	return entities, nil

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


