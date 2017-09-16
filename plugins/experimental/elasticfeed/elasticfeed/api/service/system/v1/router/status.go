package router

import (
	"github.com/roscopecoltran/feedify"
	"github.com/roscopecoltran/elasticfeed/service/system/v1/controller"
)

func InitRouters() {
	feedify.Router("/v1/system/status", &controller.StatusController{}, "get:Get")
}
