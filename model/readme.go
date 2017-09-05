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

// https://github.com/grimmer0125/search-github-starred/blob/develop/githubAPI.go

// Readme  represents a readme in the database
type Readme  struct {
	gorm.Model
	Name      	string
	ReadmeCount int    `gorm:"-"`
	StarCount 	int    `gorm:"-"`
	Stars     	[]Star `gorm:"many2many:star_readmes;"`
}

// FindReadmes finds all readmes
func FindReadmes(db *gorm.DB) ([]Readme , error) {
	var readmes []Readme 
	db.Order("name").Find(&readmes)
	return readmes, db.Error
}

// FindReadmesWithStarCount finds all readmes and gets their count of stars
func FindReadmesWithStarCount(db *gorm.DB) ([]Readme , error) {
	var readmes []Readme 

	// Create resources from GORM-backend model
	// Admin.AddResource(&Readme {})

	rows, err := db.Raw(`
		SELECT T.NAME, COUNT(ST.README_ID) AS STARCOUNT
		FROM READMES T
		LEFT JOIN STAR_READMES ST ON T.ID = ST.README_ID
		WHERE T.DELETED_AT IS NULL
		GROUP BY T.ID
		ORDER BY T.NAME`).Rows()

	if err != nil {
		return readmes, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	for rows.Next() {
		var readme Readme 
		if err = rows.Scan(&readme.Name, &readme.StarCount); err != nil {
			return readmes, err
		}
		readmes = append(readmes, readme)
	}
	return readmes, db.Error
}

// FindReadme ByName finds a readme by name
func FindReadmeByName(db *gorm.DB, name string) (*Readme , error) {
	var readme Readme 
	if db.Where("lower(name) = ?", strings.ToLower(name)).First(&readme).RecordNotFound() {
		return nil, db.Error
	}
	return &readme, db.Error
}

// FindOrCreateReadme ByName finds a readme by name, creating if it doesn't exist
func FindOrCreateReadmeByName(db *gorm.DB, name string) (*Readme , bool, error) {
	var readme Readme 
	if db.Where("lower(name) = ?", strings.ToLower(name)).First(&readme).RecordNotFound() {
		readme.Name = name
		err := db.Create(&readme).Error
		return &readme, true, err
	}
	return &readme, false, nil
}

// LoadStars loads the stars for a readme
func (readme *Readme ) LoadStars(db *gorm.DB, match string) error {
	// Make sure readme exists in database, or we will panic
	var existing Readme 
	if db.Where("id = ?", readme.ID).First(&existing).RecordNotFound() {
		return fmt.Errorf("Readme  '%d' not found", readme.ID)
	}

	if match != "" {
		var stars []Star
		db.Raw(`
			SELECT *
			FROM STARS S
			INNER JOIN STAR_READMES ST ON S.ID = ST.STAR_ID
			WHERE S.DELETED_AT IS NULL
			AND ST.README_ID = ?
			AND LOWER(S.FULL_NAME) LIKE ?
			ORDER BY S.FULL_NAME`,
			readme.ID,
			fmt.Sprintf("%%%s%%", strings.ToLower(match))).Scan(&stars)
		readme.Stars = stars
		return db.Error
	}
	return db.Model(readme).Association("Stars").Find(&readme.Stars).Error
}

// Rename renames a readme -- new name must not already exist
func (readme *Readme ) Rename(db *gorm.DB, name string) error {
	// Can't rename to the same name
	if name == readme.Name {
		return errors.New("You can't rename to the same name")
	}

	// If they're just changing case, allow. Otherwise, block the change
	if strings.ToLower(name) != strings.ToLower(readme.Name) {
		existing, err := FindReadmeByName(db, name)
		if err != nil {
			return err
		}
		if existing != nil {
			return fmt.Errorf("Readme  '%s' already exists", existing.Name)
		}
	}

	readme.Name = name
	return db.Save(readme).Error
}

// Delete deletes a readme and disassociates it from any stars
func (readme *Readme ) Delete(db *gorm.DB) error {
	if err := db.Model(readme).Association("Stars").Clear().Error; err != nil {
		return err
	}
	return db.Delete(readme).Error
}
