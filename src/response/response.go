package response

import (
	errors "customClothing/src/error"
	"github.com/gin-gonic/gin"
	"net/http"
)

type response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Resp200(c *gin.Context, err errors.Error, data interface{}) {
	res := response{
		Code: err.Code,
		Msg:  err.Msg,
		Data: data,
	}

	//为-1单独添加msg，handler中可以不用加，省事
	if res.Code == errors.INTERNAL_ERROR {
		res.Msg = "内部错误"
	}

	c.JSON(http.StatusOK, res)
	c.Next()
}

func Resp400(c *gin.Context, code int, msg string) {
	res := response{
		Code: code,
		Msg:  msg,
		Data: nil,
	}

	c.JSON(http.StatusBadRequest, res)
	c.Next()
}
