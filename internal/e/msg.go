package e

var MsgFlags = map[int]string{
	SUCCESS:        "success",
	ERROR:          "fail",
	INVALID_PARAMS: "请求参数错误",
	UNAUTHORIZED:   "没有接口访问权限, 请登录",

	ERROR_EXIST_USER:          "用户已存在",
	ERROR_INVALID_CREDENTIALS: "用户名或密码错误",

	ERROR_AUTH:                     "Token错误",
	ERROR_AUTH_CHECK_TOKEN_FAIL:    "Token鉴权失败",
	ERROR_AUTH_CHECK_TOKEN_EXPIRED: "Token已过期",

	ERROR_NOT_EXIST_COMMUNITY: "社区不存在",

	ERROR_NOT_EXIST_POST: "帖子不存在",

	ERROR_VOTE_TIME_EXPIRED: "超过投票时间，无法投票",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
