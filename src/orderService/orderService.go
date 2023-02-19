package orderService

import (
	"context"
	errors "customClothing/src/error"
)

type OrderService interface {
	ListOrders(ctx context.Context, req *ListOrdersReq) (*ListOrdersResp, errors.Error)
	GetSingleOrder(ctx context.Context, req *GetOrderReq) (*GetOrderResp, errors.Error)
	UpdateCost(ctx context.Context, req *UpdateCostReq) (*UpdateCostResp, errors.Error)
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

}
