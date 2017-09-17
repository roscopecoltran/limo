package admin

import (
	"github.com/roscopecoltran/feedify/contextor"
)


/**
 * @apiDefine AdminGetListRequest
 *
 */
func RequestGetList(input *contextor.Input) {

}

/**
 * @apiDefine AdminGetRequest
 *
 * @apiParam {String} adminId  The admin user id
 */
func RequestGet(input *contextor.Input) {

}

/**
 * @apiDefine AdminPostRequest
 *
 * @apiParam {String}    mail           The E-Mail Address of the admin user
 * @apiParam {Object[]}  roleList       A array of all roles
 * @apiParam {Int}       roleList.id    Role id (see full list at Appendix)
 * @apiParam {String}    roleList.name  Role name
 */
func RequestPost(input *contextor.Input) {

}

/**
 * @apiDefine AdminPutRequest
 *
 * @apiParam {String}    adminId        The admin user id
 * @apiParam {String}    mail           The E-Mail Address of the admin user
 * @apiParam {Object[]}  roleList       A array of all roles
 * @apiParam {Int}       roleList.id    Role id (see full list at Appendix)
 * @apiParam {String}    roleList.name  Role name
 */
func RequestPut(input *contextor.Input) {

}

/**
 * @apiDefine AdminDeleteRequest
 *
 * @apiParam {String}  adminId  The admin user id
 */
func RequestDelete(input *contextor.Input) {

}
