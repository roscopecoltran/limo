package router

import (
	"github.com/roscopecoltran/feedify"
	"github.com/roscopecoltran/elasticfeed/service/store/v1/controller"
)

func InitAdminRouters() {
	feedify.Router("/v1/admin", &controller.AdminController{}, "get:GetList;post:Post")
	feedify.Router("/v1/admin/:adminId:string", &controller.AdminController{}, "get:Get;delete:Delete;put:Put")
}
