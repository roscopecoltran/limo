package plugin

import (
	"errors"
	"github.com/roscopecoltran/feedify/contextor"
	"github.com/roscopecoltran/elasticfeed/service/system/v1/template"
)

func CheckRequiredParams() {
	// pluginId
}

func GetResponseDefinition(input *contextor.Input) (*template.ResponseDefinition) {
	return template.NewResponseDefinition(input)
}

/**
 * @apiDefine PluginGetListRequest
 *
 */
func RequestGetList(input *contextor.Input) (formatter *template.ResponseDefinition, err error) {
	if template.QueryParamsCount(input.Request.URL) > 4 {
		return nil, errors.New("Too many params in URI query")
	}
	return GetResponseDefinition(input), nil
}

/**
 * @apiDefine PluginGetRequest
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
 * @apiDefine PluginPostRequest
 */
func RequestPost(input *contextor.Input) (formatter *template.ResponseDefinition, err error) {
	if template.QueryParamsCount(input.Request.URL) != 0 {
		return nil, errors.New("Too many params in URI query")
	}
	return GetResponseDefinition(input), nil
}

/**
 * @apiDefine PluginPutRequest
 *
 * @apiParam {String}    pluginId        The plugin id
 */
func RequestPut(input *context.Input) (formatter *template.ResponseDefinition, err error) {
	if template.QueryParamsCount(input.Request.URL) != 1 {
		return nil, errors.New("Too many params in URI query")
	}
	return GetResponseDefinition(input), nil
}

/**
 * @apiDefine PluginDeleteRequest
 *
 * @apiParam {String}  pluginId  The plugin id
 */
func RequestDelete(input *contextor.Input) (formatter *template.ResponseDefinition, err error) {
	if template.QueryParamsCount(input.Request.URL) != 1 {
		return nil, errors.New("Too many params in URI query")
	}
	return GetResponseDefinition(input), nil
}
