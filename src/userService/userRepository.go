package userService

import (
	"context"
	"customClothing/src/db"
	errors "customClothing/src/error"
	"github.com/jinzhu/gorm"
)

type UserRepoService interface {
	AddUser(c context.Context, req *RegisterUserReq) (*RegisterUserResp, errors.Error)
	GetUserByPhone(c context.Context, phone string) (*User, errors.Error)
	GetUserById(c context.Context, id string) (*User, errors.Error)
	GetRoleByCode(c context.Context, code int) (*Role, errors.Error)
	AddKyc(c context.Context, req *Kyc) errors.Error
	AddMargin(c context.Context, req *AddMarginReq) (*AddMarginResp, errors.Error)
	GetMargin(c context.Context, id string) (uint, errors.Error)
	WithdrawMargin(c context.Context, req *WithdrawMarginReq) (*WithdrawMarginResp, errors.Error)
	DeductMargin(c context.Context, req *DeductMarginReq) (*DeductMarginResp, errors.Error)
	UpdateUserStatus(c context.Context, req *UpdateUserStatusReq) (*UpdateUserStatusResp, errors.Error)
}

type userRepoSvc struct {
	db *gorm.DB
}

func MakeUserRepoService() UserRepoService {
	return &userRepoSvc{
		db: db.Db().Table(User{}.TableName()),
	}
}

func (r *userRepoSvc) AddUser(c context.Context, req *RegisterUserReq) (*RegisterUserResp, errors.Error) {
	//检查用户是否已存在
	var count int64
	if err := r.db.Where("phone = ?", req.Phone).Count(&count); err != nil {
		return nil, errors.New(errors.INTERNAL_ERROR, "")
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

func (r *userRepoSvc) GetUserByPhone(c context.Context, phone string) (*User, errors.Error) {
	u := User{}
	if err := r.db.Where("phone = ?", phone).First(&u); err != nil {
		return nil, errors.New(errors.USER_NOT_EXIST, "用户不存在")
	}

	return &u, nil
}

func (r *userRepoSvc) GetRoleByCode(c context.Context, code int) (*Role, errors.Error) {
	role := &Role{}
	if err := r.db.Where("code = ?", code); err != nil {
		return nil, errors.New(errors.ROLE_NOT_EXIST, "角色不存在")
	}

	return role, nil
}
func (r *userRepoSvc) AddKyc(c context.Context, k *Kyc) errors.Error {
	if err := db.Db().Create(k).Error; err != nil {
		return errors.New(errors.INTERNAL_ERROR, "")
	}

	return nil
}

// GetUserById 根据uuid获取用户
func (r *userRepoSvc) GetUserById(c context.Context, id string) (*User, errors.Error) {
	u := User{}

	if err := r.db.Where("user_id = ?", id).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.USER_NOT_EXIST, "该用户不存在")
		}
		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	return &u, nil
}

func (r *userRepoSvc) AddMargin(c context.Context, req *AddMarginReq) (*AddMarginResp, errors.Error) {
	u := User{}
	if err := r.db.Where("user_id = ?", req.UserId).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.USER_NOT_EXIST, "用户不存在")
		}
		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	u.Margin += req.Amount
	r.db.Save(&u)

	return &AddMarginResp{}, nil
}

func (r *userRepoSvc) GetMargin(c context.Context, id string) (uint, errors.Error) {
	u := User{}
	if err := r.db.Where("user_id = ?", id).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, errors.New(errors.USER_NOT_EXIST, "用户不存在")
		}
		return 0, errors.New(errors.INTERNAL_ERROR, "")
	}

	return u.Margin, nil
}

func (r *userRepoSvc) WithdrawMargin(c context.Context, req *WithdrawMarginReq) (*WithdrawMarginResp, errors.Error) {
	u := User{}
	if err := r.db.Where("user_id = ?", req.UserId).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.USER_NOT_EXIST, "用户不存在")
		}
		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	//todo:转账给用户

	u.Margin -= req.Amount
	r.db.Save(&u)

	return &WithdrawMarginResp{}, nil
}

func (r *userRepoSvc) DeductMargin(c context.Context, req *DeductMarginReq) (*DeductMarginResp, errors.Error) {
	u := User{}
	if err := r.db.Where("user_id = ?", req.UserId).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.USER_NOT_EXIST, "用户不存在")
		}
		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	u.Margin -= req.Amount
	r.db.Save(&u)

	return &DeductMarginResp{}, nil
}

func (r *userRepoSvc) UpdateUserStatus(c context.Context, req *UpdateUserStatusReq) (*UpdateUserStatusResp, errors.Error) {
	status := false

	switch req.Status {
	case 1: //拉黑
		status = true
	case 2: //接触拉黑
		status = false
	}

	if err := r.db.Where("user_id = ?", req.UserId).Update("ban = ?", status).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.USER_NOT_EXIST, "用户不存在")
		}
		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	return &UpdateUserStatusResp{}, nil
}
