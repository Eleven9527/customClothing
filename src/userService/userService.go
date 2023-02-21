package userService

import (
	"context"
	errors "customClothing/src/error"
	"customClothing/src/utils/token"
)

type UserService interface {
	RegisterUser(c context.Context, req *RegisterUserReq) (*RegisterUserResp, errors.Error)
	GetAuthCode(c context.Context, req *GetAuthCodeReq) (*GetAuthCodeResp, errors.Error)
	Login(c context.Context, req *LoginReq) (*LoginResp, errors.Error)
	Kyc(c context.Context, req *KycReq) (*KycResp, errors.Error)
	VerifyToken(tk []byte) errors.Error
	PayMargin(c context.Context, req *PayMarginReq) (*PayMarginResp, errors.Error)
	GetMargin(c context.Context, req *GetMarginReq) (*GetMarginResp, errors.Error)
	WithdrawMargin(c context.Context, req *WithdrawMarginReq) (*WithdrawMarginResp, errors.Error)
	DeductMargin(c context.Context, req *DeductMarginReq) (*DeductMarginResp, errors.Error)
}

type UserSvc struct {
	userRepo   UserRepoService
	marginRepo MarginRecordRepoService
	cache      CacheService
}

func MakeUserService() UserService {
	return &UserSvc{
		userRepo:   MakeUserRepoService(),
		marginRepo: MakeMarginRecordRepoService(),
		cache:      MakeCacheService(),
	}
}

func (u *UserSvc) RegisterUser(c context.Context, req *RegisterUserReq) (*RegisterUserResp, errors.Error) {
	return u.userRepo.AddUser(c, req)

}

// GetAuthCode todo:生成图片验证码
func (u *UserSvc) GetAuthCode(c context.Context, req *GetAuthCodeReq) (*GetAuthCodeResp, errors.Error) {
	return &GetAuthCodeResp{AuthCode: "12345"}, nil
}

func (u *UserSvc) Login(c context.Context, req *LoginReq) (*LoginResp, errors.Error) {
	//检测用户是否存在
	user, err := u.userRepo.GetUserByPhone(c, req.Phone)
	if err != nil {
		return nil, err
	}

	//生成token
	token, err := token.EncodeToken(user.Phone, user.DisplayName)
	if err != nil {
		return nil, err
	}

	//保存token到redis
	err = u.cache.SetToken(c, &SetTokenReq{Phone: user.Phone, Token: token})
	if err != nil {
		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	return &LoginResp{
		User:  user,
		Token: string(token),
	}, nil
}

func (u *UserSvc) Kyc(c context.Context, req *KycReq) (*KycResp, errors.Error) {
	role, err := u.userRepo.GetRoleByCode(c, req.RoleCode)
	if err != nil {
		return nil, err
	}

	//检测用户是否存在
	_, err = u.userRepo.GetUserByPhone(c, req.Phone)
	if err != nil {
		return nil, err
	}

	kyc := Kyc{
		UserId:   req.UserId,
		Name:     req.Name,
		Phone:    req.Phone,
		IdCard:   "", //todo:保存图片到oss，再返回url
		RoleCode: role.Code,
		Status:   KYC_STATUS_PENDING,
	}

	return nil, u.userRepo.AddKyc(c, &kyc)
}

func (u *UserSvc) VerifyToken(tk []byte) errors.Error {
	cacheTk, err := u.cache.GetToken(tk)
	if err != nil {
		return err
	}

	if cacheTk != string(tk) {
		return errors.New(errors.TOKEN_VERIFY_ERROR, "token验证失败")
	}

	return nil
}

func (u *UserSvc) PayMargin(c context.Context, req *PayMarginReq) (*PayMarginResp, errors.Error) {
	//修改保证金余额
	_, err := u.userRepo.AddMargin(c, &AddMarginReq{
		UserId: req.UserId,
		Amount: req.Amount,
	})

	if err != nil {
		return nil, err
	}

	//添加记录
	_, err = u.marginRepo.AddMarginRecord(c, &AddMarginRecordReq{
		UserId:      req.UserId,
		Amount:      req.Amount,
		OperateType: MARGIN_OP_ADD,
	})

	return &PayMarginResp{}, err
}

func (u *UserSvc) GetMargin(c context.Context, req *GetMarginReq) (*GetMarginResp, errors.Error) {
	amount, err := u.userRepo.GetMargin(c, req.UserId)
	if err != nil {
		return nil, err
	}

	return &GetMarginResp{Amount: amount}, nil
}

func (u *UserSvc) WithdrawMargin(c context.Context, req *WithdrawMarginReq) (*WithdrawMarginResp, errors.Error) {
	//提现
	_, err := u.userRepo.WithdrawMargin(c, req)
	if err != nil {
		return nil, err
	}

	//记录
	_, err = u.marginRepo.AddMarginRecord(c, &AddMarginRecordReq{
		UserId:      req.UserId,
		Amount:      req.Amount,
		OperateType: MARGIN_OP_WITHDRAW,
	})

	return &WithdrawMarginResp{}, nil
}

func (u *UserSvc) DeductMargin(c context.Context, req *DeductMarginReq) (*DeductMarginResp, errors.Error) {
	//扣除保证金
	_, err := u.userRepo.DeductMargin(c, req)
	if err != nil {
		return nil, err
	}

	//记录
	_, err = u.marginRepo.AddMarginRecord(c, &AddMarginRecordReq{
		UserId:      req.UserId,
		Amount:      req.Amount,
		OperateType: MARGIN_OP_DEDUCT,
	})

	return &DeductMarginResp{}, nil
}
