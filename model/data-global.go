package model

import (
	//"fmt"
	//"strings"
	//"time"
	"github.com/jinzhu/gorm"
	//"github.com/sirupsen/logrus"
)

/*
	General definitions
*/
type Global_Profile struct {
	gorm.Model
	Debug        		bool   						`default:"false" json:"debug,omitempty" yaml:"debug,omitempty"`
	InstanceName 		string 						`default:"sniperkit" json:"instance_name,omitempty" yaml:"instance_name,omitempty"`
	Language     		string 						`default:"en" json:"language,omitempty" yaml:"language,omitempty"`
}

type Global_SearchOptions struct {
	// gorm.Model
	Disable 		bool 							`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	Offset  		int 							`default:"1" yaml:"offset,omitempty" json:"offset,omitempty"`			// 
	MinPage  		int 							`default:"1" yaml:"min_page,omitempty" json:"min_page,omitempty"`		// 
	MaxPage  		int 							`default:"100" yaml:"max_page,omitempty" json:"max_page,omitempty"`		// 
	PerPage  		int 							`default:"100" yaml:"per_page,omitempty" json:"per_page,omitempty"`		// 
	Sort 			string 							`default:"updated" yaml:"sort,omitempty" json:"sort,omitempty"` 		// Default is to sort by best match.
	Order 			string 							`default:"desc" yaml:"order,omitempty" json:"order,omitempty"` 			// Sort order if sort parameter is provided. Possible values are: asc, desc. Default is desc.
	TextMatch 		bool 							`yaml:"-" json:"-"`	// Whether to retrieve text match metadata with a query
}

type Global_OutgoingOptions struct {
	// gorm.Model
	Disable 			bool 						`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	PoolConnections 	int  						`default:"100" json:"pool_connections,omitempty" yaml:"pool_connections,omitempty"`
	PoolMaxsize     	int  						`default:"10" json:"pool_maxsize,omitempty" yaml:"pool_maxsize,omitempty"`
	RequestTimeout  	int  						`default:"2.0" json:"request_timeout,omitempty" yaml:"request_timeout,omitempty"`
	UseragentSuffix 	string 						`default:"Sniperkit-X" json:"useragent_suffix,omitempty" yaml:"useragent_suffix,omitempty"`
}

type Global_ServerOptions struct {
	// gorm.Model
	Disable 			bool 						`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	BaseURL             bool   						`default:"false" json:"base_url" yaml:"base_url"`
	BindAddress         string 						`default:"127.0.0.1" json:"bind_address" yaml:"bind_address"`
	HTTPProtocolVersion string 						`default:"1.1" json:"http_protocol_version" yaml:"http_protocol_version"`
	ImageProxy          bool   						`default:"false" json:"image_proxy" yaml:"image_proxy"`
	Port                int  						`default:"8888" json:"port" yaml:"port"`
	SecretKey           string 						`default:"ultradeepsecretkey" json:"secret_key" yaml:"secret_key"`
} 