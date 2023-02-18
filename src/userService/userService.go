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
}

type UserSvc struct {
	repo  RepoService
	cache CacheService
}

func MakeUserService() UserService {
	return &UserSvc{
		repo:  MakeRepoService(),
		cache: MakeCacheService(),
	}
}

func (u *UserSvc) RegisterUser(c context.Context, req *RegisterUserReq) (*RegisterUserResp, errors.Error) {
	return u.repo.AddUser(c, req)

}

// todo:生成图片验证码
func (u *UserSvc) GetAuthCode(c context.Context, req *GetAuthCodeReq) (*GetAuthCodeResp, errors.Error) {
	return &GetAuthCodeResp{AuthCode: "12345"}, nil
}

func (u *UserSvc) Login(c context.Context, req *LoginReq) (*LoginResp, errors.Error) {
	//检测用户是否存在
	user, err := u.repo.GetUserByPhone(c, req.Phone)
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
