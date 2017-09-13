package model

import (
	//"fmt"
	//"strings"
	//"time"
	// "github.com/jinzhu/gorm"
	//"github.com/sirupsen/logrus"
)

/*
	Gateway - Restful APIs
*/
type GatewayExt_Restful_Endpoint struct {
	// gorm.Model
	Disable 			bool 								`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	SearchGlobalOptions GatewayGlobal 						`json:"global,omitempty" yaml:"global,omitempty"`
	Paging 				GatewayGlobal_PagingOptions 		`json:"paging,omitempty" yaml:"paging,omitempty"` 		// paging options
	Request 			GatewayExt_Restful_RequestOptions 	`json:"request,omitempty" yaml:"request,omitempty"` 	// request options
	Response 			GatewayExt_Restful_ResponseOptions 	`json:"response,omitempty" yaml:"response,omitempty"` 	// response received
	Ignores 			[]Patterns_IgnoreRules 				`yaml:"ignores,omitempty" json:"ignores,omitempty"` 	// ignore options
}

// GatewayGlobal_Request

type GatewayExt_Restful_RequestOptions struct {
	AutoSave  			bool 						`default:"true" json:"auto_save,omitempty" yaml:"auto_save,omitempty"`
	Disable 			bool 						`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	CacheTime  			int 						`default:"3600*24*3" json:"cache_time,omitempty" yaml:"cache_time,omitempty"`
	CacheValidate  		bool 						`default:"true" json:"cache_validate,omitempty" yaml:"cache_validate,omitempty"`
	CheckSSL        	bool 						`default:"true" json:"check_ssl,omitempty" yaml:"check_ssl,omitempty"`
	ExpirationTime  	int 						`default:"3600*24*30" json:"expiration_time,omitempty" yaml:"expiration_time,omitempty"`
	Headers 			map[string]string   		`json:"headers,omitempty" yaml:"headers,omitempty"`
	MimeType 			string   					`default:"json" json:"mimetype,omitempty" yaml:"mimetype,omitempty"` // json, multipart, xml
	PayLoad 			string   					`json:"payload,omitempty" yaml:"payload,omitempty"`
	Timeout  			float64 					`default:"2.0" json:"timeout,omitempty" yaml:"timeout,omitempty"`
}

type GatewayExt_Restful_Response struct {
	// gorm.Model
	Disable 			bool 						`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	Tasks 				map[string]Task_Options 	`json:"tasks,omitempty" yaml:"tasks,omitempty"`
	Mapping 			GatewayExt_Restful_Mapping 	`json:"mapping,omitempty" yaml:"mapping,omitempty"`
}

type GatewayExt_Restful_ResponseOptions struct {
	AutoSave  			bool 						`default:"true" json:"auto_save,omitempty" yaml:"auto_save,omitempty"`
	Disable 			bool 						`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	MimeType 			string   					`default:"json" json:"mimetype,omitempty" yaml:"mimetype,omitempty"`
	MinLength  			int 						`default:"1000" json:"min_lenght,omitempty" yaml:"min_lenght,omitempty"`
	MaxLength  			int 						`default:"50000000" json:"max_lenght,omitempty" yaml:"max_lenght,omitempty"`
	PostProcess  		bool 						`default:"false" json:"post_process,omitempty" yaml:"post_process,omitempty"`
	PreProcess  		bool 						`default:"false" json:"pre_process,omitempty" yaml:"pre_process,omitempty"`
	Sanitize  			bool 						`default:"true" json:"sanitize,omitempty" yaml:"sanitize,omitempty"`
	Timeout  			float64 					`default:"2.0" json:"timeout,omitempty" yaml:"timeout,omitempty"`
}

type GatewayExt_Restful_Mapping struct { 				// dynamic parsing, eg: outter.inner.value1, ref: https://github.com/Jeffail/gabs#parsing-and-searching-json
	// gorm.Model
	Disable 			bool 						`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	Title      			string 						`json:"title,omitempty" yaml:"title,omitempty"` 					// 
	Results    			string 						`json:"results,omitempty" yaml:"results,omitempty"` 				// 
	Items 				GatewayExt_Restful_Items 	`json:"items,omitempty" yaml:"items,omitempty"`
	Item 				GatewayExt_Restful_Item 	`json:"item,omitempty" yaml:"item,omitempty"`
	Content 			string 						`json:"content,omitempty" yaml:"content,omitempty"` 				// core content
	OutboundLink 		string 						`json:"outbound_link,omitempty" yaml:"outbound_link,omitempty"`		// main link to highlight
	OutboundLinks   	[]GatewayLink 				`json:"outbound_links,omitempty" yaml:"outbound_links,omitempty"`	// string list of links fetched
	Extra 				map[string]string 			`json:"extra,omitempty" yaml:"extra,omitempty"`						// 
}


type GatewayExt_Restful_Items struct {
	// gorm.Model
	Disable 			bool 						`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	Selector  			string 						`json:"selector,omitempty" yaml:"selector,omitempty"`
	Tasks 				map[string]Task_Options 	`json:"tasks,omitempty" yaml:"tasks,omitempty"`
}

type GatewayExt_Restful_Item struct {
	// gorm.Model
	Disable 			bool 						 `default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	PrefixName 			string 						 `json:"prefix_name,omitempty" yaml:"prefix_name,omitempty"`
	PrefixSlug 			string 						 `json:"prefix_slug,omitempty" yaml:"prefix_slug,omitempty"`
	// processing content
	PreProcess  		bool 						 `json:"pre_process,omitempty" yaml:"pre_process,omitempty"`
	PostProcess  		bool 						 `json:"post_process,omitempty" yaml:"post_process,omitempty"`
	Tasks 				map[string]Task_Options 	 `json:"tasks,omitempty" yaml:"tasks,omitempty"`
	Sanitize  			bool 						 `json:"sanitize,omitempty" yaml:"sanitize,omitempty"`
	Selectors 			GatewayExt_Restful_Selectors `json:"selectors,omitempty" yaml:"selectors,omitempty"`
}

type GatewayExt_Restful_Selectors struct {
	// Parse XML doc patterns
	Title      			string 						`json:"title,omitempty" yaml:"title,omitempty"`
	Result  			string 						`json:"result,omitempty" yaml:"result,omitempty"`
	Content 	    	string 						`json:"content,omitempty" yaml:"content,omitempty"`
	Tags 	    		[]Tag 						`json:"tags,omitempty" yaml:"tags,omitempty"`
	Extra 	    		string 						`json:"extra,omitempty" yaml:"extra,omitempty"`
	URL 	       		GatewayLink 				`json:"url,omitempty" yaml:"url,omitempty"`
}

