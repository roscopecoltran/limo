package model

import (
	//"fmt"
	//"strings"
	//"time"
	"github.com/jinzhu/gorm"
	//"github.com/sirupsen/logrus"
)

/*
	Domain discovery definitions
*/
type GatewayDomainDiscovery struct {
	gorm.Model
	Profile 			GatewayDomainDiscovery_Profile 		`yaml:"profile,omitempty" json:"profile,omitempty"`
	Methods 			GatewayDomainDiscovery_Methods 		`yaml:"methods,omitempty" json:"methods,omitempty"`
	Providers 			GatewayDomainDiscovery_Providers 	`yaml:"providers,omitempty" json:"providers,omitempty"`
} 

type GatewayDomainDiscovery_Providers struct {
	// gorm.Model
	Github 				GatewayBucket_Github 				`yaml:"github,omitempty" json:"github,omitempty"`
	Gitlab 				GatewayBucket_Gitlab				`yaml:"gitlab,omitempty" json:"gitlab,omitempty"`
	//Bitbucket 		GatewayExt_Restful_Bitbucket		`yaml:"bitbucket,omitempty" json:"bitbucket,omitempty"`
	Meta 				GatewayMeta 						`yaml:"sniperkit,omitempty" json:"sniperkit,omitempty"`
	//Dynamic 			GatewayExt_Restful_Dynamic 			`yaml:"dynamic,omitempty" json:"dynamic,omitempty"`	
}

type GatewayDomainDiscovery_Profile struct {
	// gorm.Model
	Disable 			bool 								`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	Username 			string 								`yaml:"username,omitempty" json:"username,omitempty"`			// 
	Email 				string 								`yaml:"email,omitempty" json:"email,omitempty"`					// 
	Organization 		string 								`yaml:"organization,omitempty" json:"organization,omitempty"`	// 
	Keywords			[]string 							`yaml:"keywords,omitempty" json:"keywords,omitempty"`
	Patterns			[]string 							`yaml:"patterns,omitempty" json:"patterns,omitempty"`
}

type GatewayDomainDiscovery_Methods struct {
	// gorm.Model
	Disable 			bool 								`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	Expansion 			GatewayDomainDiscovery_Expansion 	`yaml:"expansion,omitempty" json:"expansion,omitempty"`			// 
	Options 			[]Global_SearchOptions 				`yaml:"options,omitempty" json:"options,omitempty"` 			// 	
}

type GatewayDomainDiscovery_Expansion struct {
	// gorm.Model
	Disable 			bool 								`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	Methods 			map[string]string  					`yaml:"mehtods,omitempty" json:"mehtods,omitempty"`
}
