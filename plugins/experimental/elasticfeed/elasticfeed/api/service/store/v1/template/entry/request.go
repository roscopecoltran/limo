package entry

import (
	"github.com/roscopecoltran/feedify/contextor"
)


/**
 * @apiDefine EntryGetListByFeedRequest
 *
 * @apiParam {String}  applicationId  The application id
 * @apiParam {String}  feedId         The application id
 */
func RequestGetListByFeed(input *contextor.Input) {

}

/**
 * @apiDefine EntryGetRequest
 *
 * @apiParam {String}  applicationId  The application id
 * @apiParam {String}  entryId        The entry id
 */
/**
 * @apiDefine EntryGetByFeedRequest
 *
 * @apiParam {String}  applicationId  The application id
 * @apiParam {String}  feedId         The feed id
 * @apiParam {String}  entryId        The entry id
 */
func RequestGet(input *contextor.Input) {

}

/**
 * @apiDefine EntryPostRequest
 *
 * @apiParam {String}    applicationId  The application id
 * @apiParam {String}    data           The data of the entry
 * @apiParam {String[]}  [tagList]      Tags of the entry
 */
func RequestPost(input *contextor.Input) {

}

/**
 * @apiDefine EntryPostToFeedRequest
 *
 * @apiParam {String}    applicationId  The application id
 * @apiParam {String}    feedId         The feed id
 * @apiParam {String}    data           The data of the entry
 * @apiParam {String[]}  [tagList]      Tags of the entry
 */
/**
 * @apiDefine EntryAddToFeedRequest
 *
 * @apiParam {String}  applicationId  The application id
 * @apiParam {String}  feedId         The feed id
 * @apiParam {String}  entryId        The entry id
 */
func RequestPostToFeed(input *contextor.Input) {

}

/**
 * @apiDefine EntryPutRequest
 *
 * @apiParam {String}    applicationId  The application id
 * @apiParam {String}    entryId        The entry id
 * @apiParam {String}    data           The data of the entry
 * @apiParam {String[]}  [tagList]      Tags of the entry
 */
func RequestPut(input *contextor.Input) {

}

/**
 * @apiDefine EntryDeleteRequest
 *
 * @apiParam {String}  applicationId  The application id
 * @apiParam {String}  entryId        The entry id
 */
func RequestDelete(input *contextor.Input) {

}

/**
 * @apiDefine EntryRemoveRequest
 *
 * @apiParam {String}  applicationId  The application id
 * @apiParam {String}  feedId         The feed id
 * @apiParam {String}  entryId        The entry id
 */
func RequestRemove(input *contextor.Input) {

}
