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

// Pkg represents a pkg in the database
type Pkg struct {
	gorm.Model
	Name      	string
	PkgCount 	int    `gorm:"-"`
	StarCount 	int    `gorm:"-"`
	Stars     	[]Star `gorm:"many2many:star_pkgs;"`
}

// FindPkgs finds all pkgs
func FindPkgs(db *gorm.DB) ([]Pkg, error) {
	var pkgs []Pkg
	db.Order("name").Find(&pkgs)
	return pkgs, db.Error
}

// FindPkgsWithStarCount finds all pkgs and gets their count of stars
func FindPkgsWithStarCount(db *gorm.DB) ([]Pkg, error) {
	var pkgs []Pkg

	// Create resources from GORM-backend model
	// Admin.AddResource(&Pkg{})

	rows, err := db.Raw(`
		SELECT T.NAME, COUNT(ST.PKG_ID) AS STARCOUNT
		FROM PKGS T
		LEFT JOIN STAR_PKGS ST ON T.ID = ST.PKG_ID
		WHERE T.DELETED_AT IS NULL
		GROUP BY T.ID
		ORDER BY T.NAME`).Rows()

	if err != nil {
		return pkgs, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	for rows.Next() {
		var pkg Pkg
		if err = rows.Scan(&pkg.Name, &pkg.StarCount); err != nil {
			return pkgs, err
		}
		pkgs = append(pkgs, pkg)
	}
	return pkgs, db.Error
}

// FindPkgByName finds a pkg by name
func FindPkgByName(db *gorm.DB, name string) (*Pkg, error) {
	var pkg Pkg
	if db.Where("lower(name) = ?", strings.ToLower(name)).First(&pkg).RecordNotFound() {
		return nil, db.Error
	}
	return &pkg, db.Error
}

// FindOrCreatePkgByName finds a pkg by name, creating if it doesn't exist
func FindOrCreatePkgByName(db *gorm.DB, name string) (*Pkg, bool, error) {
	var pkg Pkg
	if db.Where("lower(name) = ?", strings.ToLower(name)).First(&pkg).RecordNotFound() {
		pkg.Name = name
		err := db.Create(&pkg).Error
		return &pkg, true, err
	}
	return &pkg, false, nil
}

// LoadStars loads the stars for a pkg
func (pkg *Pkg) LoadStars(db *gorm.DB, match string) error {
	// Make sure pkg exists in database, or we will panic
	var existing Pkg
	if db.Where("id = ?", pkg.ID).First(&existing).RecordNotFound() {
		return fmt.Errorf("Pkg '%d' not found", pkg.ID)
	}

	if match != "" {
		var stars []Star
		db.Raw(`
			SELECT *
			FROM STARS S
			INNER JOIN STAR_PKGS ST ON S.ID = ST.STAR_ID
			WHERE S.DELETED_AT IS NULL
			AND ST.PKG_ID = ?
			AND LOWER(S.FULL_NAME) LIKE ?
			ORDER BY S.FULL_NAME`,
			pkg.ID,
			fmt.Sprintf("%%%s%%", strings.ToLower(match))).Scan(&stars)
		pkg.Stars = stars
		return db.Error
	}
	return db.Model(pkg).Association("Stars").Find(&pkg.Stars).Error
}

// Rename renames a pkg -- new name must not already exist
func (pkg *Pkg) Rename(db *gorm.DB, name string) error {
	// Can't rename to the same name
	if name == pkg.Name {
		return errors.New("You can't rename to the same name")
	}

	// If they're just changing case, allow. Otherwise, block the change
	if strings.ToLower(name) != strings.ToLower(pkg.Name) {
		existing, err := FindPkgByName(db, name)
		if err != nil {
			return err
		}
		if existing != nil {
			return fmt.Errorf("Pkg '%s' already exists", existing.Name)
		}
	}

	pkg.Name = name
	return db.Save(pkg).Error
}

// Delete deletes a pkg and disassociates it from any stars
func (pkg *Pkg) Delete(db *gorm.DB) error {
	if err := db.Model(pkg).Association("Stars").Clear().Error; err != nil {
		return err
	}
	return db.Delete(pkg).Error
}
