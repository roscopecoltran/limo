package model

import (
	//"fmt"
	//"strings"
	//"time"
	// "github.com/jinzhu/gorm"
	"github.com/google/go-github/github"
	"github.com/xanzy/go-gitlab"
	//"github.com/sirupsen/logrus"
)

/*
	Catalogs - vcs catalogs
*/
type GatewayCatalog_Github struct {
	Disable 			bool 								`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"` // Catalog status
	Catalog				GatewayBucket_Github 				`yaml:"catalog,omitempty" json:"catalog,omitempty"`					// Github API v3 - bucket of responses
	User 				*github.User 						`yaml:"user,omitempty" json:"user,omitempty"`						// https://github.com/google/go-github/blob/master/github/users.go#L20-L68
}

type GatewayCatalog_Gitlab struct {
	Disable 			bool 								`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`	// Catalog status
  	Catalog 			GatewayBucket_Gitlab  				`yaml:"catalog,omitempty" json:"catalog,omitempty"` 				// Gitlab API - bucket of responses
	User 				*gitlab.User 						`yaml:"user,omitempty" json:"user,omitempty"`						// https://github.com/google/go-github/blob/master/github/users.go#L20-L68
}

type GatewayCatalog_Bitbucket struct {
	Disable 			bool 								`default:"true" yaml:"disable,omitempty" json:"disable,omitempty"`	// Catalog status
}