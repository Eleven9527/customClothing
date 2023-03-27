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
	WithdrawMarginApplication(c context.Context, req *WithdrawMarginApplicationReq) (*WithdrawMarginApplicationResp, errors.Error)
	DeductMargin(c context.Context, req *DeductMarginReq) (*DeductMarginResp, errors.Error)
	UpdateUserStatus(c context.Context, req *UpdateUserStatusReq) (*UpdateUserStatusResp, errors.Error)
	ReviewMarginWithdrawApplication(c context.Context, req *ReviewMarginWithdrawApplicationReq) (*ReviewMarginWithdrawApplicationResp, errors.Error)
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

// @Summary		注册用户
// @Description	注册用户
// @Tags			user模块
// @Accept			json
// @Produce		json
// @Param			request	body		RegisterUserReq	true	"请求"
// @Success		200		{object}	RegisterUserResp
// @Failure		400		{object}	response.response
// @Failure		404		{object}	response.response
// @Failure		500		{object}	response.response
// @Router			/user [post]
func (u *UserSvc) RegisterUser(c context.Context, req *RegisterUserReq) (*RegisterUserResp, errors.Error) {
	return u.userRepo.AddUser(c, req)

}

//	@Summary		获取验证码
//	@Description	获取验证码
//	@Tags			user模块
//	@Accept			json
//	@Produce		json
//	@Param			None	query		string	false	"无参数"
//	@Success		200		{object}	RegisterUserResp
//	@Failure		400		{object}	response.response
//	@Failure		404		{object}	response.response
//	@Failure		500		{object}	response.response
//	@Router			/user/authcode [get]
//
// GetAuthCode todo:生成图片验证码
func (u *UserSvc) GetAuthCode(c context.Context, req *GetAuthCodeReq) (*GetAuthCodeResp, errors.Error) {
	return &GetAuthCodeResp{AuthCode: "12345"}, nil
}

// @Summary		登录
// @Description	登录
// @Tags			user模块
// @Accept			json
// @Produce		json
// @Param			request	body		LoginReq	true	"请求"
// @Success		200		{object}	LoginResp
// @Failure		400		{object}	response.response
// @Failure		404		{object}	response.response
// @Failure		500		{object}	response.response
// @Router			/user/login [post]
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

// @Summary		Kyc
// @Description	Kyc
// @Tags			user模块
// @Accept			json
// @Produce		json
// @Param			request			body		KycReq	true	"请求"
// @Param			Authorization	header		string	true	"token"
// @Success		200				{object}	KycResp
// @Failure		400				{object}	response.response
// @Failure		404				{object}	response.response
// @Failure		500				{object}	response.response
// @Router			/user/kyc [post]
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

// @Summary		缴纳保证金
// @Description	乙方缴纳保证金
// @Tags			user模块
// @Accept			json
// @Produce		json
// @Param			request			body		PayMarginReq	true	"请求"
// @Param			Authorization	header		string			true	"token"
// @Success		200				{object}	PayMarginResp
// @Failure		400				{object}	response.response
// @Failure		404				{object}	response.response
// @Failure		500				{object}	response.response
// @Router			/user/margin [post]
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

// @Summary		查询保证金
// @Description	查询保证金
// @Tags			user模块
// @Accept			json
// @Produce		json
// @Param			userId			query		string	true	"用户uuid"
// @Param			Authorization	header		string	true	"token"
// @Success		200				{object}	GetMarginResp
// @Failure		400				{object}	response.response
// @Failure		404				{object}	response.response
// @Failure		500				{object}	response.response
// @Router			/user/margin [get]
func (u *UserSvc) GetMargin(c context.Context, req *GetMarginReq) (*GetMarginResp, errors.Error) {
	amount, err := u.userRepo.GetMargin(c, req.UserId)
	if err != nil {
		return nil, err
	}

	return &GetMarginResp{Amount: amount}, nil
}

