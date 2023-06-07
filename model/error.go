package model

const (
	IDENTITY_EXIST = "user.identity.exist"

	MOBILE_INVALID             = "mobile.invalid"
	EMAIL_INVALID              = "email.invalid"
	VERIFYCODE_INVALID         = "verifycode.invalid"
	SIGNUP_MOBILE_CLAIMED      = "signup.mobile.claimed"
	SIGNUP_EMAIL_CLAIMED       = "signup.email.claimed"
	MOBILE_VERIFY_INVALIDCODE  = "mobile.verify.invalidcode"
	SIGNUP_MOBILE_INCONSISTENT = "signup.mobile.inconsistent"
	SIGNUP_INIT_FAILED         = "signup.init.failed"

	ID_INVALID = "common.id.invalid"

	LOGIN_IDENTITY_NOTEXIST = "login.identity.notexist"
	LOGIN_IDENTITY_INVALID  = "login.identity.invalid"
	LOGIN_IDENTITY_CONFILCT = "login.identity.conflict"

	LOGIN_FAILED = "login.failed"

	MOBILE_NOTEXIST           = "mobile.notexist"
	AUTH_PASSWORD_INVALID     = "auth.password.invalid"
	LOGIN_OAUTH_TYPE_NOTEXIST = "login.oauthtype.notexist"

	USER_NOTEXIST                   = "user.notexist"
	USER_FREEZE                     = "user.freeze"
	USER_DELETED                    = "user.deleted"
	USER_FIND_ERROR                 = "user.find.error"
	USER_FEEDBACK_CONTENT_NIL       = "user.feedback.content.nil"
	USER_FEEDBACK_SAVE_FAILED       = "user.feedback.save.failed"
	USER_DEVICE_TOKEN_SAVE_FAILED   = "user.devicetoken.save.failed"
	USER_SOURCE_TYPE_INVALID        = "user.source.type.invalid"
	USER_DEVICE_TOKEN_UPDATE_FAILED = "user.devicetoken.update.failed"
	USER_DEVICE_TOKEN_PARAM_EMPTY   = "user.devicetoken.param.empty"
	USER_IMAGE_UPLOAD_EMPTY         = "user.image.upload.empty"
	USER_IMAGE_SAVE_FAILED          = "user.image.save.failed"
	USER_IMAGE_TYPE_INVALID         = "user.image.type.invalid"
	USER_USERNAME_IS_EXIST          = "user.username.isexist"
	USER_USERNAME_INVALID           = "user.username.invalid"
	USER_NICKNAME_INVALID           = "user.nickname.invalid"

	OAUTH_REQUEST_ACCESS_TOKEN_FAILED = "oauth.request.accesstoken.failed"

	SNS_FOLLOW_FAILED        = "sns.follow.failed"
	SNS_QUERY_FANS_FAILED    = "sns.query.fans.failed"
	SNS_QUERY_FOLLOWS_FAILED = "sns.query.follows.failed"
	SNS_QUERY_FRIENDS_FAILED = "sns.query.friends.failed"
	SNS_SEND_PMSG_FAILED     = "sns.send.pmsg.failed"
	SNS_QUERY_MSG_LUT_EMPTY  = "sns.query.msg.lut.empty"
	SNS_QUERY_MSG_FAILED     = "sns.query.msg.failed"

	SERVER_INTERNAL_ERROR = "s.500"
	BAD_REQUEST           = "400"

	UNAUTHORIZED = "unauthorized"

	/**
		operate database error
	**/

)
