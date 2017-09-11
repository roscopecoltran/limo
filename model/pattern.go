package model

import (
	"errors"
	"fmt"
	//"log"
	"strings"
    // "github.com/qor/qor"
    // "github.com/qor/admin"
	"github.com/jinzhu/gorm"
)

// Pattern represents a pattern in the database
type Pattern struct {
	gorm.Model
	Name      	string
	PatternCount 	int    `gorm:"-"`
	StarCount 	int    `gorm:"-"`
	Stars     	[]Star `gorm:"many2many:star_patterns;"`
}

// FindPatterns finds all patterns
func FindPatterns(db *gorm.DB) ([]Pattern, error) {
	var patterns []Pattern
	db.Order("name").Find(&patterns)
	return patterns, db.Error
}

// FindPatternsWithStarCount finds all patterns and gets their count of stars
func FindPatternsWithStarCount(db *gorm.DB) ([]Pattern, error) {
	var patterns []Pattern

	// Create resources from GORM-backend model
	// Admin.AddResource(&Pattern{})

	rows, err := db.Raw(`
		SELECT T.NAME, COUNT(ST.PATTERN_ID) AS STARCOUNT
		FROM PATTERNS T
		LEFT JOIN STAR_PATTERNS ST ON T.ID = ST.PATTERN_ID
		WHERE T.DELETED_AT IS NULL
		GROUP BY T.ID
		ORDER BY T.NAME`).Rows()

	if err != nil {
		return patterns, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	for rows.Next() {
		var pattern Pattern
		if err = rows.Scan(&pattern.Name, &pattern.StarCount); err != nil {
			return patterns, err
		}
		patterns = append(patterns, pattern)
	}
	return patterns, db.Error
}

// FindPatternByName finds a pattern by name
func FindPatternByName(db *gorm.DB, name string) (*Pattern, error) {
	var pattern Pattern
	if db.Where("lower(name) = ?", strings.ToLower(name)).First(&pattern).RecordNotFound() {
		return nil, db.Error
	}
	return &pattern, db.Error
}

// FindOrCreatePatternByName finds a pattern by name, creating if it doesn't exist
func FindOrCreatePatternByName(db *gorm.DB, name string) (*Pattern, bool, error) {
	var pattern Pattern
	if db.Where("lower(name) = ?", strings.ToLower(name)).First(&pattern).RecordNotFound() {
		pattern.Name = name
		err := db.Create(&pattern).Error
		return &pattern, true, err
	}
	return &pattern, false, nil
}

// LoadStars loads the stars for a pattern
func (pattern *Pattern) LoadStars(db *gorm.DB, match string) error {
	// Make sure pattern exists in database, or we will panic
	var existing Pattern
	if db.Where("id = ?", pattern.ID).First(&existing).RecordNotFound() {
		return fmt.Errorf("Pattern '%d' not found", pattern.ID)
	}

	if match != "" {
		var stars []Star
		db.Raw(`
			SELECT *
			FROM STARS S
			INNER JOIN STAR_PATTERNS ST ON S.ID = ST.STAR_ID
			WHERE S.DELETED_AT IS NULL
			AND ST.PATTERN_ID = ?
			AND LOWER(S.FULL_NAME) LIKE ?
			ORDER BY S.FULL_NAME`,
			pattern.ID,
			fmt.Sprintf("%%%s%%", strings.ToLower(match))).Scan(&stars)
		pattern.Stars = stars
		return db.Error
	}
	return db.Model(pattern).Association("Stars").Find(&pattern.Stars).Error
}

// Rename renames a pattern -- new name must not already exist
func (pattern *Pattern) Rename(db *gorm.DB, name string) error {
	// Can't rename to the same name
	if name == pattern.Name {
		return errors.New("You can't rename to the same name")
	}

	// If they're just changing case, allow. Otherwise, block the change
	if strings.ToLower(name) != strings.ToLower(pattern.Name) {
		existing, err := FindPatternByName(db, name)
		if err != nil {
			return err
		}
		if existing != nil {
			return fmt.Errorf("Pattern '%s' already exists", existing.Name)
		}
	}

	pattern.Name = name
	return db.Save(pattern).Error
}

// Delete deletes a pattern and disassociates it from any stars
func (pattern *Pattern) Delete(db *gorm.DB) error {
	if err := db.Model(pattern).Association("Stars").Clear().Error; err != nil {
		return err
	}
	return db.Delete(pattern).Error
}
