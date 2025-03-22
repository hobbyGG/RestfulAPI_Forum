package code

type Code int16

const (
	Success Code = iota
	InvalidParam
	EmialExisted
	NeedAuth
	AuthType
	InvalidToken
	ServeBusy
	NotLogin
)

var codeMsgMap = map[Code]string{
	Success:      "success",
	InvalidParam: "参数错误",
	EmialExisted: "邮箱已被占用",
	NeedAuth:     "需要认证",
	AuthType:     "auth认证类型错误",
	ServeBusy:    "服务繁忙",
	NotLogin:     "用户未登录",
}

func (c Code) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		return codeMsgMap[ServeBusy]
	}
	return msg
}
