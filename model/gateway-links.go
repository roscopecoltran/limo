package model

import (
	//"fmt"
	//"strings"
	//"time"
	//"github.com/jinzhu/gorm"
	//"github.com/sirupsen/logrus"
)

/*
	Links definition
*/
type GatewayLink struct {
	Disable 			bool 						 	`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	Title 				string 						 	`json:"title,omitempty" yaml:"title,omitempty"`
	Href 				string 						 	`json:"href,omitempty" yaml:"href,omitempty"`
	Meta 				GatewayClassification 			`json:"meta,omitempty" yaml:"meta,omitempty"`
	Attr 				GatewayLink_AttributesOptions	`json:"attributes,omitempty" yaml:"attributes,omitempty"`
	Metrics 			GatewayLink_MetricsOptions 	`json:"metrics,omitempty" yaml:"metrics,omitempty"`
}

type GatewayLink_RanksOptions struct {
	Disable 			bool 							`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	Global 				map[string]float64 				`json:"global,omitempty" yaml:"global,omitempty"` 			// alexa, google, bing
	Internal 			float64 							`json:"internal,omitempty" yaml:"internal,omitempty"`
	Page 				map[string]float64 				`json:"page,omitempty" yaml:"page,omitempty"`
	Site 				map[string]float64 				`json:"site,omitempty" yaml:"site,omitempty"`
}

type GatewayLink_BoundsOptions struct {
	Disable 			bool 							`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	Inbounds 			map[string]int 					`json:"inbounds,omitempty" yaml:"inbounds,omitempty"` 		// web, social, ecommerce, knowledge
	Outbounds 			map[string]int 					`json:"outbounds,omitempty" yaml:"outbounds,omitempty"` 	// web, social, ecommerce, knowledge
}

type GatewayLink_MetricsOptions struct {
	Disable 			bool 							`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	Weight          	int  							`default:"1" json:"weight,omitempty" yaml:"weight,omitempty"`
	Ranks 				GatewayLink_RanksOptions 		`json:"ranks,omitempty" yaml:"ranks,omitempty"`
	Links 				GatewayLink_BoundsOptions 		`json:"links,omitempty" yaml:"links,omitempty"`
}

type GatewayLink_AttributesOptions struct {
	Disable 			bool 							`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	Target 				string 							`json:"target,omitempty" yaml:"target,omitempty"`
	HTML5 				map[string]string 				`json:"data,omitempty" yaml:"data,omitempty"`
}