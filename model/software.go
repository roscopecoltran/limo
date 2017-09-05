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

// Software represents a software in the database
type Software struct {
	gorm.Model
	Name      		string
	SoftwareCount 	int    `gorm:"-"`
	StarCount 		int    `gorm:"-"`
	Stars     		[]Star `gorm:"many2many:star_softwares;"`
}

// FindSoftwares finds all softwares
func FindSoftwares(db *gorm.DB) ([]Software, error) {
	var softwares []Software
	db.Order("name").Find(&softwares)
	return softwares, db.Error
}

// FindSoftwaresWithStarCount finds all softwares and gets their count of stars
func FindSoftwaresWithStarCount(db *gorm.DB) ([]Software, error) {
	var softwares []Software

	// Create resources from GORM-backend model
	// Admin.AddResource(&Software{})

	rows, err := db.Raw(`
		SELECT T.NAME, COUNT(ST.SOFTWARE_ID) AS STARCOUNT
		FROM SOFTWARES T
		LEFT JOIN STAR_SOFTWARES ST ON T.ID = ST.SOFTWARE_ID
		WHERE T.DELETED_AT IS NULL
		GROUP BY T.ID
		ORDER BY T.NAME`).Rows()

	if err != nil {
		return softwares, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	for rows.Next() {
		var software Software
		if err = rows.Scan(&software.Name, &software.StarCount); err != nil {
			return softwares, err
		}
		softwares = append(softwares, software)
	}
	return softwares, db.Error
}

// FindSoftwareByName finds a software by name
func FindSoftwareByName(db *gorm.DB, name string) (*Software, error) {
	var software Software
	if db.Where("lower(name) = ?", strings.ToLower(name)).First(&software).RecordNotFound() {
		return nil, db.Error
	}
	return &software, db.Error
}

// FindOrCreateSoftwareByName finds a software by name, creating if it doesn't exist
func FindOrCreateSoftwareByName(db *gorm.DB, name string) (*Software, bool, error) {
	var software Software
	if db.Where("lower(name) = ?", strings.ToLower(name)).First(&software).RecordNotFound() {
		software.Name = name
		err := db.Create(&software).Error
		return &software, true, err
	}
	return &software, false, nil
}

// LoadStars loads the stars for a software
func (software *Software) LoadStars(db *gorm.DB, match string) error {
	// Make sure software exists in database, or we will panic
	var existing Software
	if db.Where("id = ?", software.ID).First(&existing).RecordNotFound() {
		return fmt.Errorf("Software '%d' not found", software.ID)
	}

	if match != "" {
		var stars []Star
		db.Raw(`
			SELECT *
			FROM SOFTWARES S
			INNER JOIN STAR_SOFTWARES ST ON S.ID = ST.STAR_ID
			WHERE S.DELETED_AT IS NULL
			AND ST.SOFTWARE_ID = ?
			AND LOWER(S.FULL_NAME) LIKE ?
			ORDER BY S.FULL_NAME`,
			software.ID,
			fmt.Sprintf("%%%s%%", strings.ToLower(match))).Scan(&stars)
		software.Stars = stars
		return db.Error
	}
	return db.Model(software).Association("Stars").Find(&software.Stars).Error
}

// Rename renames a software -- new name must not already exist
func (software *Software) Rename(db *gorm.DB, name string) error {
	// Can't rename to the same name
	if name == software.Name {
		return errors.New("You can't rename to the same name")
	}

	// If they're just changing case, allow. Otherwise, block the change
	if strings.ToLower(name) != strings.ToLower(software.Name) {
		existing, err := FindSoftwareByName(db, name)
		if err != nil {
			return err
		}
		if existing != nil {
			return fmt.Errorf("Software '%s' already exists", existing.Name)
		}
	}

	software.Name = name
	return db.Save(software).Error
}

// Delete deletes a software and disassociates it from any stars
func (software *Software) Delete(db *gorm.DB) error {
	if err := db.Model(software).Association("Stars").Clear().Error; err != nil {
		return err
	}
	return db.Delete(software).Error
}
