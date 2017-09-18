package model

import (
	"errors"
	"github.com/jinzhu/gorm"
	//"github.com/qor/sorting"
	"github.com/kylelemons/godebug/pretty"
	"github.com/sirupsen/logrus"
)

// https://github.com/yoru9zine/starlink/blob/master/main.go
// https://github.com/importre/mecca/blob/master/polymer.go

// Service represents a hosting service like Github
type Language struct {
	gorm.Model `json:"-" yaml:"-"`
	//sorting.SortingDESC
	Name         string         `gorm:"type:varchar(128);not null" yaml:"name,omitempty" json:"name,omitempty"`
	ServiceID    uint           `gorm:"column:service_id" yaml:"service_id,omitempty" json:"service_id,omitempty"`
	ByteCode     int            `gorm:"column:byte_code" yaml:"byte_code,omitempty" json:"byte_code,omitempty"`
	Count        int            `gorm:"-" yaml:"count,omitempty" json:"count,omitempty"`
	LanguageType []LanguageType `gorm:"many2many:language_type;" yaml:"language_type,omitempty" json:"language_type,omitempty"` // is a dev language, human language or sign language ?!
}

type LanguageResult struct {
	Language *Language
	Error    error
}

func NewLanguageFromGithub(langs map[string]int, remoteId int, userId int, remoteUri string) ([]Language, error) {
	log.WithFields(
		logrus.Fields{"service": "NewLanguageFromGithub",
			"languages": langs}).Info("")
	if len(langs) > 0 {
		log.WithFields(
			logrus.Fields{"service": "NewLanguageFromGithub",
				"step":  "check_language_empty",
				"langs": langs}).Infof("language detected")
	} else {
		err := errors.New("Language map is empty.")
		log.WithError(err).WithFields(
			logrus.Fields{"service": "NewLanguageFromGithub",
				"step": "check_language_empty"}).Warn("")
		return nil, err
	}
	var languages []Language
	if len(langs) > 0 {
		for langName, byteCode := range langs {
			languages = append(languages, Language{Name: langName, ByteCode: byteCode})
		}
		log.WithFields(
			logrus.Fields{"remoteUri": remoteUri,
				"languages": languages}).Warn("")
	}
	pretty.Print(languages)
	return languages, nil
}

// CreateOrUpdateLanguage creates or updates a language and returns true if the star was created (vs updated)
func CreateOrUpdateLanguage(db *gorm.DB, language *Language, service *Service) (bool, error) {
	// Get existing by remote ID and service ID
	var existing Language
	if db.Where("service_id = ?", service.ID).First(&existing).RecordNotFound() {
		language.ServiceID = service.ID
		err := db.Create(language).Error
		return err == nil, err
	}
	language.ID = existing.ID
	language.ServiceID = service.ID
	language.CreatedAt = existing.CreatedAt
	return false, db.Save(language).Error
}
