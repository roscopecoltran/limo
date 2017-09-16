package router

import (
	"github.com/roscopecoltran/feedify"
	"github.com/roscopecoltran/elasticfeed/service/predict/v1/controller"
)

func InitPredictRouters() {
	feedify.Router("/v1/predict/predict", &controller.DefaultController{}, "get:Get")
}
