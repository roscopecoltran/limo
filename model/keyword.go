package model

import (
	"errors"
	"fmt"
	//"log"
	"strings"
    // "github.com/qor/qor"
    // "github.com/qor/admin"
	// "github.com/sirupsen/logrus"
	"github.com/jinzhu/gorm"
	// "github.com/sirupsen/logrus"
)

// Keyword represents a keyword in the database
type Keyword struct {
	gorm.Model
	Name      	string
	KeywordCount 	int    `gorm:"-"`
	StarCount 	int    `gorm:"-"`
	Stars     	[]Star `gorm:"many2many:star_keywords;"`
}

// FindKeywords finds all keywords
func FindKeywords(db *gorm.DB) ([]Keyword, error) {
	var keywords []Keyword
	db.Order("name").Find(&keywords)
	return keywords, db.Error
}

// FindKeywordsWithStarCount finds all keywords and gets their count of stars
func FindKeywordsWithStarCount(db *gorm.DB) ([]Keyword, error) {
	var keywords []Keyword

	// Create resources from GORM-backend model
	// Admin.AddResource(&Keyword{})

	rows, err := db.Raw(`
		SELECT T.NAME, COUNT(ST.KEYWORD_ID) AS STARCOUNT
		FROM KEYWORDS T
		LEFT JOIN STAR_KEYWORDS ST ON T.ID = ST.KEYWORD_ID
		WHERE T.DELETED_AT IS NULL
		GROUP BY T.ID
		ORDER BY T.NAME`).Rows()

	if err != nil {
		return keywords, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	for rows.Next() {
		var keyword Keyword
		if err = rows.Scan(&keyword.Name, &keyword.StarCount); err != nil {
			return keywords, err
		}
		keywords = append(keywords, keyword)
	}
	return keywords, db.Error
}

// FindKeywordByName finds a keyword by name
func FindKeywordByName(db *gorm.DB, name string) (*Keyword, error) {
	var keyword Keyword
	if db.Where("lower(name) = ?", strings.ToLower(name)).First(&keyword).RecordNotFound() {
		return nil, db.Error
	}
	return &keyword, db.Error
}

// FindOrCreateKeywordByName finds a keyword by name, creating if it doesn't exist
func FindOrCreateKeywordByName(db *gorm.DB, name string) (*Keyword, bool, error) {
	var keyword Keyword
	if db.Where("lower(name) = ?", strings.ToLower(name)).First(&keyword).RecordNotFound() {
		keyword.Name = name
		err := db.Create(&keyword).Error
		return &keyword, true, err
	}
	return &keyword, false, nil
}

// LoadStars loads the stars for a keyword
func (keyword *Keyword) LoadStars(db *gorm.DB, match string) error {
	// Make sure keyword exists in database, or we will panic
	var existing Keyword
	if db.Where("id = ?", keyword.ID).First(&existing).RecordNotFound() {
		return fmt.Errorf("Keyword '%d' not found", keyword.ID)
	}

	if match != "" {
		var stars []Star
		db.Raw(`
			SELECT *
			FROM STARS S
			INNER JOIN STAR_KEYWORDS ST ON S.ID = ST.STAR_ID
			WHERE S.DELETED_AT IS NULL
			AND ST.KEYWORD_ID = ?
			AND LOWER(S.FULL_NAME) LIKE ?
			ORDER BY S.FULL_NAME`,
			keyword.ID,
			fmt.Sprintf("%%%s%%", strings.ToLower(match))).Scan(&stars)
		keyword.Stars = stars
		return db.Error
	}
	return db.Model(keyword).Association("Stars").Find(&keyword.Stars).Error
}

// Rename renames a keyword -- new name must not already exist
func (keyword *Keyword) Rename(db *gorm.DB, name string) error {
	// Can't rename to the same name
	if name == keyword.Name {
		return errors.New("You can't rename to the same name")
	}

	// If they're just changing case, allow. Otherwise, block the change
	if strings.ToLower(name) != strings.ToLower(keyword.Name) {
		existing, err := FindKeywordByName(db, name)
		if err != nil {
			return err
		}
		if existing != nil {
			return fmt.Errorf("Keyword '%s' already exists", existing.Name)
		}
	}

	keyword.Name = name
	return db.Save(keyword).Error
}

// Delete deletes a keyword and disassociates it from any stars
func (keyword *Keyword) Delete(db *gorm.DB) error {
	if err := db.Model(keyword).Association("Stars").Clear().Error; err != nil {
		return err
	}
	return db.Delete(keyword).Error
}
