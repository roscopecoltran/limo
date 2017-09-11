package components

import (
	"github.com/jinzhu/configor"
)

type DockerConfig struct {
	Containers 			[]ContainerConfig
}

type ContainerConfig struct {
	Active 				bool 			`default:"false"`	
    ContainerName 		string 			`json:"container_name" yaml:"container_name"` 
    Image 				string 			`required:"true" json:"image" yaml:"image"`
}