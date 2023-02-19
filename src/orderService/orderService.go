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
