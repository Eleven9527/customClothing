package orderService

import (
	"context"
	errors "customClothing/src/error"
)

type OrderService interface {
	ListOrders(ctx context.Context, req *ListOrdersReq) (*ListOrdersResp, errors.Error)
	GetSingleOrder(ctx context.Context, req *GetOrderReq) (*GetOrderResp, errors.Error)
	UpdateCost(ctx context.Context, req *UpdateCostReq) (*UpdateCostResp, errors.Error)
	CancelOrder(ctx context.Context, req *CancelOrderReq) (*CancelOrderResp, errors.Error)
	ConfirmOrder(ctx context.Context, req *ConfirmOrderReq) (*ConfirmOrderResp, errors.Error)
	ReportOrder(ctx context.Context, req *ReportOrderReq) (*ReportOrderResp, errors.Error)
	UploadDesignArtwork(ctx context.Context, req *UploadDesignArtworkReq) (*UploadDesignArtworkResp, errors.Error)
	UploadPatternArtwork(ctx context.Context, req *UploadPatternArtworkReq) (*UploadPatternArtworkResp, errors.Error)
	UploadPatternMakingProcess(ctx context.Context, req *UploadPatternMakingProcessReq) (*UploadPatternMakingProcessResp, errors.Error)
	UploadSampleImage(ctx context.Context, req *UploadSampleImageReq) (*UploadSampleImageResp, errors.Error)
	UploadShowVideo(ctx context.Context, req *UploadShowVideoReq) (*UploadShowVideoResp, errors.Error)
	GetOrderSum(ctx context.Context, req *GetOrderSumReq) (*GetOrderSumResp, errors.Error)
	PublishOrder(ctx context.Context, req *PublishOrderReq) (*PublishOrderResp, errors.Error)
}

type OrderSvc struct {
	orderRepo RepoService
	cache     CacheService
	//userRepo  userService.RepoService
}

func MakeUserService() OrderService {
	return &OrderSvc{
		orderRepo: MakeRepoService(),
		cache:     MakeCacheService(),
		//userRepo:  userService.MakeRepoService(),
	}
}

//	@Summary		获取所有需求
//	@Description	查询自己发布的or接手的需求订单
//	@Tags			order模块
//	@Accept			json
//	@Produce		json
//	@Param			userId			query		number	true	"用户uuid"
//	@Param			roleCode		query		number	true	"用户角色code: 0 = 管理员，1 = 普通用户，2 = 设计师，3 = 版型师"
//	@Param			status			query		number	true	"订单状态: 1 = 待接单, 2 = 已过期, 3 = 进行中, 4 = 已完成, 5 = 已取消"
//	@Param			pageNum			query		number	true	"页数"
//	@Param			pageSize		query		number	true	"每页条数"
//	@Param			Authorization	header		string	true	"token"
//	@Success		200				{object}	ListOrdersResp
//	@Failure		400				{object}	response.response
//	@Failure		404				{object}	response.response
//	@Failure		500				{object}	response.response
//	@Router			/order/all [get]
func (o *OrderSvc) ListOrders(ctx context.Context, req *ListOrdersReq) (*ListOrdersResp, errors.Error) {
	return o.orderRepo.ListOrders(ctx, req)
}

//	@Summary		获取单个需求
//	@Description	查询自己发布的or接手的需求订单
//	@Tags			order模块
//	@Accept			json
//	@Produce		json
//	@Param			orderId			query		string	true	"订单uuid"
//	@Param			Authorization	header		string	true	"token"
//	@Success		200				{object}	GetOrderResp
//	@Failure		400				{object}	response.response
//	@Failure		404				{object}	response.response
//	@Failure		500				{object}	response.response
//	@Router			/order/single [get]
func (o *OrderSvc) GetSingleOrder(ctx context.Context, req *GetOrderReq) (*GetOrderResp, errors.Error) {
	return o.orderRepo.GetSingleOrder(ctx, req)
}

//	@Summary		修改订单费用
//	@Description	修改订单费用
//	@Tags			order模块
//	@Accept			json
//	@Produce		json
//	@Param			request			body		UpdateCostReq	true	"请求"
//	@Param			Authorization	header		string			true	"token"
//	@Success		200				{object}	UpdateCostResp
//	@Failure		400				{object}	response.response
//	@Failure		404				{object}	response.response
//	@Failure		500				{object}	response.response
//	@Router			/order/cost [post]
func (o *OrderSvc) UpdateCost(ctx context.Context, req *UpdateCostReq) (*UpdateCostResp, errors.Error) {
	//只有待接单中的订单可以修改费用
	order, err := o.orderRepo.GetSingleOrder(ctx, &GetOrderReq{OrderId: req.OrderId})
	if err != nil {
		return nil, err
	}

	if order.Order.Status != STATUS_PENDING {
		return nil, errors.New(errors.ORDER_STATUS_ERROR, "只有待接单中的订单可以修改费用")
	}

	return o.orderRepo.UpdateCost(ctx, req)
}

