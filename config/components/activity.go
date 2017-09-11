package components

import (
	"github.com/jinzhu/configor"
)

type ElkStackConfig struct {
	Elk ElkConfig
}

type ElkConfig struct {
	Active 				bool 			`default:"false"`
	ElasticSearch 		ElasticSearchConfig	
	Kibana 				KibanaConfig
	FileBeat 			FileBeatConfig
	LogStash 			LogStashConfig
}

type ElasticSearchConfig struct {
	Active 				bool 			`default:"false"`	
}

type KibanaConfig struct {
	Active 				bool 			`default:"false"`	
	ProxiedBy 			bool 			`default:"nginx"`	
	ReverseProxy 		struct { 
		Configuration 	struct { 
			Nginx 		NginxConfig
			Caddy 		CaddyConfig
			Apache2 	Apache2Config
		}
	}
}

type LogStashConfig struct {
	Active 				bool 			`default:"false"`	
}

type FileBeatConfig struct {
	Active 				bool 			`default:"false"`	
}

