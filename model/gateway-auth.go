package model

// https://github.com/euclid1990/gstats/blob/develop/utilities/google.go
// https://github.com/euclid1990/gstats/blob/develop/utilities/redmine.go
// https://github.com/euclid1990/gstats/blob/develop/utilities/github.go
// https://github.com/hfurubotten/autograder/blob/master/global/global.go

// git subtree add --prefix ./shared/models/tensorflow https://github.com/tensorflow/models master --squash

import (
	//"fmt"
	//"strings"
	//"time"
	"github.com/jinzhu/gorm"
	//"github.com/sirupsen/logrus"
)

/*
	Auth connectors definitions
*/
type GatewayAuth_Credentials struct {
	gorm.Model
	Disable         bool   `default:"false" gorm:"column:disable" yaml:"disable,omitempty" json:"disable,omitempty"`
	Provider        string `gorm:"column:provider" json:"provider,omitempty" yaml:"provider,omitempty"`
	PersonalToken   string `gorm:"column:auth_personal_token" json:"auth_personal_token,omitempty" yaml:"auth_personal_token,omitempty"`
	ClientID        string `gorm:"column:auth_client_id" json:"auth_client_id,omitempty" yaml:"auth_client_id,omitempty"`
	ClientSecretKey string `gorm:"column:auth_client_secret" json:"auth_client_secret,omitempty" yaml:"auth_client_secret,omitempty"`
	LoginURL        string `gorm:"column:auth_login_url" json:"auth_login_url,omitempty" yaml:"auth_login_url,omitempty"`
	CallbackURL     string `gorm:"column:auth_callback_url" json:"auth_callback_url,omitempty" yaml:"auth_callback_url,omitempty"`
}
