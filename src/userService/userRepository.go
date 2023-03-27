package userService

import (
	"context"
	"customClothing/src/db"
	errors "customClothing/src/error"
	"customClothing/src/utils/uuid"
	"fmt"
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
	WithdrawMarginApplication(c context.Context, req *WithdrawMarginApplicationReq) (*WithdrawMarginApplicationResp, errors.Error)
	DeductMargin(c context.Context, req *DeductMarginReq) (*DeductMarginResp, errors.Error)
	UpdateUserStatus(c context.Context, req *UpdateUserStatusReq) (*UpdateUserStatusResp, errors.Error)
	ReviewMarginWithdrawApplication(c context.Context, req *ReviewMarginWithdrawApplicationReq) (*ReviewMarginWithdrawApplicationResp, errors.Error)
}

type userRepoSvc struct {
	userDb              *gorm.DB
	marginWithdrawAppDb *gorm.DB
}

func MakeUserRepoService() UserRepoService {
	return &userRepoSvc{
		userDb:              db.Db().Table(User{}.TableName()),
		marginWithdrawAppDb: db.Db().Table(MarginWithdrawApplication{}.TableName()),
	}
}

func (r *userRepoSvc) AddUser(c context.Context, req *RegisterUserReq) (*RegisterUserResp, errors.Error) {
	//检查用户是否已存在
	var count int64
	if err := r.userDb.Where("phone = ?", req.Phone).Count(&count).Error; err != nil {
		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	if count > 0 {
		return nil, errors.New(errors.INTERNAL_ERROR, "用户已存在，请使用不同的手机号注册")
	}

	//新增用户
	user := User{
		UserId:      uuid.BuildUuid(),
		DisplayName: req.DisplayName,
		Password:    req.Password,
		Avatar:      "", //todo
		RoleCode:    ROLE_NORMAL,
		Phone:       req.Phone,
	}
	if err := r.userDb.Create(&user).Error; err != nil {
		fmt.Println(err)
		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	return &RegisterUserResp{
		Code: 0,
		Msg:  "注册成功",
	}, nil
}

func (r *userRepoSvc) GetUserByPhone(c context.Context, phone string) (*User, errors.Error) {
	u := User{}
	if err := r.userDb.Where("phone = ?", phone).First(&u); err != nil {
		return nil, errors.New(errors.USER_NOT_EXIST, "用户不存在")
	}

	return &u, nil
}

func (r *userRepoSvc) GetRoleByCode(c context.Context, code int) (*Role, errors.Error) {
	role := &Role{}
	if err := r.userDb.Where("code = ?", code); err != nil {
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

	if err := r.userDb.Where("user_id = ?", id).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.USER_NOT_EXIST, "该用户不存在")
		}
		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	return &u, nil
}

func (r *userRepoSvc) AddMargin(c context.Context, req *AddMarginReq) (*AddMarginResp, errors.Error) {
	u := User{}
	if err := r.userDb.Where("user_id = ?", req.UserId).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.USER_NOT_EXIST, "用户不存在")
		}
		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	u.Margin += req.Amount
	r.userDb.Save(&u)

	return &AddMarginResp{}, nil
}

func (r *userRepoSvc) GetMargin(c context.Context, id string) (uint, errors.Error) {
	u := User{}
	if err := r.userDb.Where("user_id = ?", id).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, errors.New(errors.USER_NOT_EXIST, "用户不存在")
		}
		return 0, errors.New(errors.INTERNAL_ERROR, "")
	}

	return u.Margin, nil
}

func (r *userRepoSvc) WithdrawMarginApplication(c context.Context, req *WithdrawMarginApplicationReq) (*WithdrawMarginApplicationResp, errors.Error) {
	app := MarginWithdrawApplication{
		ApplicationId: uuid.BuildUuid(),
		UserId:        req.UserId,
		Amount:        req.Amount,
		Status:        MarginWithdrawStatus[1],
	}

	if err := r.marginWithdrawAppDb.Create(&app).Error; err != nil {
		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	return &WithdrawMarginApplicationResp{}, nil
}

func (r *userRepoSvc) DeductMargin(c context.Context, req *DeductMarginReq) (*DeductMarginResp, errors.Error) {
	u := User{}
	if err := r.userDb.Where("user_id = ?", req.UserId).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.USER_NOT_EXIST, "用户不存在")
		}
		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	u.Margin -= req.Amount
	r.userDb.Save(&u)

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

	if err := r.userDb.Where("user_id = ?", req.UserId).Update("ban = ?", status).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.USER_NOT_EXIST, "用户不存在")
		}
		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	return &UpdateUserStatusResp{}, nil
}

func (r *userRepoSvc) ReviewMarginWithdrawApplication(c context.Context, req *ReviewMarginWithdrawApplicationReq) (*ReviewMarginWithdrawApplicationResp, errors.Error) {
	if err := r.marginWithdrawAppDb.Where("application_id = ?", req.ApplicationId).Update("status = ?", MarginWithdrawStatus[req.Status]).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.MARGIN_WITHDRAW_APPLICATION_NOT_EXIST, "保证金提现申请记录不存在")
		}
		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	if req.Status == 2 {
		u := User{}
		if err := r.userDb.Where("user_id = ?", req.UserId).First(&u).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, errors.New(errors.USER_NOT_EXIST, "用户不存在")
			}
			return nil, errors.New(errors.INTERNAL_ERROR, "")
		}

		u.Margin -= req.Amount
		r.userDb.Save(&u)

		//todo:转账给用户
	}

	return &ReviewMarginWithdrawApplicationResp{}, nil
}
