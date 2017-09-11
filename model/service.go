package model

import (
	"github.com/jinzhu/gorm"
    // "github.com/qor/qor"
    // "github.com/qor/admin"
	// "github.com/sirupsen/logrus"
)

// https://github.com/yoru9zine/starlink/blob/master/main.go
// https://github.com/importre/mecca/blob/master/polymer.go
// 

// Service represents a hosting service like Github
type Service struct {
	gorm.Model
	Name  		string
	Description *string
	Homepage    *string
	URL         *string
	Stars 		[]Star
}

// FindOrCreateServiceByName returns a service with the specified name, creating if necessary
func FindOrCreateServiceByName(db *gorm.DB, name string) (*Service, bool, error) {
	var service Service
	if db.Where("name = ?", name).First(&service).RecordNotFound() {
		service.Name = name
		err := db.Create(&service).Error
		return &service, true, err
	}
	return &service, false, nil
}
