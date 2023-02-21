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

	r.GET("/all", ListOrdersHandler)                                   //查询所有订单
	r.GET("/single", GetSingleOrderHandler)                            //查询单个订单
	r.PUT("/cost", UpdateCostHandler)                                  //修改订单费用
	r.DELETE("", CancelOrderHandler)                                   //取消订单
	r.POST("/confirm", ConfirmOrderHandler)                            //确认订单
	r.POST("/reporter", ReportOrderHandler)                            //举报
	r.POST("/designArtwork", UploadDesignArtworkHandler)               //上传设计图稿
	r.POST("/patternArtwork", UploadPatternArtworkHandler)             //上传版型图稿
	r.POST("/patternMakingProcess", UploadPatternMakingProcessHandler) //上传版型制作工艺
	r.POST("/sampleImage", UploadSampleImageHandler)                   //上传样品成衣图
	r.POST("/showVideo", UploadShowVideoHandler)                       //上传模特展示视频
}

func ListOrdersHandler(c *gin.Context) {
	roleCode := c.GetInt("roleCode")
	status := c.GetInt("status")
	pageNum := c.GetInt("pageNum")
	pageSize := c.GetInt("pageSize")
	req := orderService.ListOrdersReq{
		UserId:   c.Query("userId"),
		RoleCode: roleCode,
		Status:   status,
		PageNum:  uint(pageNum),
		PageSize: uint(pageSize),
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
	req := orderService.GetOrderReq{
		OrderId: c.Query("orderId"),
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

func UploadDesignArtworkHandler(c *gin.Context) {
	req := orderService.UploadDesignArtworkReq{}
	if err := c.ShouldBind(&req); err != nil {
		response.RespError(http.StatusBadRequest, c, errors.REQ_PARAMETER_ERROR, "参数错误")
	}

	//todo:最多6张图片

	//验证用户是否登录
	err := UserSvc.VerifyToken([]byte(c.Query(config.Cfg().TokenCfg.HeaderKey)))
	if err != nil {
		response.RespError(http.StatusBadRequest, c, err.Code(), err.Msg())
	}

	resp, err := OrderSvc.UploadDesignArtwork(c, &req)
	response.Success(c, err, resp)
}

func UploadPatternArtworkHandler(c *gin.Context) {
	req := orderService.UploadPatternArtworkReq{}
	if err := c.ShouldBind(&req); err != nil {
		response.RespError(http.StatusBadRequest, c, errors.REQ_PARAMETER_ERROR, "参数错误")
	}

	//todo:最多6张图片

	//验证用户是否登录
	err := UserSvc.VerifyToken([]byte(c.Query(config.Cfg().TokenCfg.HeaderKey)))
	if err != nil {
		response.RespError(http.StatusBadRequest, c, err.Code(), err.Msg())
	}

	resp, err := OrderSvc.UploadPatternArtwork(c, &req)
	response.Success(c, err, resp)
}

func UploadPatternMakingProcessHandler(c *gin.Context) {
	req := orderService.UploadPatternMakingProcessReq{}
	if err := c.ShouldBind(&req); err != nil {
		response.RespError(http.StatusBadRequest, c, errors.REQ_PARAMETER_ERROR, "参数错误")
	}

	//限制内容长度
	if len(req.Content) > 500 {
		response.RespError(http.StatusBadRequest, c, errors.CONTENT_LENGTH_ERROR, "内容长度过大，请限制在500以内")
	}

	//验证用户是否登录
	err := UserSvc.VerifyToken([]byte(c.Query(config.Cfg().TokenCfg.HeaderKey)))
	if err != nil {
		response.RespError(http.StatusBadRequest, c, err.Code(), err.Msg())
	}

	resp, err := OrderSvc.UploadPatternMakingProcess(c, &req)
	response.Success(c, err, resp)
}

func UploadSampleImageHandler(c *gin.Context) {
	req := orderService.UploadSampleImageReq{}
	if err := c.ShouldBind(&req); err != nil {
		response.RespError(http.StatusBadRequest, c, errors.REQ_PARAMETER_ERROR, "参数错误")
	}

	//todo:最多6张图片

	//验证用户是否登录
	err := UserSvc.VerifyToken([]byte(c.Query(config.Cfg().TokenCfg.HeaderKey)))
	if err != nil {
		response.RespError(http.StatusBadRequest, c, err.Code(), err.Msg())
	}

	resp, err := OrderSvc.UploadSampleImage(c, &req)
	response.Success(c, err, resp)
}

func UploadShowVideoHandler(c *gin.Context) {
	req := orderService.UploadShowVideoReq{}
	if err := c.ShouldBind(&req); err != nil {
		response.RespError(http.StatusBadRequest, c, errors.REQ_PARAMETER_ERROR, "参数错误")
	}

	//todo:每个视频最大20m

	//验证用户是否登录
	err := UserSvc.VerifyToken([]byte(c.Query(config.Cfg().TokenCfg.HeaderKey)))
	if err != nil {
		response.RespError(http.StatusBadRequest, c, err.Code(), err.Msg())
	}

	resp, err := OrderSvc.UploadShowVideo(c, &req)
	response.Success(c, err, resp)
}
