package main

import (
	"github.com/roscopecoltran/elasticfeed/elasticfeed"
)

func main() {
	engine := elasticfeed.NewElasticfeed()
	engine.Run()
}
