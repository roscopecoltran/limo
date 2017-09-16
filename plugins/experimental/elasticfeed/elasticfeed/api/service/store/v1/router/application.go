package router

import (
	"github.com/roscopecoltran/feedify"
	"github.com/roscopecoltran/elasticfeed/service/store/v1/controller"
)

func InitApplicationRouters() {
	feedify.Router("/v1/application", &controller.ApplicationController{}, "get:GetList;post:Post")
	feedify.Router("/v1/application/:applicationId:string", &controller.ApplicationController{}, "get:Get;delete:Delete;put:Put")
}
