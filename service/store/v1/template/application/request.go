package application

import (
	"github.com/roscopecoltran/feedify/contextor"
)

/**
 * @apiDefine ApplicationGetListRequest
 *
 */
func RequestGetList(input *contextor.Input) {

}

/**
 * @apiDefine ApplicationGetRequest
 *
 * @apiParam {String}  applicationId  The application id
 */
func RequestGet(input *contextor.Input) {

}

/**
 * @apiDefine ApplicationPostRequest
 *
 * @apiParam {String}  name           The name of the application
 * @apiParam {String}  [description]  The description of the application
 *
 * @apiExample {json} Example post body (json):
 *     {
 *       "name": "DragonBall",
 *       "description": "The DragonBall application. This will contain feeds for all DragonBall characters."
 *     }
 */
func RequestPost(input *contextor.Input) {

}

/**
 * @apiDefine ApplicationPutRequest
 *
 * @apiParam {String}  applicationId  The application id
 * @apiParam {String}  name           The name of the application
 * @apiParam {String}  [description]  The description of the application
 *
 * @apiExample {json} Example post body (json):
 *     {
 *       "name": "DragonBall",
 *       "description": "The DragonBall application. This will contain feeds for all DragonBall characters."
 *     }
 */
func RequestPut(input *contextor.Input) {

}

/**
 * @apiDefine ApplicationDeleteRequest
 *
 * @apiParam {String}  applicationId  The application id
 */
func RequestDelete(input *contextor.Input) {

}
