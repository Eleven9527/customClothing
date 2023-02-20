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
	CancelOrder(ctx context.Context, req *CancelOrderReq) (*CancelOrderResp, errors.Error)
	ConfirmOrder(ctx context.Context, req *ConfirmOrderReq) (*ConfirmOrderResp, errors.Error)
	AddReporter(ctx context.Context, req *ReportOrderReq) (*ReportOrderResp, errors.Error)
	UpdateDesignArtwork(ctx context.Context, req *UpdateDesignArtworkReq) (*UpdateDesignArtworkResp, errors.Error)
	UpdatePatternArtwork(ctx context.Context, req *UpdatePatternArtworkReq) (*UpdatePatternArtworkResp, errors.Error)
	UpdatePatternMakingProcess(ctx context.Context, req *UpdatePatternMakingProcessReq) (*UpdatePatternMakingProcessResp, errors.Error)
	UpdateSampleImage(ctx context.Context, req *UpdateSampleImageReq) (*UpdateSampleImageResp, errors.Error)
	UpdateShowVideo(ctx context.Context, req *UpdateShowVideoReq) (*UpdateShowVideoResp, errors.Error)
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

func (r *repoSvc) CancelOrder(ctx context.Context, req *CancelOrderReq) (*CancelOrderResp, errors.Error) {
	if err := db.Db().Where("order_id = ?", req.OrderId).Update("status = ?", STATUS_CANCEL).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.ORDER_NOT_EXIST, "订单不存在")
		}

		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	return &CancelOrderResp{}, nil
}

func (r *repoSvc) ConfirmOrder(ctx context.Context, req *ConfirmOrderReq) (*ConfirmOrderResp, errors.Error) {
	if err := db.Db().Where("order_id = ?", req.OrderId).Update("status = ?", STATUS_COMPLETE).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.ORDER_NOT_EXIST, "订单不存在")
		}

		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	//todo:把托管的钱转入乙方钱包

	return &ConfirmOrderResp{}, nil
}

func (r *repoSvc) AddReporter(ctx context.Context, req *ReportOrderReq) (*ReportOrderResp, errors.Error) {
	reporter := Reporter{
		Whistleblower:  req.Whistleblower,
		ReportedPerson: req.ReportedPerson,
		Description:    req.Description,
		OrderId:        req.OrderId,
	}

	if err := db.Db().Create(&reporter).Error; err != nil {
		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	return &ReportOrderResp{}, nil
}

func (r *repoSvc) UpdateDesignArtwork(ctx context.Context, req *UpdateDesignArtworkReq) (*UpdateDesignArtworkResp, errors.Error) {
	if err := db.Db().Where("order_id = ?", req.OrderId).Update("design_artwork = ?", req.Url).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.ORDER_NOT_EXIST, "订单不存在")
		}
		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	return &UpdateDesignArtworkResp{}, nil
}

func (r *repoSvc) UpdatePatternArtwork(ctx context.Context, req *UpdatePatternArtworkReq) (*UpdatePatternArtworkResp, errors.Error) {
	if err := db.Db().Where("order_id = ?", req.OrderId).Update("pattern_artwork = ?", req.Url).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.ORDER_NOT_EXIST, "订单不存在")
		}
		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	return &UpdatePatternArtworkResp{}, nil
}

func (r *repoSvc) UpdatePatternMakingProcess(ctx context.Context, req *UpdatePatternMakingProcessReq) (*UpdatePatternMakingProcessResp, errors.Error) {
	if err := db.Db().Where("order_id = ?", req.OrderId).Update("pattern_making_process = ?", req.Content).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.ORDER_NOT_EXIST, "订单不存在")
		}
		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	return &UpdatePatternMakingProcessResp{}, nil
}

func (r *repoSvc) UpdateSampleImage(ctx context.Context, req *UpdateSampleImageReq) (*UpdateSampleImageResp, errors.Error) {
	if err := db.Db().Where("order_id = ?", req.OrderId).Update("sample_image = ?", req.Url).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.ORDER_NOT_EXIST, "订单不存在")
		}
		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	return &UpdateSampleImageResp{}, nil
}

func (r *repoSvc) UpdateShowVideo(ctx context.Context, req *UpdateShowVideoReq) (*UpdateShowVideoResp, errors.Error) {
	if err := db.Db().Where("order_id = ?", req.OrderId).Update("show_video = ?", req.Url).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.ORDER_NOT_EXIST, "订单不存在")
		}
		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	return &UpdateShowVideoResp{}, nil
}
