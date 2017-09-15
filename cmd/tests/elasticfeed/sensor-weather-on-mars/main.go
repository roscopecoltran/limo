package main

import (
  "github.com/roscopecoltran/elasticfeed/plugin"
  sensor "github.com/roscopecoltran/elasticfeed-plugins/sensor/weather"
)

func main() {
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}
	server.RegisterSensor(new(sensor.Sensor))
	server.Serve()
}