// @Summary		保证金退回
// @Description	乙方申请保证金退回
// @Tags			user模块
// @Accept			json
// @Produce		json
// @Param			request			body		WithdrawMarginApplicationReq	true	"请求"
// @Param			Authorization	header		string							true	"token"
// @Success		200				{object}	WithdrawMarginApplicationResp
// @Failure		400				{object}	response.response
// @Failure		404				{object}	response.response
// @Failure		500				{object}	response.response
// @Router			/user/marginApplication [post]
func (u *UserSvc) WithdrawMarginApplication(c context.Context, req *WithdrawMarginApplicationReq) (*WithdrawMarginApplicationResp, errors.Error) {
	//添加提现申请
	_, err := u.userRepo.WithdrawMarginApplication(c, req)
	if err != nil {
		return nil, err
	}

	//记录
	_, err = u.marginRepo.AddMarginRecord(c, &AddMarginRecordReq{
		UserId:      req.UserId,
		Amount:      req.Amount,
		OperateType: MARGIN_OP_WITHDRAW,
	})

	return &WithdrawMarginApplicationResp{}, nil
}

// @Summary		扣除保证金
// @Description	管理员扣除保证金
// @Tags			user模块
// @Accept			json
// @Produce		json
// @Param			request			body		DeductMarginReq	true	"请求"
// @Param			Authorization	header		string			true	"token"
// @Success		200				{object}	DeductMarginResp
// @Failure		400				{object}	response.response
// @Failure		404				{object}	response.response
// @Failure		500				{object}	response.response
// @Router			/user/margin [delete]
func (u *UserSvc) DeductMargin(c context.Context, req *DeductMarginReq) (*DeductMarginResp, errors.Error) {
	//只有管理员可以扣除保证金
	user, err := u.userRepo.GetUserById(c, req.UserId)
	if err != nil {
		return nil, err
	}
	if user.RoleCode != ROLE_ADMIN {
		return nil, errors.New(errors.ROLE_ERROR, "只有管理员可以扣除保证金")
	}

	//扣除保证金
	_, err = u.userRepo.DeductMargin(c, req)
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

// @Summary		拉黑、解除拉黑用户
// @Description	管理员拉黑、解除拉黑用户
// @Tags			user模块
// @Accept			json
// @Produce		json
// @Param			request			body		UpdateUserStatusReq	true	"请求"
// @Param			Authorization	header		string				true	"token"
// @Success		200				{object}	UpdateUserStatusResp
// @Failure		400				{object}	response.response
// @Failure		404				{object}	response.response
// @Failure		500				{object}	response.response
// @Router			/user/status [put]
func (u *UserSvc) UpdateUserStatus(c context.Context, req *UpdateUserStatusReq) (*UpdateUserStatusResp, errors.Error) {
	//只有管理员可以拉黑用户
	user, err := u.userRepo.GetUserById(c, req.UserId)
	if err != nil {
		return nil, err
	}
	if user.RoleCode != ROLE_ADMIN {
		return nil, errors.New(errors.ROLE_ERROR, "只有管理员可以拉黑用户")
	}

	_, err = u.userRepo.UpdateUserStatus(c, req)

	return &UpdateUserStatusResp{}, err
}

// @Summary		审核保证金退回申请
// @Description	管理员审核保证金退回申请
// @Tags			user模块
// @Accept			json
// @Produce		json
// @Param			request			body		ReviewMarginWithdrawApplicationReq	true	"请求"
// @Param			Authorization	header		string								true	"token"
// @Success		200				{object}	ReviewMarginWithdrawApplicationResp
// @Failure		400				{object}	response.response
// @Failure		404				{object}	response.response
// @Failure		500				{object}	response.response
// @Router			/user/marginApplication [put]
func (u *UserSvc) ReviewMarginWithdrawApplication(c context.Context, req *ReviewMarginWithdrawApplicationReq) (*ReviewMarginWithdrawApplicationResp, errors.Error) {
	//只有管理员可以审核提现申请
	user, err := u.userRepo.GetUserById(c, req.AdminId)
	if err != nil {
		return nil, err
	}
	if user.RoleCode != ROLE_ADMIN {
		return nil, errors.New(errors.ROLE_ERROR, "只有管理员可以审核提现申请")
	}

	//提现
	_, err = u.userRepo.ReviewMarginWithdrawApplication(c, req)

	return &ReviewMarginWithdrawApplicationResp{}, err
}
