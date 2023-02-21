package userService

import (
	"context"
	"customClothing/src/db"
	errors "customClothing/src/error"
	"github.com/jinzhu/gorm"
)

type MarginRecordRepoService interface {
	AddMarginRecord(c context.Context, req *AddMarginRecordReq) (*AddMarginRecordResp, errors.Error)
}

type marginRecordRepoSvc struct {
	db *gorm.DB
}

func MakeMarginRecordRepoService() MarginRecordRepoService {
	return &marginRecordRepoSvc{
		db: db.Db().Table(MarginRecord{}.TableName()),
	}
}

func (r *marginRecordRepoSvc) AddMarginRecord(c context.Context, req *AddMarginRecordReq) (*AddMarginRecordResp, errors.Error) {
	record := MarginRecord{
		UserId:      req.UserId,
		Amount:      req.Amount,
		OperateType: MARGIN_OP_ADD,
	}

	if err := r.db.Where("user_id = ?", req.UserId).Create(&record).Error; err != nil {
		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	return &AddMarginRecordResp{}, nil
}
