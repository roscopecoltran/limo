package main

import (
	"github.com/roscopecoltran/elasticfeed-plugins/sensor/weather"
	"github.com/roscopecoltran/elasticfeed/plugin"
)

func main() {
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}
	server.RegisterPipeline(new(weather.Sensor))
	server.Serve()
}
