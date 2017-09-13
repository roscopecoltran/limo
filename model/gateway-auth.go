package model

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
	Disable 			bool 						`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	Provider   			string 						`json:"provider,omitempty" yaml:"provider,omitempty"`
	PersonalToken   	string 						`json:"auth_personal_token,omitempty" yaml:"auth_personal_token,omitempty"`
	ClientID       		string 						`json:"auth_client_id,omitempty" yaml:"auth_client_id,omitempty"`
	ClientSecretKey 	string 						`json:"auth_client_secret,omitempty" yaml:"auth_client_secret,omitempty"`
	LoginURL       		string 						`json:"auth_login_url,omitempty" yaml:"auth_login_url,omitempty"`
	CallbackURL     	string 						`json:"auth_callback_url,omitempty" yaml:"auth_callback_url,omitempty"`
}

