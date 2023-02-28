package orderService

import (
	"context"
	"customClothing/src/db"
	errors "customClothing/src/error"
	"customClothing/src/userService"
	"customClothing/src/utils/uuid"
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
	GetOrderSum(ctx context.Context, req *GetOrderSumReq) (*GetOrderSumResp, errors.Error)
	AddOrder(ctx context.Context, req *PublishOrderReq) (*PublishOrderResp, errors.Error)
	GetOrderDetail(ctx context.Context, req *GetOrderDetailReq) (*GetOrderDetailResp, errors.Error)
	UpdatePartB(ctx context.Context, userId, orderId string) errors.Error
	DeleteOrder(ctx context.Context, orderId string) errors.Error
	DeleteOrderDetail(ctx context.Context, orderId string) errors.Error
	ListReporters(ctx context.Context, req *ListReportersReq) (*ListReportersResp, errors.Error)
}

type repoSvc struct {
	orderTable    *gorm.DB
	detalTable    *gorm.DB
	reporterTable *gorm.DB
}

func MakeRepoService() RepoService {
	return &repoSvc{
		orderTable:    db.Db().Table(Order{}.TableName()),
		detalTable:    db.Db().Table(OrderDetail{}.TableName()),
		reporterTable: db.Db().Table(Reporter{}.TableName()),
	}
}

func (r *repoSvc) ListOrders(ctx context.Context, req *ListOrdersReq) (*ListOrdersResp, errors.Error) {
	var count int64
	if err := r.orderTable.Where("id >0").Count(&count); err != nil {
		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	list := make([]*Order, 0)
	switch req.RoleCode {
	case userService.ROLE_DESIGNER, userService.ROLE_PATTERN: //设计师、版型师
		if err := db.Db().Where("part_b = ? AND ( 0 = ? OR status = ?)", req.UserId, OrderStatus[req.Status], OrderStatus[req.Status]).
			Offset(req.PageSize * (req.PageNum - 1)).Limit(req.PageSize).Error; err != nil && err != gorm.ErrRecordNotFound {
			return nil, errors.New(errors.INTERNAL_ERROR, "")
		}
	case userService.ROLE_NORMAL: //普通用户
		if err := db.Db().Where("part_a = ? AND ( 0 = ? OR status = ?)", req.UserId, OrderStatus[req.Status], OrderStatus[req.Status]).
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

	if err := r.orderTable.Where("order_id = ?", req.OrderId).First(&o).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.ORDER_NOT_EXIST, "订单不存在")
		}
		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	return &GetOrderResp{Order: &o}, nil
}

func (r *repoSvc) UpdateCost(ctx context.Context, req *UpdateCostReq) (*UpdateCostResp, errors.Error) {
	if err := r.orderTable.Where("order_id = ?", req.OrderId).Update("cost = ?", req.Cost).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.ORDER_NOT_EXIST, "订单不存在")
		}
		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	return &UpdateCostResp{}, nil
}

func (r *repoSvc) CancelOrder(ctx context.Context, req *CancelOrderReq) (*CancelOrderResp, errors.Error) {
	if err := r.orderTable.Where("order_id = ?", req.OrderId).Update("status = ?", STATUS_CANCEL).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.ORDER_NOT_EXIST, "订单不存在")
		}

		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	return &CancelOrderResp{}, nil
}

func (r *repoSvc) ConfirmOrder(ctx context.Context, req *ConfirmOrderReq) (*ConfirmOrderResp, errors.Error) {
	if err := r.orderTable.Where("order_id = ?", req.OrderId).Update("status = ?", STATUS_COMPLETE).Error; err != nil {
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

	if err := r.reporterTable.Create(&reporter).Error; err != nil {
		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	return &ReportOrderResp{}, nil
}

func (r *repoSvc) UpdateDesignArtwork(ctx context.Context, req *UpdateDesignArtworkReq) (*UpdateDesignArtworkResp, errors.Error) {
	if err := r.orderTable.Where("order_id = ?", req.OrderId).Update("design_artwork = ?", req.Url).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.ORDER_NOT_EXIST, "订单不存在")
		}
		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	return &UpdateDesignArtworkResp{}, nil
}

func (r *repoSvc) UpdatePatternArtwork(ctx context.Context, req *UpdatePatternArtworkReq) (*UpdatePatternArtworkResp, errors.Error) {
	if err := r.orderTable.Where("order_id = ?", req.OrderId).Update("pattern_artwork = ?", req.Url).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.ORDER_NOT_EXIST, "订单不存在")
		}
		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	return &UpdatePatternArtworkResp{}, nil
}

func (r *repoSvc) UpdatePatternMakingProcess(ctx context.Context, req *UpdatePatternMakingProcessReq) (*UpdatePatternMakingProcessResp, errors.Error) {
	if err := r.orderTable.Where("order_id = ?", req.OrderId).Update("pattern_making_process = ?", req.Content).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.ORDER_NOT_EXIST, "订单不存在")
		}
		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	return &UpdatePatternMakingProcessResp{}, nil
}

