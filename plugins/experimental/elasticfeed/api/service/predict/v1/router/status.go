package router

import (
	"github.com/roscopecoltran/feedify"
	"github.com/roscopecoltran/elasticfeed/service/predict/v1/controller"
)

func InitStatusRouters() {
	feedify.Router("/v1/predict/status", &controller.StatusController{}, "get:Get")
}
