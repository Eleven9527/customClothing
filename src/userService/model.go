package userService

import (
	"customClothing/src/db"
	"fmt"
	"github.com/jinzhu/gorm"
)

type User struct {
	UserId      string `json:"userId" gorm:"comment:'用户uuid'"`   //用户uuid
	DisplayName string `json:"displayName" gorm:"comment:'昵称'"`  //昵称
	Password    string `json:"password" gorm:"comment:'密码'"`     //密码
	Avatar      string `json:"avatar" gorm:"comment:'头像图片地址'"`   //头像图片地址
	Role        Role   `json:"role" gorm:"comment:'角色'"`         //角色
	KycId       string `json:"kycId" gorm:"comment:'认证uuid'"`    //认证uuid
	WalletId    string `json:"walletId" gorm:"comment:'钱包uuid'"` //钱包uuid
	Margin      uint64 `json:"margin" gorm:"comment:'保证金'"`      //保证金
	Phone       string `json:"phone" gorm:"comment:'手机号'"`       //手机号
	gorm.Model
}

type Role struct {
	Code int    `json:"code" gorm:"comment:'角色代码'"` //角色代码
	Name string `json:"name" gorm:"comment:'角色名'"`  //角色名
	gorm.Model
}

func InitUserServiceDb() {
	if err := db.Db().AutoMigrate(&User{}).Error; err != nil {
		panic("初始化user表失败:" + err.Error())
	}
	fmt.Println("user表初始化成功!")

	if err := db.Db().AutoMigrate(&Role{}).Error; err != nil {
		panic("初始化user表失败:" + err.Error())
	}
	fmt.Println("role表初始化成功!")
}

type RegisterUserReq struct {
	DisplayName string `json:"displayName"` //昵称，长度6-20
	Password    string `json:"password"`    //密码，长度6-20
	Phone       string `json:"phone"`       //手机号
}

type RegisterUserResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type GetAuthCodeReq struct {
}

type GetAuthCodeResp struct {
	AuthCode string `json:"authCode"`
}

type LoginReq struct {
	Phone    string `json:"phone"`    //手机号
	Password string `json:"password"` //密码
	AuthCode string `json:"authCode"` //验证码
}

type LoginResp struct {
	User  *User
	Token string
}

type SetTokenReq struct {
	Phone string
	Token []byte
}

type SetTokenResp struct {
	Token string
}
