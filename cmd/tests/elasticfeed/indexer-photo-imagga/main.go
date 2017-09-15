package main

import (
	sensor "github.com/roscopecoltran/elasticfeed-plugins/indexer/photo-imagga"
	"github.com/roscopecoltran/elasticfeed/plugin"
)

func main() {
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}
	server.RegisterPipeline(new(sensor.Indexer))
	server.Serve()
}