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

// LanguageType represents a language_type in the database
type LanguageType struct {
	gorm.Model
	Label      		string 		`gorm:"type:varchar(64);not null;unique"`
	Count 			int    		`gorm:"-"`
	Categories      []Category 	`gorm:"many2many:language_categories;"`
	Topics        	[]Topic 	`gorm:"many2many:language_topics;"` 				// is a dev language, human language or sign language ?!
}

type LanguageTypeResult struct {
	LanguageType  	*LanguageType
	Error 	error
}

// should provide a map[string]map[string]
func TestLanguageTypesGraph(query string) (taggraph.Tagger, error) {
	tagg.AddChildToTag("shirts", "clothes")
	tagg.AddChildToTag("pants", "clothes")
	tagg.AddChildToTag("dress clothes", "clothes")
	tagg.AddChildToTag("shirts", "dress clothes")
	tagg.AddChildToTag("shirts", "tops")
	tagg.AddChildToTag("tops", "casual")
	tagg.AddChildToTag("casual", "clothes")
	entities, ok := tagg.GetTag(query)
	if !ok {
		log.WithFields(
			logrus.Fields{	"action": 	"PrintLanguageTypesGraph", 
							"step": 	"GetTag", 
							"model": 	"LanguageType", 
							"query": 	query}).Warnf("Tag language_type not found", query)
		return nil, fmt.Errorf("Tag language_type not found", query)
	}
	// iterate
	paths := entities.PathsToAllAncestorsAsString("->")
	for _, path := range paths {
		log.WithFields(
			logrus.Fields{	"action": 	"PrintLanguageTypesGraph", 
							"step": 	"PathsToAllAncestorsAsString", 
							"model": 	"LanguageType", 
							"query": 	query, 
							"path": 	path}).Info("New path discovered.")
	}
	return entities, nil

}

// FindLanguageTypes finds all language_types
func FindLanguageTypes(db *gorm.DB) ([]LanguageType, error) {
	var language_types []LanguageType
	db.Order("label").Find(&language_types)
	return language_types, db.Error
}

// FindLanguageTypeByLabel finds a language_type by label
func FindLanguageTypeByLabel(db *gorm.DB, label string) (*LanguageType, error) {
	var language_type LanguageType
	if db.Where("lower(label) = ?", strings.ToLower(label)).First(&language_type).RecordNotFound() {
		return nil, db.Error
	}
	return &language_type, db.Error
}

// FindOrCreateLanguageTypeByLabel finds a language_type by label, creating if it doesn't exist
func FindOrCreateLanguageTypeByLabel(db *gorm.DB, label string) (*LanguageType, bool, error) {
	var language_type LanguageType
	if db.Where("lower(label) = ?", strings.ToLower(label)).First(&language_type).RecordNotFound() {
		language_type.Label = label
		err := db.Create(&language_type).Error
		return &language_type, true, err
	}
	return &language_type, false, nil
}

// Rename renames a language_type -- new label must not already exist
func (language_type *LanguageType) Rename(db *gorm.DB, label string) error {
	// Can't rename to the same label
	if label == language_type.Label {
		return errors.New("You can't rename to the same label")
	}

	// If they're just changing case, allow. Otherwise, block the change
	if strings.ToLower(label) != strings.ToLower(language_type.Label) {
		existing, err := FindLanguageTypeByLabel(db, label)
		if err != nil {
			log.WithError(err).WithFields(
				logrus.Fields{	"section:": "model", 
								"typology": "tag", 
								"step": 	"Rename"}).Warnf("%#s", err)
			return err
		}
		if existing != nil {
			err := fmt.Errorf("LanguageType '%s' already exists", existing.Label)
			log.WithError(err).WithFields(
				logrus.Fields{	"section:": "model", 
								"typology": "tag", 
								"step": 	"Rename"}).Errorf("LanguageType '%s' already exists", existing.Label)
			return err
		}
	}

	language_type.Label = label
	return db.Save(language_type).Error
}

// Delete deletes a language_type and disassociates it from any stars
func (language_type *LanguageType) Delete(db *gorm.DB) error {
	if err := db.Model(language_type).Association("Stars").Clear().Error; err != nil {
		return err
	}
	return db.Delete(language_type).Error
}