//	@Summary		取消订单
//	@Description	甲方取消订单
//	@Tags			order模块
//	@Accept			json
//	@Produce		json
//	@Param			request			body		CancelOrderReq	true	"请求"
//	@Param			Authorization	header		string			true	"token"
//	@Success		200				{object}	CancelOrderResp
//	@Failure		400				{object}	response.response
//	@Failure		404				{object}	response.response
//	@Failure		500				{object}	response.response
//	@Router			/order [delete]
func (o *OrderSvc) CancelOrder(ctx context.Context, req *CancelOrderReq) (*CancelOrderResp, errors.Error) {
	//只有待接单中的订单可以取消
	order, err := o.orderRepo.GetSingleOrder(ctx, &GetOrderReq{OrderId: req.OrderId})
	if err != nil {
		return nil, err
	}

	if order.Order.Status != STATUS_PENDING {
		return nil, errors.New(errors.ORDER_STATUS_ERROR, "只有待接单中的订单可以取消")
	}

	return o.orderRepo.CancelOrder(ctx, req)
}

//	@Summary		确认订单
//	@Description	甲方确认订单
//	@Tags			order模块
//	@Accept			json
//	@Produce		json
//	@Param			request			body		ConfirmOrderReq	true	"请求"
//	@Param			Authorization	header		string			true	"token"
//	@Success		200				{object}	ConfirmOrderResp
//	@Failure		400				{object}	response.response
//	@Failure		404				{object}	response.response
//	@Failure		500				{object}	response.response
//	@Router			/confirm [post]
func (o *OrderSvc) ConfirmOrder(ctx context.Context, req *ConfirmOrderReq) (*ConfirmOrderResp, errors.Error) {
	//只有进行中中的订单可以确认
	order, err := o.orderRepo.GetSingleOrder(ctx, &GetOrderReq{OrderId: req.OrderId})
	if err != nil {
		return nil, err
	}

	if order.Order.Status != STATUS_PROCESS {
		return nil, errors.New(errors.ORDER_STATUS_ERROR, "只有进行中中的订单可以确认")
	}

	return o.orderRepo.ConfirmOrder(ctx, req)
}

//	@Summary		举报
//	@Description	举报
//	@Tags			order模块
//	@Accept			json
//	@Produce		json
//	@Param			request			body		ReportOrderReq	true	"请求"
//	@Param			Authorization	header		string			true	"token"
//	@Success		200				{object}	ReportOrderResp
//	@Failure		400				{object}	response.response
//	@Failure		404				{object}	response.response
//	@Failure		500				{object}	response.response
//	@Router			/order/reporter [post]
func (o *OrderSvc) ReportOrder(ctx context.Context, req *ReportOrderReq) (*ReportOrderResp, errors.Error) {
	//只有进行中中的订单可以举报
	order, err := o.orderRepo.GetSingleOrder(ctx, &GetOrderReq{OrderId: req.OrderId})
	if err != nil {
		return nil, err
	}

	if order.Order.Status != STATUS_PROCESS {
		return nil, errors.New(errors.ORDER_STATUS_ERROR, "只有进行中中的订单可以举报")
	}

	return o.orderRepo.AddReporter(ctx, req)
}

//	@Summary		上传设计图稿
//	@Description	乙方上传设计图稿
//	@Tags			order模块
//	@Accept			json
//	@Produce		json
//	@Param			request			body		UploadDesignArtworkReq	true	"请求"
//	@Param			Authorization	header		string					true	"token"
//	@Success		200				{object}	UploadDesignArtworkResp
//	@Failure		400				{object}	response.response
//	@Failure		404				{object}	response.response
//	@Failure		500				{object}	response.response
//	@Router			/order/designArtwork [post]
func (o *OrderSvc) UploadDesignArtwork(ctx context.Context, req *UploadDesignArtworkReq) (*UploadDesignArtworkResp, errors.Error) {
	//todo:存到图片到oss并返回url

	//存储url到数据库
	_, err := o.orderRepo.UpdateDesignArtwork(ctx, &UpdateDesignArtworkReq{
		OrderId: req.OrderId,
		Url:     "", //todo:url
	})

	return &UploadDesignArtworkResp{}, err
}

