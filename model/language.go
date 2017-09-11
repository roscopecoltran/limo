package model

import (
	"github.com/jinzhu/gorm"
	//"fmt"
	"errors"
    // "github.com/qor/qor"
    // "github.com/qor/admin"
	"github.com/sirupsen/logrus"
)

// https://github.com/yoru9zine/starlink/blob/master/main.go
// https://github.com/importre/mecca/blob/master/polymer.go
// 

// Service represents a hosting service like Github
type Language struct {
	gorm.Model
	Name  			string `gorm:"type:varchar(128);not null;unique"`
	ServiceID   	uint   
	Count 			int    `gorm:"-"` 
}

type LanguageResult struct {
	Language  	*Language
	Error 	error
}

func NewLanguageFromGithub(language string) (*Language, error) {
	log.WithFields(logrus.Fields{"service": "NewLanguageFromGithub", "language": language}).Info("")
	if language == "" {
		return nil, errors.New("Language is empty.")
	}
	languageInfo :=	&Language{
		Name:    	 language,
	}
	return languageInfo, nil
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

