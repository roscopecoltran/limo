package vcs

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
type Mention struct {
	gorm.Model
	Name  			string
	Slug  			string
	ContextTags 	[]string
	Description 	*string
	URL         	*string
	DocType 		string
}
