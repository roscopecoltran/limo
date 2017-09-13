package model

import (
	//"fmt"
	//"strings"
	//"time"
	// "github.com/jinzhu/gorm"
	//"github.com/sirupsen/logrus"
)

/*
	Entities classification
*/
type GatewayClassification struct {
	Disable 			bool 						`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	Topics 				[]Topic 					`json:"topics,omitempty" yaml:"topics,omitempty"`
	Tags 				[]Tag 						`json:"tags,omitempty" yaml:"tags,omitempty"`	
}

