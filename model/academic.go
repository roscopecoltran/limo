package model

import (
	"errors"
	"fmt"
	//"log"
	"strings"
    // "github.com/qor/qor"
    // "github.com/qor/admin"
	"github.com/jinzhu/gorm"
	// "github.com/sirupsen/logrus"
)

// Academic represents a readmes in the database
type Academic struct {
	gorm.Model
	Name      		string
	Description 	*string
	Homepage    	*string
	AcademicCount 	int    `gorm:"-"`
	StarCount 		int    `gorm:"-"`
	Stars     		[]Star `gorm:"many2many:star_academics;"`
}

// FindAcademics finds all academics
func FindAcademics(db *gorm.DB) ([]Academic, error) {
	var academics []Academic
	db.Order("name").Find(&academics)
	return academics, db.Error
}

// FindAcademicsWithStarCount finds all academics and gets their count of stars
func FindAcademicsWithStarCount(db *gorm.DB) ([]Academic, error) {
	var academics []Academic

	// Create resources from GORM-backend model
	// Admin.AddResource(&Academic{})

	rows, err := db.Raw(`
		SELECT T.NAME, COUNT(ST.ACADEMIC_ID) AS STARCOUNT
		FROM ACADEMICS T
		LEFT JOIN STAR_ACADEMICS ST ON T.ID = ST.ACADEMIC_ID
		WHERE T.DELETED_AT IS NULL
		GROUP BY T.ID
		ORDER BY T.NAME`).Rows()

	if err != nil {
		return academics, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	for rows.Next() {
		var readmes Academic
		if err = rows.Scan(&readmes.Name, &readmes.StarCount); err != nil {
			return academics, err
		}
		academics = append(academics, readmes)
	}
	return academics, db.Error
}

// FindAcademicByName finds a readmes by name
func FindAcademicByName(db *gorm.DB, name string) (*Academic, error) {
	var readmes Academic
	if db.Where("lower(name) = ?", strings.ToLower(name)).First(&readmes).RecordNotFound() {
		return nil, db.Error
	}
	return &readmes, db.Error
}

// FindOrCreateAcademicByName finds a readmes by name, creating if it doesn't exist
func FindOrCreateAcademicByName(db *gorm.DB, name string) (*Academic, bool, error) {
	var readmes Academic
	if db.Where("lower(name) = ?", strings.ToLower(name)).First(&readmes).RecordNotFound() {
		readmes.Name = name
		err := db.Create(&readmes).Error
		return &readmes, true, err
	}
	return &readmes, false, nil
}

// LoadStars loads the stars for a readmes
func (readmes *Academic) LoadStars(db *gorm.DB, match string) error {
	// Make sure readmes exists in database, or we will panic
	var existing Academic
	if db.Where("id = ?", readmes.ID).First(&existing).RecordNotFound() {
		return fmt.Errorf("Academic '%d' not found", readmes.ID)
	}

	if match != "" {
		var stars []Star
		db.Raw(`
			SELECT *
			FROM ACADEMICS S
			INNER JOIN STAR_ACADEMICS ST ON S.ID = ST.STAR_ID
			WHERE S.DELETED_AT IS NULL
			AND ST.ACADEMIC_ID = ?
			AND LOWER(S.FULL_NAME) LIKE ?
			ORDER BY S.FULL_NAME`,
			readmes.ID,
			fmt.Sprintf("%%%s%%", strings.ToLower(match))).Scan(&stars)
		readmes.Stars = stars
		return db.Error
	}
	return db.Model(readmes).Association("Stars").Find(&readmes.Stars).Error
}

// Rename renames a readmes -- new name must not already exist
func (readmes *Academic) Rename(db *gorm.DB, name string) error {
	// Can't rename to the same name
	if name == readmes.Name {
		return errors.New("You can't rename to the same name")
	}

	// If they're just changing case, allow. Otherwise, block the change
	if strings.ToLower(name) != strings.ToLower(readmes.Name) {
		existing, err := FindAcademicByName(db, name)
		if err != nil {
			return err
		}
		if existing != nil {
			return fmt.Errorf("Academic '%s' already exists", existing.Name)
		}
	}

	readmes.Name = name
	return db.Save(readmes).Error
}

// Delete deletes a readmes and disassociates it from any stars
func (readmes *Academic) Delete(db *gorm.DB) error {
	if err := db.Model(readmes).Association("Stars").Clear().Error; err != nil {
		return err
	}
	return db.Delete(readmes).Error
}