func (r *repoSvc) UpdateSampleImage(ctx context.Context, req *UpdateSampleImageReq) (*UpdateSampleImageResp, errors.Error) {
	if err := r.orderTable.Where("order_id = ?", req.OrderId).Update("sample_image = ?", req.Url).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.ORDER_NOT_EXIST, "订单不存在")
		}
		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	return &UpdateSampleImageResp{}, nil
}

func (r *repoSvc) UpdateShowVideo(ctx context.Context, req *UpdateShowVideoReq) (*UpdateShowVideoResp, errors.Error) {
	if err := r.orderTable.Where("order_id = ?", req.OrderId).Update("show_video = ?", req.Url).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.ORDER_NOT_EXIST, "订单不存在")
		}
		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	return &UpdateShowVideoResp{}, nil
}

func (r *repoSvc) GetOrderSum(ctx context.Context, req *GetOrderSumReq) (*GetOrderSumResp, errors.Error) {
	orders := make([]*Order, 0)

	if err := r.orderTable.Where("id > 0").Find(&orders).Error; err != nil {
		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	var amount float64
	for _, v := range orders {
		amount += v.Cost
	}

	return &GetOrderSumResp{
		Sum:    len(orders),
		Amount: amount,
	}, nil
}

func (r *repoSvc) AddOrder(ctx context.Context, req *PublishOrderReq) (*PublishOrderResp, errors.Error) {
	//添加订单
	order := Order{
		OrderId: uuid.BuildUuid(),
		PartA:   req.UserId,
		Status:  STATUS_PENDING,
		Cost:    req.Cost,
		Model:   gorm.Model{},
	}

	if err := r.orderTable.Create(&order).Error; err != nil {
		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	//添加订单详情
	detail := OrderDetail{
		OrderId:                order.OrderId,
		Cost:                   req.Cost,
		DeadLine:               req.DeadLine,
		ProvideFabric:          req.ProvideFabric,
		Part:                   Part[req.Part],
		ReferenceSize:          req.ReferenceSize,
		WearingOccasion:        WearingOccasion[req.WearingOccasion],
		PatternRequest:         req.PatternRequest,
		Style:                  Style[req.Style],
		ConsumptionPositioning: req.ConsumptionPositioning,
		ProductionTime:         req.ProductionTime,
		SampleCost:             req.SampleCost,
		PatternCost:            req.PatternCost,
		Address:                req.Address,
		Description:            req.Description,
	}

	//自填风格描述
	if detail.Style == "" {
		detail.Style = req.CustomStyle
	}

	//需要确认的步骤
	if len(req.ConfirmSteps) > 0 {
		steps := ""
		for _, v := range req.ConfirmSteps {
			steps += Step[v]
			steps += ","
		}
		//去除末尾逗号
		if len(steps) > 0 {
			steps = steps[:len(steps)-1]
		}
		detail.ConfirmSteps = steps
	}

	if err := r.detalTable.Create(&detail).Error; err != nil {
		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	return &PublishOrderResp{}, nil
}

func (r *repoSvc) GetOrderDetail(ctx context.Context, req *GetOrderDetailReq) (*GetOrderDetailResp, errors.Error) {
	d := OrderDetail{}

	if err := r.detalTable.Where("order_id = ?", req.OrderId).First(&d).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.ORDER_DETAIL_NOT_EXIST, "订单详情不存在")
		}
		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	return &GetOrderDetailResp{Detail: &d}, nil
}

func (r *repoSvc) UpdatePartB(ctx context.Context, userId, orderId string) errors.Error {
	if err := r.orderTable.Where("order_id = ?", orderId).Update("part_b = ?", userId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New(errors.ORDER_NOT_EXIST, "订单不存在")
		}
		return errors.New(errors.INTERNAL_ERROR, "")
	}

	return nil
}

func (r *repoSvc) DeleteOrder(ctx context.Context, orderId string) errors.Error {
	o := Order{OrderId: orderId}
	if err := r.orderTable.Delete(&o).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New(errors.ORDER_NOT_EXIST, "订单不存在")
		}
		return errors.New(errors.INTERNAL_ERROR, "")
	}

	return nil
}

func (r *repoSvc) DeleteOrderDetail(ctx context.Context, orderId string) errors.Error {
	o := OrderDetail{OrderId: orderId}
	if err := r.detalTable.Delete(&o).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New(errors.ORDER_NOT_EXIST, "订单不存在")
		}
		return errors.New(errors.INTERNAL_ERROR, "")
	}

	return nil
}

func (r *repoSvc) ListReporters(ctx context.Context, req *ListReportersReq) (*ListReportersResp, errors.Error) {
	list := make([]*Reporter, 0)
	var count int64

	if err := r.reporterTable.Where("deleted_at != null").Count(&count).Error; err != nil {
		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	if err := r.reporterTable.Where("created_at >= ? AND created_at <= ?", req.StartTime, req.EndTime).Offset(req.PageSize * (req.PageNum - 1)).Error; err != nil {
		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	return &ListReportersResp{
		Total: count,
		List:  list,
	}, nil
}
