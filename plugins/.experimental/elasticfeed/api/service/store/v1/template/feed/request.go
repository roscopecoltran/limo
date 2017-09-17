package feed

import (
	"github.com/roscopecoltran/feedify/contextor"
)


/**
 * @apiDefine FeedGetListRequest
 *
 * @apiParam {String}  applicationId  The application id
 */
func RequestGetList(input *contextor.Input) {

}

/**
 * @apiDefine FeedGetRequest
 *
 * @apiParam {String}  applicationId  The application id
 * @apiParam {String}  feedId         The feed id
 */
func RequestGet(input *contextor.Input) {

}

/**
 * @apiDefine FeedPostRequest
 *
 * @apiParam {String}    applicationId  The application id
 * @apiParam {String}    name           The name of the feed
 * @apiParam {String}    [description]  The description of the feed
 * @apiParam {String[]}  [tagList]      Tags of the feed
 */
func RequestPost(input *contextor.Input) {

}

/**
 * @apiDefine FeedPutRequest
 *
 * @apiParam {String}    applicationId  The application id
 * @apiParam {String}    feedId         The feed id
 * @apiParam {String}    name           The name of the feed
 * @apiParam {String}    [description]  The description of the feed
 * @apiParam {String[]}  [tagList]      Tags of the feed
 */
func RequestPut(input *contextor.Input) {

}

/**
 * @apiDefine FeedDeleteRequest
 *
 * @apiParam {String}  applicationId  The application id
 * @apiParam {String}  feedId         The feed id
 */
func RequestDelete(input *contextor.Input) {

}
