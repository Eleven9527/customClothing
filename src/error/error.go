package errors

const (
	//公用错误码
	INTERNAL_ERROR      = -1  //内部错误
	REQ_PARAMETER_ERROR = 400 //参数错误

	//用户模块错误码
	USER_DISPLAYNAME_LENGTH               = 10011 //用户名长度错误
	USER_PASSWORD_LENGTH                  = 10012 //密码长度错误
	USER_PHONE_ERROR                      = 10013 //手机号格式错误
	USER_EXIST                            = 10014 //用户已存在
	AUTHCODE_ERROR                        = 10015 //验证码错误
	USER_NOT_EXIST                        = 10016 //用户不存在
	USER_UUID_ERROR                       = 10017 //uuid长度错误
	USER_NAME_LENGTH_EEEOR                = 10018 //用户真实姓名长度错误
	ROLE_NOT_EXIST                        = 10019 //角色不存在
	TOKEN_VERIFY_ERROR                    = 10020 //token验证失败
	USER_NOT_LOGIN                        = 10021 //用户未登录
	ORDER_NOT_EXIST                       = 10022 //订单不存在
	ORDER_STATUS_ERROR                    = 10023 //订单状态错误
	CONTENT_LENGTH_ERROR                  = 10024 //文本长度错误
	ROLE_ERROR                            = 10025 //用户角色错误
	ORDER_DETAIL_NOT_EXIST                = 10026 //订单详情不存在
	MARGIN_ERROR                          = 10027 //用户未缴纳保证金
	MARGIN_WITHDRAW_APPLICATION_NOT_EXIST = 10028 //保证金提现申请记录不存在
)

type Error interface {
	Code() int
	Msg() string
}

type error struct {
	code int    `json:"code"`
	msg  string `json:"msg"`
}

func New(code int, msg string) Error {
	return &error{
		code: code,
		msg:  msg,
	}
}

func (e *error) Code() int {
	return e.code
}

func (e *error) Msg() string {
	return e.msg
}
