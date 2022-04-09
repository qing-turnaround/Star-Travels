package controller

type ResCode int64

// 状态码
const (
	CodeSuccess         ResCode = 1000 + iota // 成功
	CodeInvalidParam                          // 参数错误
	CodeUserExist                             // 用户已经存在
	CodeUserNotExist                          // 用户不存在
	CodeInvalidPassword                       // 用户密码错误
	CodeServerBusy                            // 服务器繁忙
	CodeInvalidToken                          // 无效的token
	CodeNeedLogin                             // 需要登录
	CodeNotAdmin                              // 需要管理员身份才能访问
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:         "success",
	CodeInvalidParam:    "请求参数错误",
	CodeUserExist:       "用户已经存在",
	CodeUserNotExist:    "用户不存在",
	CodeInvalidPassword: "用户密码错误",
	CodeServerBusy:      "服务器繁忙",
	CodeInvalidToken:    "无效的token",
	CodeNeedLogin:       "需要登录",
	CodeNotAdmin:        "需要管理员身份才能访问",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}
