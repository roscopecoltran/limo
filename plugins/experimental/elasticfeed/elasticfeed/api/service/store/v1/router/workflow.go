package router

import (
	"github.com/roscopecoltran/feedify"
	"github.com/roscopecoltran/elasticfeed/service/store/v1/controller"
)

func InitEntryWorkflows() {
	feedify.Router("/v1/application/:applicationId:string/feed/:feedId:int/workflow", &controller.WorkflowController{}, "get:GetList;post:Post")
	feedify.Router("/v1/application/:applicationId:string/feed/:feedId:int/workflow/:feedWorkflowId:int", &controller.WorkflowController{}, "get:Get;delete:Delete;put:Put")
}
