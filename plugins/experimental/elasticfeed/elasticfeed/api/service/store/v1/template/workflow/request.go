package workflow

import (
	"errors"
	"github.com/roscopecoltran/feedify/contextor"
	"github.com/roscopecoltran/elasticfeed/service/store/v1/template"
)

func CheckRequiredParams() {
	// workflowId
}

func GetResponseDefinition(input *contextor.Input) (*template.ResponseDefinition) {
	return template.NewResponseDefinition(input)
}

/**
 * @apiDefine WorkflowGetListRequest
 *
 */
func RequestGetList(input *contextor.Input) (formatter *template.ResponseDefinition, err error) {
	if template.QueryParamsCount(input.Request.URL) > 4 {
		return nil, errors.New("Too many params in URI query")
	}
	return GetResponseDefinition(input), nil
}

/**
 * @apiDefine WorkflowGetRequest
 *
 * @apiParam {String} pluginId  The plugin id
 */
func RequestGet(input *contextor.Input) (formatter *template.ResponseDefinition, err error) {
	if template.QueryParamsCount(input.Request.URL) != 1 {
		return nil, errors.New("Too many params in URI query")
	}
	return GetResponseDefinition(input), nil
}

/**
 * @apiDefine WorkflowPostRequest
 */
func RequestPost(input *contextor.Input) (formatter *template.ResponseDefinition, err error) {
	if template.QueryParamsCount(input.Request.URL) != 0 {
		return nil, errors.New("Too many params in URI query")
	}
	return GetResponseDefinition(input), nil
}

/**
 * @apiDefine WorkflowPutRequest
 *
 * @apiParam {String}    pluginId        The plugin id
 */
func RequestPut(input *contextor.Input) (formatter *template.ResponseDefinition, err error) {
	if template.QueryParamsCount(input.Request.URL) > 4 {
		return nil, errors.New("Too many params in URI query")
	}
	return GetResponseDefinition(input), nil
}

/**
 * @apiDefine WorkflowDeleteRequest
 *
 * @apiParam {String}  pluginId  The plugin id
 */
func RequestDelete(input *contextor.Input) (formatter *template.ResponseDefinition, err error) {
	if template.QueryParamsCount(input.Request.URL) != 1 {
		return nil, errors.New("Too many params in URI query")
	}
	return GetResponseDefinition(input), nil
}
