package model

import (
	"github.com/jinzhu/gorm"
	"errors"
	"github.com/sirupsen/logrus"
	//"strings"
	"github.com/kylelemons/godebug/pretty"
)

// https://github.com/yoru9zine/starlink/blob/master/main.go
// https://github.com/importre/mecca/blob/master/polymer.go

// Service represents a hosting service like Github
type Language struct {
	gorm.Model
	Name  			string 			`gorm:"type:varchar(128);not null;unique"`
	ServiceID   	uint   
	ByteCode 		int    			 
	Count 			int    			`gorm:"-"` 
	LanguageType    []LanguageType 	`gorm:"many2many:language_type;"` 			// is a dev language, human language or sign language ?!
}

type LanguageResult struct {
	Language  	*Language
	Error 		error
}

func NewLanguageFromGithub(langs map[string]int, remoteId int, userId int, remoteUri string) ([]Language, error) {
	log.WithFields(
		logrus.Fields{	"service": 		"NewLanguageFromGithub", 
						"languages": 	langs}).Info("")
	if len(langs) > 0 {
		log.WithFields(
			logrus.Fields{	"service": 	"NewLanguageFromGithub", 
							"step": 	"check_language_empty", 
							"langs": 	langs}).Infof("language detected")		
	} else {
		err := errors.New("Language map is empty.")
		log.WithError(err).WithFields(
			logrus.Fields{	"service": 	"NewLanguageFromGithub", 
							"step": 	"check_language_empty"}).Warn("")
		return nil, err
	}
	var languages []Language
	if len(langs) > 0 {
		for langName, byteCode := range langs {
			languages = append(languages, Language{Name: langName, ByteCode: byteCode})
		}
		log.WithFields(
			logrus.Fields{	"remoteUri": remoteUri, 
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

