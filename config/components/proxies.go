package components

import (
	"github.com/jinzhu/configor"
)

type CaddyConfig struct {
	Active 				bool 			`default:"false"`	
}

type Apache2Config struct {
	Active 				bool 			`default:"false"`	
}

type NginxConfig struct {
	Active 				bool 			`default:"false"`	
}

