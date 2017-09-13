package model

import (
	//"fmt"
	//"strings"
	//"time"
	//"github.com/jinzhu/gorm"
	//"github.com/sirupsen/logrus"
)

/*
	Gateway aggregation neurons
*/
type GatewayNeurons struct {
	// gorm.Model
	Disable 			bool 						`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	Name 				string 						`json:"name,omitempty" yaml:"name,omitempty"`
	Topics 				[]string 					`json:"topics,omitempty" yaml:"topics,omitempty"`
	Channels 			map[string]string 			`json:"chanels,omitempty" yaml:"chanels,omitempty"`
}
