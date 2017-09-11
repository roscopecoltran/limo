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
		log.WithFields(logrus.Fields{"action": "PrintLanguageTypesGraph", "step": "GetTag", "model": "LanguageType", "query": query}).Warnf("Tag language_type not found", query)
		return nil, fmt.Errorf("Tag language_type not found", query)
	}
	// iterate
	paths := entities.PathsToAllAncestorsAsString("->")
	for _, path := range paths {
		log.WithFields(logrus.Fields{"action": "PrintLanguageTypesGraph", "step": "PathsToAllAncestorsAsString", "model": "LanguageType", "query": query, "path": path}).Info("New path discovered.")
	}
	return entities, nil

}

// FindLanguageTypes finds all language_types
func FindLanguageTypes(db *gorm.DB) ([]LanguageType, error) {
	var language_types []LanguageType
	db.Order("label").Find(&language_types)
	return language_types, db.Error
}

// FindLanguageTypesWithStarCount finds all language_types and gets their count of stars
func FindLanguageTypesWithStarCount(db *gorm.DB) ([]LanguageType, error) {
	var language_types []LanguageType
	rows, err := db.Raw(`
		SELECT LT.LABEL, COUNT(ST.LANGUAGE_TYPE_ID) AS STARCOUNT
		FROM LANGUAGE_TYPES LT
		LEFT JOIN STAR_LANGUAGE_TYPES ST ON LT.ID = ST.LANGUAGE_TYPE_ID
		WHERE LT.DELETED_AT IS NULL
		GROUP BY LT.ID
		ORDER BY LT.LABEL`).Rows()

	if err != nil {
		return language_types, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	for rows.Next() {
		var language_type LanguageType
		if err = rows.Scan(&language_type.Label); err != nil {
			return language_types, err
		}
		language_types = append(language_types, language_type)
	}
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

// LoadStars loads the stars for a language_type
func (language_type *LanguageType) LoadStars(db *gorm.DB, match string) error {
	// Make sure language_type exists in database, or we will panic
	var existing LanguageType
	if db.Where("id = ?", language_type.ID).First(&existing).RecordNotFound() {
		err := fmt.Errorf("LanguageType '%d' not found", language_type.ID)
		log.WithError(err).WithFields(logrus.Fields{"section:": "model", "typology": "tag", "step": "LoadStars"}).Errorf("LanguageType '%d' not found", language_type.ID)
		return err
	}
	if match != "" {
		var stars []Star
		db.Raw(`
			SELECT *
			FROM STARS S
			INNER JOIN STAR_LANGUAGE_TYPES SLT ON S.ID = SLT.STAR_ID
			WHERE S.DELETED_AT IS NULL
			AND SLT.LANGUAGE_TYPE_ID = ?
			AND LOWER(S.FULL_NAME) LIKE ?
			ORDER BY S.FULL_NAME`,
			language_type.ID,
			fmt.Sprintf("%%%s%%", strings.ToLower(match))).Scan(&stars)
		language_type.Stars = stars
		return db.Error
	}
	return db.Model(language_type).Association("Stars").Find(&language_type.Stars).Error
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
			log.WithError(err).WithFields(logrus.Fields{"section:": "model", "typology": "tag", "step": "Rename"}).Warnf("%#s", err)
			return err
		}
		if existing != nil {
			err := fmt.Errorf("LanguageType '%s' already exists", existing.Label)
			log.WithError(err).WithFields(logrus.Fields{"section:": "model", "typology": "tag", "step": "Rename"}).Errorf("LanguageType '%s' already exists", existing.Label)
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
