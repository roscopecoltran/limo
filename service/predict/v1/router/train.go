package router

import (
	"github.com/roscopecoltran/feedify"
	"github.com/roscopecoltran/elasticfeed/service/predict/v1/controller"
)

func InitTrainRouters() {
	feedify.Router("/v1/predict/train", &controller.DefaultController{}, "get:Get")
}
