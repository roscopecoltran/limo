package model

import (
	//"fmt"
	//"strings"
	//"time"
	"github.com/jinzhu/gorm"
	//"github.com/sirupsen/logrus"
)

type GatewayMeta struct {
	// gorm.Model
	Disable 			bool 									`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	General 			GatewayGlobal 							`json:"general,omitempty" yaml:"general,omitempty"`
	Locales 			GatewayGlobal_Locales 					`json:"locales,omitempty" yaml:"locales,omitempty"`
	Outgoing 			GatewayGlobal_Outgoing 					`json:"outgoing,omitempty" yaml:"outgoing,omitempty"`
	Server 				Global_ServerOptions 					`json:"server,omitempty" yaml:"server,omitempty"`
	Frontend 			Web_UIOptions 							`json:"frontend,omitempty" yaml:"frontend,omitempty"`
	Backend 			Web_UIOptions 							`json:"backend,omitempty" yaml:"backend,omitempty"`
	Engines 			[]GatewayMeta_Engine 					`json:"engines,omitempty" yaml:"engines,omitempty"`
}

type GatewayMeta_EngineExtractors struct {
	// gorm.Model
	Disable 			bool 									`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	Tasks 				map[string]Task_Options 				`json:"tasks,omitempty" yaml:"tasks,omitempty"`
	Endpoints struct {
		Xpath 			map[string]GatewayExt_XPath_Endpoint 	`json:"xpath,omitempty" yaml:"xpath,omitempty"`
		Restful 		map[string]GatewayExt_Restful_Endpoint 	`json:"restful,omitempty" yaml:"restful,omitempty"`
		//GraphQL 		map[string]GraphQLOutput 				`json:"graphql,omitempty" yaml:"graphql,omitempty"`
	} `json:"endpoints,omitempty" yaml:"endpoints,omitempty"`
}

type GatewayMeta_Engine struct {
	// gorm.Model
	Disable 			bool 									`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	Name            	string 									`json:"name,omitempty" yaml:"name,omitempty"`
	Disabled        	bool   									`default:"false" json:"disabled,omitempty" yaml:"disabled,omitempty"`
	Profile 			Web_SiteProfile 						`json:"profile,omitempty" yaml:"profile,omitempty"`
	Content 			GatewayMeta_EngineExtractors 			`json:"content,omitempty" yaml:"content,omitempty"`
	Options 			Global_SearchOptions 					`json:"options,omitempty" yaml:"options,omitempty"`
	Ranks 				GatewayMeta_EngineRankOptions 			`json:"ranks,omitempty" yaml:"ranks,omitempty"`
	Auth 				[]GatewayAuth_Credentials 				`json:"auth,omitempty" yaml:"auth,omitempty"`
	// Query   // xpath, api (krakend ?!)
	// Results // xpath, api (krakend ?!)
	// Content // xpath, api (krakend ?!)
	// Engines // xpath, api (krakend ?!)
	// Options // xpath, api (krakend ?!)
}

type GatewayMeta_SearxEngine struct {
	BaseURL         	string 									`json:"base_url" yaml:"base_url"`
	Categories      	string 									`json:"categories" yaml:"categories"`
	ContentQuery    	string 									`json:"content_query" yaml:"content_query"`
	ContentXpath    	string 									`json:"content_xpath" yaml:"content_xpath"`
	Disabled        	bool   									`default:"false" json:"disabled" yaml:"disabled"`
	Engine          	string 									`json:"engine" yaml:"engine"`
	FirstPageNum    	int  									`default:"1" json:"first_page_num" yaml:"first_page_num"`
	Name            	string 									`json:"name" yaml:"name"`
	NumberOfResults 	int  									`default:"100" json:"number_of_results" yaml:"number_of_results"`
	PageSize        	int  									`default:"100" json:"page_size" yaml:"page_size"`
	Paging          	bool   									`default:"100" json:"paging" yaml:"paging"`
	ResultsQuery    	string 									`json:"results_query" yaml:"results_query"`
	ResultsXpath    	string 									`json:"results_xpath" yaml:"results_xpath"`
	SearchType      	string 									`json:"search_type" yaml:"search_type"`
	SearchURL       	string 									`json:"search_url" yaml:"search_url"`
	Shortcut        	string 									`json:"shortcut" yaml:"shortcut"`
	SuggestionXpath 	string 									`json:"suggestion_xpath" yaml:"suggestion_xpath"`
	Timeout         	float64  								`default:"2.0" json:"timeout" yaml:"timeout"`
	TitleQuery      	string 									`json:"title_query" yaml:"title_query"`
	TitleXpath      	string 									`json:"title_xpath" yaml:"title_xpath"`
	URL             	string 									`json:"url" yaml:"url"`
	URLQuery        	string 									`json:"url_query" yaml:"url_query"`
	URLXpath        	string 									`json:"url_xpath" yaml:"url_xpath"`
	Weight          	int  									`json:"weight" yaml:"weight"`
}

type GatewayMeta_EngineRankOptions struct {
	// gorm.Model
	Disable 			bool 									`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	Weight      		int  									`default:"1" json:"weight,omitempty" yaml:"weight,omitempty"`			
	Neurons 			[]GatewayNeurons 						`json:"nodes,omitempty" yaml:"nodes,omitempty"`
}

type GatewayMeta_EnginesGroups struct {
	gorm.Model
	Disable 			bool 									`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	Name 				string 									`json:"name,omitempty" yaml:"name,omitempty"`
}
