package handler

import (
	"customClothing/src/config"
	errors "customClothing/src/error"
	"customClothing/src/orderService"
	"customClothing/src/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

var OrderSvc orderService.OrderService

func RegisterOrderHandlers(r *gin.RouterGroup) {
	OrderSvc = orderService.MakeUserService()

	r.GET("/all", ListOrdersHandler)        //查询所有订单
	r.GET("/single", GetSingleOrderHandler) //查询单个订单
	r.PUT("/cost", UpdateCostHandler)       //修改订单费用
	r.DELETE("", CancelOrderHandler)        //取消订单
	r.POST("/confirm", ConfirmOrderHandler) //确认订单
	r.POST("/reporter", ReportOrderHandler) //举报
}

func ListOrdersHandler(c *gin.Context) {
	req := orderService.ListOrdersReq{}
	if err := c.ShouldBind(&req); err != nil {
		response.RespError(http.StatusBadRequest, c, errors.REQ_PARAMETER_ERROR, "参数错误")
	}

	if req.PageNum <= 0 || req.PageSize <= 0 {
		response.RespError(http.StatusBadRequest, c, errors.REQ_PARAMETER_ERROR, "参数错误")
	}

	//验证用户是否登录
	err := UserSvc.VerifyToken([]byte(c.Query(config.Cfg().TokenCfg.HeaderKey)))
	if err != nil {
		response.RespError(http.StatusBadRequest, c, err.Code(), err.Msg())
	}

	resp, err := OrderSvc.ListOrders(c, &req)
	response.Success(c, err, resp)
}

func GetSingleOrderHandler(c *gin.Context) {
	req := orderService.GetOrderReq{}
	if err := c.ShouldBind(&req); err != nil {
		response.RespError(http.StatusBadRequest, c, errors.REQ_PARAMETER_ERROR, "参数错误")
	}

	if len(req.OrderId) == 0 {
		response.RespError(http.StatusBadRequest, c, errors.REQ_PARAMETER_ERROR, "参数错误")
	}

	//验证用户是否登录
	err := UserSvc.VerifyToken([]byte(c.Query(config.Cfg().TokenCfg.HeaderKey)))
	if err != nil {
		response.RespError(http.StatusBadRequest, c, err.Code(), err.Msg())
	}

	resp, err := OrderSvc.GetSingleOrder(c, &req)
	response.Success(c, err, resp)
}

func UpdateCostHandler(c *gin.Context) {
	req := orderService.UpdateCostReq{}
	if err := c.ShouldBind(&req); err != nil {
		response.RespError(http.StatusBadRequest, c, errors.REQ_PARAMETER_ERROR, "参数错误")
	}

	if len(req.OrderId) == 0 {
		response.RespError(http.StatusBadRequest, c, errors.REQ_PARAMETER_ERROR, "参数错误")
	}

	if req.Cost < 0 {
		response.RespError(http.StatusBadRequest, c, errors.REQ_PARAMETER_ERROR, "参数错误")
	}

	//验证用户是否登录
	err := UserSvc.VerifyToken([]byte(c.Query(config.Cfg().TokenCfg.HeaderKey)))
	if err != nil {
		response.RespError(http.StatusBadRequest, c, err.Code(), err.Msg())
	}

	resp, err := OrderSvc.UpdateCost(c, &req)
	response.Success(c, err, resp)
}

func CancelOrderHandler(c *gin.Context) {
	req := orderService.CancelOrderReq{}
	if err := c.ShouldBind(&req); err != nil {
		response.RespError(http.StatusBadRequest, c, errors.REQ_PARAMETER_ERROR, "参数错误")
	}

	if len(req.OrderId) == 0 {
		response.RespError(http.StatusBadRequest, c, errors.REQ_PARAMETER_ERROR, "参数错误")
	}

	//验证用户是否登录
	err := UserSvc.VerifyToken([]byte(c.Query(config.Cfg().TokenCfg.HeaderKey)))
	if err != nil {
		response.RespError(http.StatusBadRequest, c, err.Code(), err.Msg())
	}

	resp, err := OrderSvc.CancelOrder(c, &req)
	response.Success(c, err, resp)
}

func ConfirmOrderHandler(c *gin.Context) {
	req := orderService.ConfirmOrderReq{}
	if err := c.ShouldBind(&req); err != nil {
		response.RespError(http.StatusBadRequest, c, errors.REQ_PARAMETER_ERROR, "参数错误")
	}

	if len(req.OrderId) == 0 {
		response.RespError(http.StatusBadRequest, c, errors.REQ_PARAMETER_ERROR, "参数错误")
	}

	//验证用户是否登录
	err := UserSvc.VerifyToken([]byte(c.Query(config.Cfg().TokenCfg.HeaderKey)))
	if err != nil {
		response.RespError(http.StatusBadRequest, c, err.Code(), err.Msg())
	}

	resp, err := OrderSvc.ConfirmOrder(c, &req)
	response.Success(c, err, resp)
}

func ReportOrderHandler(c *gin.Context) {
	req := orderService.ReportOrderReq{}
	if err := c.ShouldBind(&req); err != nil {
		response.RespError(http.StatusBadRequest, c, errors.REQ_PARAMETER_ERROR, "参数错误")
	}

	if len(req.OrderId) == 0 {
		response.RespError(http.StatusBadRequest, c, errors.REQ_PARAMETER_ERROR, "参数错误")
	}

	//验证用户是否登录
	err := UserSvc.VerifyToken([]byte(c.Query(config.Cfg().TokenCfg.HeaderKey)))
	if err != nil {
		response.RespError(http.StatusBadRequest, c, err.Code(), err.Msg())
	}

	resp, err := OrderSvc.ReportOrder(c, &req)
	response.Success(c, err, resp)
}
