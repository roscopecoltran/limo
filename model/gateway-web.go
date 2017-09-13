package model

import (
	//"fmt"
	//"strings"
	//"time"
	"github.com/jinzhu/gorm"
	//"github.com/sirupsen/logrus"
)

/*
	Website definitions
*/
type Web_SiteProfile struct {
	gorm.Model
	Disable 			bool 						`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	HomePageURL     	string 						`json:"home_page_url,omitempty" yaml:"home_page_url,omitempty"`
	FaqPageURL     		string 						`json:"faq_page_url,omitempty" yaml:"faq_page_url,omitempty"`
	ForumPageURL    	string 						`json:"forum_page_url,omitempty" yaml:"forum_page_url,omitempty"`
	BaseURL         	string 						`json:"base_url,omitempty" yaml:"base_url,omitempty"`
	Engine          	string 						`default:"xpath" json:"engine,omitempty" yaml:"engine,omitempty"`
	Topics      		[]string 					`json:"topics,omitempty" yaml:"topics,omitempty"`
	Categories      	[]string 					`json:"categories,omitempty" yaml:"categories,omitempty"`
}

/*
	WebUI definitions
*/
type Web_UIOptions struct {
	// gorm.Model
	Disable 			bool 						`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	DefaultLocale 		string 						`default:"en" json:"default_locale,omitempty" yaml:"default_locale,omitempty"`
	DefaultTheme  		string 						`default:"sniperkit" json:"default_theme,omitempty" yaml:"default_theme,omitempty"`
	StaticPath    		string 						`default:"" json:"static_path,omitempty" yaml:"static_path,omitempty"`
	TemplatesPath 		string 						`default:"" json:"templates_path,omitempty" yaml:"templates_path,omitempty"`
}
