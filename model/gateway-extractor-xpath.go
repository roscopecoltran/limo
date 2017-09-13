package model

import (
	//"fmt"
	//"strings"
	//"time"
	//"github.com/jinzhu/gorm"
	//"github.com/sirupsen/logrus"
)

/*
	Gateway - XPath endpoints
	- 1. Endpoint
	- 2. Request
	- 3. RequestOptions
	- 4. Response
	- 5. ResponseOptions
	- 6. Mapping
	- 7. Items
	- 8. Item
	- 9. Selectors
*/
type GatewayExt_XPath_Endpoint struct {
	// gorm.Model
	Disable 			bool 							`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	Options 			GatewayGlobal 					`json:"options,omitempty" yaml:"options,omitempty"`
	Request 			GatewayExt_XPath_Request		`json:"request,omitempty" yaml:"request,omitempty"`		
	Response 			GatewayExt_XPath_Response 		`json:"results,omitempty" yaml:"results,omitempty"`
}

type GatewayExt_XPath_Request struct {
	// gorm.Model
	Disable 			bool 							`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	Tasks 				map[string]Task_Options 		`json:"tasks,omitempty" yaml:"tasks,omitempty"`
	Request				GatewayGlobal_Request 			`json:"request,omitempty" yaml:"request,omitempty"`
}

type GatewayExt_XPath_Response struct {
	// gorm.Model
	Disable 			bool 							`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	Tasks 				map[string]Task_Options 		`json:"tasks,omitempty" yaml:"tasks,omitempty"`
	Mapping 			GatewayExt_XPath_Mapping 		`json:"mapping,omitempty" yaml:"mapping,omitempty"`
}

type GatewayExt_XPath_Mapping struct {
	// gorm.Model
	Disable 			bool 							`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	Items  				GatewayExt_XPath_Items 		`json:"items,omitempty" yaml:"items,omitempty"`
	Item 				GatewayExt_XPath_Item 			`json:"item,omitempty" yaml:"item,omitempty"`
}

type GatewayExt_XPath_Items struct {
	// gorm.Model
	Disable 			bool 							`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	Selector  			string 							`json:"selector,omitempty" yaml:"selector,omitempty"`
	Tasks 				map[string]Task_Options 		`json:"tasks,omitempty" yaml:"tasks,omitempty"`
}

type GatewayExt_XPath_Item struct {
	// gorm.Model
	Disable 			bool 							`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	PrefixName 			string 							`json:"prefix_name,omitempty" yaml:"prefix_name,omitempty"`
	PrefixSlug 			string 							`json:"prefix_slug,omitempty" yaml:"prefix_slug,omitempty"`
	// processing content
	PreProcess  		bool 							`json:"pre_process,omitempty" yaml:"pre_process,omitempty"`
	PostProcess  		bool 							`json:"post_process,omitempty" yaml:"post_process,omitempty"`
	Tasks 				map[string]Task_Options 		`json:"tasks,omitempty" yaml:"tasks,omitempty"`
	Sanitize  			bool 							`json:"sanitize,omitempty" yaml:"sanitize,omitempty"`
	Selectors 			GatewayExt_XPath_Selectors 	`json:"selectors,omitempty" yaml:"selectors,omitempty"`
}

type GatewayExt_XPath_Selectors struct {
	// Parse XML doc patterns
	Title      			string 							`json:"title,omitempty" yaml:"title,omitempty"`
	Result  			string 							`json:"result,omitempty" yaml:"result,omitempty"`
	Content 	    	string 							`json:"content,omitempty" yaml:"content,omitempty"`
	Tags 	    		[]Tag 							`json:"tags,omitempty" yaml:"tags,omitempty"`
	Extra 	    		string 							`json:"extra,omitempty" yaml:"extra,omitempty"`
	URL 	       		GatewayLink 					`json:"url,omitempty" yaml:"url,omitempty"`
}
