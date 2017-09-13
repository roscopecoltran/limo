package model

import (
	//"fmt"
	//"strings"
	//"time"
	"github.com/jinzhu/gorm"
	//"github.com/sirupsen/logrus"
)

/*
	Gateway - global definitions
*/
type GatewayGlobal struct {
	gorm.Model
	Shortcut        string 									`json:"shortcut,omitempty" yaml:"shortcut,omitempty"`
	Timeout         float64  								`default:"3.0" json:"timeout,omitempty" yaml:"timeout,omitempty"`		
	Autocomplete 	bool 									`default:"true" json:"autocomplete,omitempty" yaml:"autocomplete,omitempty"`
	SafeSearch   	bool  									`default:"true" json:"safe_search,omitempty" yaml:"safe_search,omitempty"`
}

type GatewayGlobal_PagingOptions struct {
	Disable 			bool 								`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	PageFollow      	bool   								`default:"true" json:"paging,omitempty" yaml:"paging,omitempty"`
	PageSize        	int  								`default:"25" json:"page_size,omitempty" yaml:"page_size,omitempty"`
	PerPage 			int  								`default:"100" json:"number_of_results,omitempty" yaml:"number_of_results,omitempty"`
	Offset    			int  								`default:"1" json:"first_page_num,omitempty" yaml:"first_page_num,omitempty"`
} 

type GatewayGlobal_SearchOptions struct {
	
}

type GatewayGlobal_Request struct {
	// gorm.Model
	Disable 			bool 								`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	Base        		GatewayGlobal_StandardConnection 	`default:"base" json:"protocol,omitempty" yaml:"base,omitempty"` // http, https, tcp, udp, icp
	Endpoint        	string 								`json:"endpoint,omitempty" yaml:"endpoint,omitempty"`
	Headers 			map[string]string   				`json:"headers,omitempty" yaml:"headers,omitempty"`
	Options 			Global_SearchOptions 				`json:"options,omitempty" yaml:"options,omitempty"`
} 

/*
	Gateway - protocols requests definition
*/
type GatewayGlobal_StandardConnection struct {
	Disable 			bool 								`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	ListenAddr 			string  							`json:"listen_addr,omitempty" yaml:"listen_addr,omitempty"`
	Domain 				string  							`json:"domain,omitempty" yaml:"domain,omitempty"`
	Port 				string  							`json:"port,omitempty" yaml:"port,omitempty"`
	RetryMax 			uint   								`default:"3" json:"retry_max,omitempty" yaml:"retry_max,omitempty"`
	Timeout 			float64   							`default:"3.0" json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Secured 			bool   								`default:"true" json:"secured,omitempty" yaml:"secured,omitempty"`
	VerifySSL 			bool   								`default:"true" json:"verify_ssl,omitempty" yaml:"verify_ssl,omitempty"`
}

type GatewayGlobal_Protocols struct {
	Disable 			bool 								`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	HTTP 				GatewayGlobal_StandardConnection 	`json:"http,omitempty" yaml:"http,omitempty"`
	HTTPS 				GatewayGlobal_StandardConnection 	`json:"https,omitempty" yaml:"https,omitempty"`
	WS 					GatewayGlobal_StandardConnection 	`json:"ws,omitempty" yaml:"ws,omitempty"`
	WSS 				GatewayGlobal_StandardConnection 	`json:"wss,omitempty" yaml:"wss,omitempty"`
	TCP 				GatewayGlobal_StandardConnection 	`json:"tcp,omitempty" yaml:"tcp,omitempty"`
	UDP 				GatewayGlobal_StandardConnection 	`json:"udp,omitempty" yaml:"udp,omitempty"`
	TLS 				GatewayGlobal_StandardConnection 	`json:"tls,omitempty" yaml:"tls,omitempty"`
	ICP 				GatewayGlobal_StandardConnection 	`json:"icp,omitempty" yaml:"icp,omitempty"`
	RPC 				GatewayGlobal_StandardConnection 	`json:"rpc,omitempty" yaml:"rpc,omitempty"`
}

/*
	Gateway - outgoing definition
*/
type GatewayGlobal_Outgoing struct {
	PoolConnections 	int  								`json:"pool_connections,omitempty" yaml:"pool_connections,omitempty"`
	PoolMaxsize     	int  								`json:"pool_maxsize,omitempty" yaml:"pool_maxsize,omitempty"`
	RequestTimeout  	int  								`json:"request_timeout,omitempty" yaml:"request_timeout,omitempty"`
	UseragentSuffix 	string 								`json:"useragent_suffix,omitempty" yaml:"useragent_suffix,omitempty"`
}

/*
	Gateway - locales definition
*/
type GatewayGlobal_Locales struct {
	Disable 			bool 								`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	Current 			string 								`default:"en" json:"current,omitempty" yaml:"current,omitempty"`
	Defaults 			GatewayGlobal_DefaultLocales 		`json:"defaults,omitempty" yaml:"defaults,omitempty"`
}

type GatewayGlobal_DefaultLocales struct {
	//Bg   string `json:"bg,omitempty"`
	//Cs   string `json:"cs,omitempty"`
	//De   string `json:"de,omitempty"`
	//DeDE string `json:"de_DE,omitempty"`
	//ElGR string `json:"el_GR,omitempty"`
	En   string `json:"en,omitempty" yaml:"en,omitempty"`
	//Eo   string `json:"eo,omitempty"`
	//Es   string `json:"es,omitempty"`
	//Fi   string `json:"fi,omitempty"`
	Fr   string `json:"fr,omitempty" yaml:"fr,omitempty"`
	//He   string `json:"he,omitempty"`
	//Hu   string `json:"hu,omitempty"`
	//It   string `json:"it,omitempty"`
	//Ja   string `json:"ja,omitempty"`
	//Nl   string `json:"nl,omitempty"`
	//Pt   string `json:"pt,omitempty"`
	//PtBR string `json:"pt_BR,omitempty"`
	//Ro   string `json:"ro,omitempty"`
	//Ru   string `json:"ru,omitempty"`
	//Sk   string `json:"sk,omitempty"`
	//Sv   string `json:"sv,omitempty"`
	//Tr   string `json:"tr,omitempty"`
	Uk   string `json:"uk,omitempty" yaml:"uk,omitempty"`
	//Zh   string `json:"zh,omitempty"`
}

