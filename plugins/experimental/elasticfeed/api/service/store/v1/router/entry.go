package router

import (
	"github.com/roscopecoltran/feedify"
	"github.com/roscopecoltran/elasticfeed/service/store/v1/controller"
)

func InitEntryRouters() {
	feedify.Router("/v1/application/:applicationId:string/feed/:feedId:int/entry", &controller.EntryController{}, "get:GetList;post:Post")
	feedify.Router("/v1/application/:applicationId:string/feed/:feedId:int/entry/:feedEntryId:int", &controller.EntryController{}, "get:Get;delete:Delete;put:Put")

	// not implemented yet!
	feedify.Router("/v1/application/:applicationId:string/feed/:feedId:int/entry/:feedEntryId:int/metric", &controller.EntryController{}, "get:Get;delete:Delete;put:Put;post:Post")
}
