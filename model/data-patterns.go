package model

import (
	//"fmt"
	//"strings"
	//"time"
	"github.com/jinzhu/gorm"
	//"github.com/sirupsen/logrus"
)

/*
	Patterns - Ignore(s)
*/
type Patterns_Ignore struct {
	// gorm.Model
	// A list of URLs that the bot can ignore.
	Words 			[]Patterns_IgnoreRules 					`yaml:"content,omitempty" json:"content,omitempty"` 			//
	Links 			[]Patterns_IgnoreRules 					`yaml:"links,omitempty" json:"links,omitempty"` 				//
}

type Patterns_IgnoreRules struct {
	gorm.Model
	Disable 		bool 									`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	Mode 			string 									`default:"regex" yaml:"mode,omitempty" json:"mode,omitempty"` 	// regex, strict, contains
	Attribute 		string 									`yaml:"attribute,omitempty" json:"attribute,omitempty"` 		// 
	Pattern 		string 									`yaml:"pattern,omitempty" json:"pattern,omitempty"` 			// 
	Timeout 		string 									`yaml:"timeout,omitempty" json:"timeout,omitempty"` 			//
}

/*
	Patterns - Word Lists
*/
type Patterns_WordLists struct {
	// gorm.Model
	Disable 		bool 									`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	TargetLanaguge 	string 									`yaml:"target_language,omitempty" json:"target_language,omitempty"` 
	DorksList 		[]string 								`yaml:"dorks_list,omitempty" json:"dorks_list,omitempty"` 		// 
	BlackList 		[]string 								`yaml:"black_list,omitempty" json:"black_list,omitempty"` 		// 
	WhiteList 		[]string 								`yaml:"white_list,omitempty" json:"white_list,omitempty"` 		// 
	StopLists 		[]string 								`yaml:"stop_list,omitempty" json:"stop_list,omitempty"` 		// 
	CommonWords 	[]string 								`yaml:"common_words,omitempty" json:"mode,omitempty"` 			//
}
