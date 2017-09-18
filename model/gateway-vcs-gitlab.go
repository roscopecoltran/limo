package model

//import (
//"fmt"
//"strings"
//"time"
// "github.com/jinzhu/gorm"
// "github.com/xanzy/go-gitlab"
//"github.com/sirupsen/logrus"
//)

/*
	VCS - gitlab bucket
*/
type GatewayBucket_Gitlab struct {
	// gorm.Model
	Disable bool `default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	PerPage int  `yaml:"per_page,omitempty" json:"per_page,omitempty"` //
}

type GatewayBucket_GitlabSearchOptions Global_SearchOptions
