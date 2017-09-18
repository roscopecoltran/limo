package model

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"strings"
)

// Tag represents a tag in the database
type Tag struct {
	gorm.Model `json:"-" yaml:"-"`
	Name       string `gorm:"type:varchar(64);not null;column:name" yaml:"name,omitempty" json:"name,omitempty"`
	StarCount  int    `gorm:"-" yaml:"count,omitempty" json:"count,omitempty"`
	Stars      []Star `gorm:"many2many:star_tags;" yaml:"star_tags,omitempty" json:"star_tags,omitempty"`
}

// FindTags finds all tags
func FindTags(db *gorm.DB) ([]Tag, error) {
	var tags []Tag
	db.Order("name").Find(&tags)
	return tags, db.Error
}

// FindTagsWithStarCount finds all tags and gets their count of stars
func FindTagsWithStarCount(db *gorm.DB) ([]Tag, error) {
	var tags []Tag
	rows, err := db.Raw(`
		SELECT T.NAME, COUNT(ST.TAG_ID) AS STARCOUNT
		FROM TAGS T
		LEFT JOIN STAR_TAGS ST ON T.ID = ST.TAG_ID
		WHERE T.DELETED_AT IS NULL
		GROUP BY T.ID
		ORDER BY T.NAME`).Rows()
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"section:": "model", "typology": "tag", "step": "FindTagsWithStarCount"}).Warnf("%#s", err)
		return tags, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.WithError(err).WithFields(logrus.Fields{"section:": "model", "typology": "tag", "step": "FindTagsWithStarCount"}).Fatalf("%#s", err)
			//log.Fatal(err)
		}
	}()

	for rows.Next() {
		var tag Tag
		if err = rows.Scan(&tag.Name); err != nil {
			log.WithError(err).WithFields(logrus.Fields{"section:": "model", "typology": "tag", "step": "FindTagsWithStarCount"}).Warnf("%#s", err)
			return tags, err
		}
		tags = append(tags, tag)
	}
	return tags, db.Error
}

// FindTagByName finds a tag by name
func FindTagByName(db *gorm.DB, name string) (*Tag, error) {
	var tag Tag
	if db.Where("lower(name) = ?", strings.ToLower(name)).First(&tag).RecordNotFound() {
		return nil, db.Error
	}
	return &tag, db.Error
}

// FindOrCreateTagByName finds a tag by name, creating if it doesn't exist
func FindOrCreateTagByName(db *gorm.DB, name string) (*Tag, bool, error) {
	var tag Tag
	if db.Where("lower(name) = ?", strings.ToLower(name)).First(&tag).RecordNotFound() {
		tag.Name = name
		err := db.Create(&tag).Error
		log.WithError(err).WithFields(logrus.Fields{"section:": "model", "typology": "tag", "step": "FindOrCreateTagByName"}).Warnf("%#s", err)
		return &tag, true, err
	}
	return &tag, false, nil
}

// LoadStars loads the stars for a tag
func (tag *Tag) LoadStars(db *gorm.DB, match string) error {
	// Make sure tag exists in database, or we will panic
	var existing Tag
	if db.Where("id = ?", tag.ID).First(&existing).RecordNotFound() {
		return fmt.Errorf("Tag '%d' not found", tag.ID)
	}

	if match != "" {
		var stars []Star
		db.Raw(`
			SELECT *
			FROM STARS S
			INNER JOIN STAR_TAGS ST ON S.ID = ST.STAR_ID
			WHERE S.DELETED_AT IS NULL
			AND ST.TAG_ID = ?
			AND LOWER(S.FULL_NAME) LIKE ?
			ORDER BY S.FULL_NAME`,
			tag.ID,
			fmt.Sprintf("%%%s%%", strings.ToLower(match))).Scan(&stars)
		tag.Stars = stars
		return db.Error
	}
	return db.Model(tag).Association("Stars").Find(&tag.Stars).Error
}

// Rename renames a tag -- new name must not already exist
func (tag *Tag) Rename(db *gorm.DB, name string) error {
	// Can't rename to the same name
	if name == tag.Name {
		return errors.New("You can't rename to the same name")
	}

	// If they're just changing case, allow. Otherwise, block the change
	if strings.ToLower(name) != strings.ToLower(tag.Name) {
		existing, err := FindTagByName(db, name)
		if err != nil {
			log.WithError(err).WithFields(logrus.Fields{"section:": "model", "typology": "tag", "step": "Rename"}).Warnf("%#s", err)
			return err
		}
		if existing != nil {
			err := fmt.Errorf("Tag '%s' already exists", existing.Name)
			log.WithError(err).WithFields(logrus.Fields{"section:": "model", "typology": "tag", "step": "Rename"}).Errorf("Tag '%s' already exists", existing.Name)
			return err
		}
	}

	tag.Name = name
	return db.Save(tag).Error
}

// Delete deletes a tag and disassociates it from any stars
func (tag *Tag) Delete(db *gorm.DB) error {
	if err := db.Model(tag).Association("Stars").Clear().Error; err != nil {
		return err
	}
	return db.Delete(tag).Error
}
