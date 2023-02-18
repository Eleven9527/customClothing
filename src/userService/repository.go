package userService

import (
	"context"
	"customClothing/src/db"
	errors "customClothing/src/error"
	"github.com/jinzhu/gorm"
)

type RepoService interface {
	AddUser(c context.Context, req *RegisterUserReq) (*RegisterUserResp, errors.Error)
	GetUserByPhone(c context.Context, phone string) (*User, errors.Error)
}

type repoSvc struct {
	db *gorm.DB
}

func MakeRepoService() RepoService {
	return &repoSvc{
		db: db.Db(),
	}
}

func (r *repoSvc) AddUser(c context.Context, req *RegisterUserReq) (*RegisterUserResp, errors.Error) {
	//检查用户是否已存在
	var count int64
	if err := r.db.Where("phone = ?", req.Phone).Count(&count); err != nil {
		return nil, errors.New(errors.USER_EXIST, "")
	}

	if count > 0 {
		return nil, errors.New(errors.INTERNAL_ERROR, "用户已存在，请使用不同的手机号注册")
	}

	//新增用户
	if err := r.db.Create(&req).Error; err != nil {
		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	return &RegisterUserResp{
		Code: 0,
		Msg:  "注册成功",
	}, nil
}

func (r *repoSvc) GetUserByPhone(c context.Context, phone string) (*User, errors.Error) {
	u := User{}
	if err := r.db.Where("phone = ?", phone).First(&u); err != nil {
		return nil, errors.New(errors.USER_NOT_EXIST, "用户不存在")
	}

	return &u, nil
}
