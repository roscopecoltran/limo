package org

import (
	"errors"
	"github.com/roscopecoltran/feedify/contextor"
	"github.com/roscopecoltran/elasticfeed/service/store/v1/template"
)

func CheckRequiredParams() {
	// orgId
}

func GetResponseDefinition(input *contextor.Input) (*template.ResponseDefinition) {
	return template.NewResponseDefinition(input)
}

/**
 * @apiDefine OrgGetListRequest
 *
 */
func RequestGetList(input *contextor.Input) (formatter *template.ResponseDefinition, err error) {
	if template.QueryParamsCount(input.Request.URL) > 4 {
		return nil, errors.New("Too many params in URI query")
	}
	return GetResponseDefinition(input), nil
}

/**
 * @apiDefine OrgGetRequest
 *
 * @apiParam {String} orgId  The org id
 */
func RequestGet(input *contextor.Input) (formatter *template.ResponseDefinition, err error) {
	if template.QueryParamsCount(input.Request.URL) != 1 {
		return nil, errors.New("Too many params in URI query")
	}
	return GetResponseDefinition(input), nil
}

/**
 * @apiDefine OrgPostRequest
 */
func RequestPost(input *contextor.Input) (formatter *template.ResponseDefinition, err error) {
	if template.QueryParamsCount(input.Request.URL) != 0 {
		return nil, errors.New("Too many params in URI query")
	}
	return GetResponseDefinition(input), nil
}

/**
 * @apiDefine OrgPutRequest
 *
 * @apiParam {String}    orgId        The org id
 */
func RequestPut(input *contextor.Input) (formatter *template.ResponseDefinition, err error) {
	if template.QueryParamsCount(input.Request.URL) != 1 {
		return nil, errors.New("Too many params in URI query")
	}
	return GetResponseDefinition(input), nil
}

/**
 * @apiDefine OrgDeleteRequest
 *
 * @apiParam {String}  orgId  The org id
 */
func RequestDelete(input *contextor.Input) (formatter *template.ResponseDefinition, err error) {
	if template.QueryParamsCount(input.Request.URL) != 1 {
		return nil, errors.New("Too many params in URI query")
	}
	return GetResponseDefinition(input), nil
}
