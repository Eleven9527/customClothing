package orderService

import (
	"context"
	"customClothing/src/db"
	errors "customClothing/src/error"
	"customClothing/src/userService"
	"github.com/jinzhu/gorm"
)

type RepoService interface {
	ListOrders(ctx context.Context, req *ListOrdersReq) (*ListOrdersResp, errors.Error)
	GetSingleOrder(ctx context.Context, req *GetOrderReq) (*GetOrderResp, errors.Error)
	UpdateCost(ctx context.Context, req *UpdateCostReq) (*UpdateCostResp, errors.Error)
}

type repoSvc struct {
	db *gorm.DB
}

func MakeRepoService() RepoService {
	return &repoSvc{
		db: db.Db().Table(Order{}.TableName()),
	}
}

func (r *repoSvc) ListOrders(ctx context.Context, req *ListOrdersReq) (*ListOrdersResp, errors.Error) {
	var count int64
	if err := db.Db().Where("id >0").Count(&count); err != nil {
		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	list := make([]*Order, 0)
	switch req.RoleCode {
	case userService.ROLE_DESIGNER, userService.ROLE_PATTERN: //设计师、版型师
		if err := db.Db().Where("part_b = ? AND ( 0 = ? OR status = ?)", req.UserId, req.Status, req.Status).
			Offset(req.PageSize * (req.PageNum - 1)).Limit(req.PageSize).Error; err != nil && err != gorm.ErrRecordNotFound {
			return nil, errors.New(errors.INTERNAL_ERROR, "")
		}
	case userService.ROLE_NORMAL: //普通用户
		if err := db.Db().Where("part_a = ? AND ( 0 = ? OR status = ?)", req.UserId, req.Status, req.Status).
			Offset(req.PageSize * (req.PageNum - 1)).Limit(req.PageSize).Error; err != nil && err != gorm.ErrRecordNotFound {
			return nil, errors.New(errors.INTERNAL_ERROR, "")
		}
	}

	return &ListOrdersResp{
		Total:  count,
		Orders: list,
	}, nil
}

func (r *repoSvc) GetSingleOrder(ctx context.Context, req *GetOrderReq) (*GetOrderResp, errors.Error) {
	o := Order{}

	if err := db.Db().Where("order_id = ?", req.OrderId).First(&o).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.ORDER_NOT_EXIST, "订单不存在")
		}
		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	return &GetOrderResp{Order: &o}, nil
}

func (r *repoSvc) UpdateCost(ctx context.Context, req *UpdateCostReq) (*UpdateCostResp, errors.Error) {
	if err := db.Db().Where("order_id = ?", req.OrderId).Update("cost = ?", req.Cost).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.ORDER_NOT_EXIST, "订单不存在")
		}
		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	return &UpdateCostResp{}, nil
}
