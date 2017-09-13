package model

import (
	//"fmt"
	//"strings"
	//"time"
	//"github.com/jinzhu/gorm"
	//"github.com/sirupsen/logrus"
)

/*
	Tasks definition
*/

type Task_Options struct {
	Disable 			bool 						`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	Name 				string 						`json:"name,omitempty" yaml:"name,omitempty"`
	Disabled 			bool 						`default:"false" json:"disabled,omitempty" yaml:"disabled,omitempty"`
	Sanitize 			bool 						`default:"true" json:"sanitize,omitempty" yaml:"sanitize,omitempty"`
	Language 			string 						`default:"shell" json:"language,omitempty" yaml:"language,omitempty"` 	// shell, go, perl, python
	Script 				string 						`json:"script,omitempty" yaml:"script,omitempty"`
	Hooks 				[]string 					`json:"hooks,omitempty" yaml:"hooks,omitempty"`
	Tags 				[]string 					`json:"tags,omitempty" yaml:"tags,omitempty"`
	FlowBased 			string 						`json:"fbp,omitempty" yaml:"fbp,omitempty"`
}