//	@Summary		上传版型图稿
//	@Description	乙方上传版型图稿
//	@Tags			order模块
//	@Accept			json
//	@Produce		json
//	@Param			request			body		UploadPatternArtworkReq	true	"请求"
//	@Param			Authorization	header		string					true	"token"
//	@Success		200				{object}	UploadPatternArtworkResp
//	@Failure		400				{object}	response.response
//	@Failure		404				{object}	response.response
//	@Failure		500				{object}	response.response
//	@Router			/order/patternArtwork [post]
func (o *OrderSvc) UploadPatternArtwork(ctx context.Context, req *UploadPatternArtworkReq) (*UploadPatternArtworkResp, errors.Error) {
	//todo:存到图片到oss并返回url

	//存储url到数据库
	_, err := o.orderRepo.UpdatePatternArtwork(ctx, &UpdatePatternArtworkReq{
		OrderId: req.OrderId,
		Url:     "", //todo:url
	})

	return &UploadPatternArtworkResp{}, err
}

//	@Summary		上传版型制作工艺
//	@Description	乙方上传版型制作工艺
//	@Tags			order模块
//	@Accept			json
//	@Produce		json
//	@Param			request			body		UploadPatternMakingProcessReq	true	"请求"
//	@Param			Authorization	header		string							true	"token"
//	@Success		200				{object}	UploadPatternMakingProcessResp
//	@Failure		400				{object}	response.response
//	@Failure		404				{object}	response.response
//	@Failure		500				{object}	response.response
//	@Router			/order/patternMakingProcess [post]
func (o *OrderSvc) UploadPatternMakingProcess(ctx context.Context, req *UploadPatternMakingProcessReq) (*UploadPatternMakingProcessResp, errors.Error) {
	_, err := o.orderRepo.UpdatePatternMakingProcess(ctx, &UpdatePatternMakingProcessReq{
		OrderId: req.OrderId,
		Content: req.Content,
	})

	return &UploadPatternMakingProcessResp{}, err
}

//	@Summary		上传样品成衣图
//	@Description	乙方上传样品成衣图
//	@Tags			order模块
//	@Accept			json
//	@Produce		json
//	@Param			request			body		UploadSampleImageReq	true	"请求"
//	@Param			Authorization	header		string					true	"token"
//	@Success		200				{object}	UploadSampleImageResp
//	@Failure		400				{object}	response.response
//	@Failure		404				{object}	response.response
//	@Failure		500				{object}	response.response
//	@Router			/order/sampleImage [post]
func (o *OrderSvc) UploadSampleImage(ctx context.Context, req *UploadSampleImageReq) (*UploadSampleImageResp, errors.Error) {
	//todo:存到图片到oss并返回url

	_, err := o.orderRepo.UpdateSampleImage(ctx, &UpdateSampleImageReq{
		OrderId: req.OrderId,
		Url:     "", //todo:url
	})

	return &UploadSampleImageResp{}, err
}

//	@Summary		上传模特展示视频
//	@Description	乙方上传模特展示视频
//	@Tags			order模块
//	@Accept			json
//	@Produce		json
//	@Param			request			body		UploadShowVideoReq	true	"请求"
//	@Param			Authorization	header		string				true	"token"
//	@Success		200				{object}	UploadShowVideoResp
//	@Failure		400				{object}	response.response
//	@Failure		404				{object}	response.response
//	@Failure		500				{object}	response.response
//	@Router			/order/showVideo [post]
func (o *OrderSvc) UploadShowVideo(ctx context.Context, req *UploadShowVideoReq) (*UploadShowVideoResp, errors.Error) {
	//todo:存到视频到oss并返回url

	_, err := o.orderRepo.UpdateShowVideo(ctx, &UpdateShowVideoReq{
		OrderId: req.OrderId,
		Url:     "", //todo:url
	})

	return &UploadShowVideoResp{}, err
}

//	@Summary		查询总交易订单次数+总交易金额
//	@Description	首页中，查询总交易订单次数+总交易金额
//	@Tags			order模块
//	@Accept			json
//	@Produce		json
//	@Param			None			query		string	false	"无参数"
//	@Param			Authorization	header		string	true	"token"
//	@Success		200				{object}	GetOrderSumResp
//	@Failure		400				{object}	response.response
//	@Failure		404				{object}	response.response
//	@Failure		500				{object}	response.response
//	@Router			/order/sum [get]
func (o *OrderSvc) GetOrderSum(ctx context.Context, req *GetOrderSumReq) (*GetOrderSumResp, errors.Error) {
	return o.orderRepo.GetOrderSum(ctx, req)
}

//	@Summary		发布需求
//	@Description	甲方发布需求
//	@Tags			order模块
//	@Accept			json
//	@Produce		json
//	@Param			request			body		PublishOrderReq	true	"请求"
//	@Param			Authorization	header		string			true	"token"
//	@Success		200				{object}	PublishOrderResp
//	@Failure		400				{object}	response.response
//	@Failure		404				{object}	response.response
//	@Failure		500				{object}	response.response
//	@Router			/order/sum [post]
func (o *OrderSvc) PublishOrder(ctx context.Context, req *PublishOrderReq) (*PublishOrderResp, errors.Error) {
	return o.orderRepo.AddOrder(ctx, req)
}
