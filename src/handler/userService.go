package handler

import (
	"customClothing/src/config"
	errors "customClothing/src/error"
	"customClothing/src/response"
	"customClothing/src/userService"
	"customClothing/src/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

var UserSvc userService.UserService

func RegisterUserHandlers(r *gin.RouterGroup) {
	UserSvc = userService.MakeUserService()

	r.GET("/authcode", GetAuthCodeHandler) //获取验证码
	r.POST("", RegisterHandler)            //注册用户
	r.POST("/login", LoginHandler)         //登录
	r.POST("/kyc", KycHandler)             //kyc
}
func RegisterHandler(c *gin.Context) {
	req := userService.RegisterUserReq{}
	if err := c.ShouldBind(&req); err != nil {
		response.RespError(http.StatusBadRequest, c, errors.REQ_PARAMETER_ERROR, "参数错误")
	}

	if len(req.DisplayName) < 6 || len(req.DisplayName) > 20 {
		response.RespError(http.StatusBadRequest, c, errors.USER_DISPLAYNAME_LENGTH, "用户名长度不符合要求")
	}

	if len(req.Password) < 6 || len(req.Password) > 20 {
		response.RespError(http.StatusBadRequest, c, errors.USER_PASSWORD_LENGTH, "密码长度不符合要求")
	}

	if !utils.VerifyMobileFormat(req.Phone) {
		response.RespError(http.StatusBadRequest, c, errors.USER_PHONE_ERROR, "手机号格式错误")
	}

	resp, err := UserSvc.RegisterUser(c, &req)
	response.Success(c, err, resp)
}

func GetAuthCodeHandler(c *gin.Context) {
	req := userService.GetAuthCodeReq{}

	resp, err := UserSvc.GetAuthCode(c, &req)
	response.Success(c, err, resp)
}

func LoginHandler(c *gin.Context) {
	req := userService.LoginReq{}
	if err := c.ShouldBind(&req); err != nil {
		response.RespError(http.StatusBadRequest, c, errors.REQ_PARAMETER_ERROR, "参数错误")
	}

	if !utils.VerifyMobileFormat(req.Phone) {
		response.RespError(http.StatusBadRequest, c, errors.USER_PHONE_ERROR, "手机号格式错误")
	}

	if len(req.Password) < 6 || len(req.Password) > 20 {
		response.RespError(http.StatusBadRequest, c, errors.USER_PASSWORD_LENGTH, "密码长度不符合要求")
	}

	//todo:校验图片验证码
	if req.AuthCode != "12345" {
		response.RespError(http.StatusBadRequest, c, errors.AUTHCODE_ERROR, "验证码错误")
	}

	resp, err := UserSvc.Login(c, &req)
	response.Success(c, err, resp)
}

func KycHandler(c *gin.Context) {
	req := userService.KycReq{}
	if err := c.ShouldBind(&req); err != nil {
		response.RespError(http.StatusBadRequest, c, errors.REQ_PARAMETER_ERROR, "参数错误")
	}

	if len(req.UserId) == 0 {
		response.RespError(http.StatusBadRequest, c, errors.USER_UUID_ERROR, "uuid长度错误")
	}

	if !utils.VerifyMobileFormat(req.Phone) {
		response.RespError(http.StatusBadRequest, c, errors.USER_PHONE_ERROR, "手机号格式错误")
	}

	if len(req.Name) < 2 || len(req.Name) > 10 {
		response.RespError(http.StatusBadRequest, c, errors.USER_NAME_LENGTH_EEEOR, "姓名长度错误")
	}

	//验证用户是否登录
	err := UserSvc.VerifyToken([]byte(c.Query(config.Cfg().TokenCfg.HeaderKey)))
	if err != nil {
		response.RespError(http.StatusBadRequest, c, err.Code(), err.Msg())
	}

	resp, err := UserSvc.Kyc(c, &req)
	response.Success(c, err, resp)
}
