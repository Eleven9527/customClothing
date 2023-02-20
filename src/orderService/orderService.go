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

func (o *OrderSvc) ListOrders(ctx context.Context, req *ListOrdersReq) (*ListOrdersResp, errors.Error) {
	return o.orderRepo.ListOrders(ctx, req)
}

func (o *OrderSvc) GetSingleOrder(ctx context.Context, req *GetOrderReq) (*GetOrderResp, errors.Error) {
	return o.orderRepo.GetSingleOrder(ctx, req)
}

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

func (o *OrderSvc) UploadDesignArtwork(ctx context.Context, req *UploadDesignArtworkReq) (*UploadDesignArtworkResp, errors.Error) {
	//todo:存到图片到oss并返回url

	//存储url到数据库
	_, err := o.orderRepo.UpdateDesignArtwork(ctx, &UpdateDesignArtworkReq{
		OrderId: req.OrderId,
		Url:     "", //todo:url
	})

	return &UploadDesignArtworkResp{}, err
}

func (o *OrderSvc) UploadPatternArtwork(ctx context.Context, req *UploadPatternArtworkReq) (*UploadPatternArtworkResp, errors.Error) {
	//todo:存到图片到oss并返回url

	//存储url到数据库
	_, err := o.orderRepo.UpdatePatternArtwork(ctx, &UpdatePatternArtworkReq{
		OrderId: req.OrderId,
		Url:     "", //todo:url
	})

	return &UploadPatternArtworkResp{}, err
}

func (o *OrderSvc) UploadPatternMakingProcess(ctx context.Context, req *UploadPatternMakingProcessReq) (*UploadPatternMakingProcessResp, errors.Error) {
	_, err := o.orderRepo.UpdatePatternMakingProcess(ctx, &UpdatePatternMakingProcessReq{
		OrderId: req.OrderId,
		Content: req.Content,
	})

	return &UploadPatternMakingProcessResp{}, err
}

func (o *OrderSvc) UploadSampleImage(ctx context.Context, req *UploadSampleImageReq) (*UploadSampleImageResp, errors.Error) {
	//todo:存到图片到oss并返回url

	_, err := o.orderRepo.UpdateSampleImage(ctx, &UpdateSampleImageReq{
		OrderId: req.OrderId,
		Url:     "", //todo:url
	})

	return &UploadSampleImageResp{}, err
}

func (o *OrderSvc) UploadShowVideo(ctx context.Context, req *UploadShowVideoReq) (*UploadShowVideoResp, errors.Error) {
	//todo:存到视频到oss并返回url

	_, err := o.orderRepo.UpdateShowVideo(ctx, &UpdateShowVideoReq{
		OrderId: req.OrderId,
		Url:     "", //todo:url
	})

	return &UploadShowVideoResp{}, err
}
