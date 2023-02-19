package userService

import (
	"customClothing/src/db"
	"fmt"
	"github.com/jinzhu/gorm"
)

// 认证审核状态
const (
	KYC_STATUS_PENDING = 1
	KYC_STATUS_SUCCESS = 2
	KYC_STATUS_FAIL    = 3
)

// 角色code
const (
	ROLE_ADMIN    = 0
	ROLE_NORMAL   = 1
	ROLE_DESIGNER = 2
	ROLE_PATTERN  = 3
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

func (User) TableName() string {
	return "user"
}

type Role struct {
	Code int    `json:"code" gorm:"comment:'角色代码'"` //角色代码
	Name string `json:"name" gorm:"comment:'角色名'"`  //角色名
	gorm.Model
}

func (Role) TableName() string {
	return "role"
}

type Kyc struct {
	UserId   string `json:"userId" gorm:"comment:'用户uuid'"`    //用户uuid
	Name     string `json:"name" gorm:"comment:'姓名'"`          //姓名
	Phone    string `json:"phone" gorm:"comment:'手机号'"`        //手机号
	IdCard   string `json:"idCard" gorm:"comment:'手机号'"`       //身份证照片地址
	RoleCode int    `json:"roleCode" gorm:"comment:'申请成为的角色'"` //申请成为的角色
	Status   int    `json:"status" gorm:"comment:'审核状态'"`      //审核状态
	gorm.Model
}

func (Kyc) TableName() string {
	return "kyc"
}

func InitUserServiceDb() {
	if err := db.Db().AutoMigrate(&User{}).Error; err != nil {
		panic("初始化user表失败:" + err.Error())
	}
	fmt.Println("user表初始化成功!")

	if err := db.Db().AutoMigrate(&Role{}).Error; err != nil {
		panic("初始化Role表失败:" + err.Error())
	}
	initRoles()
	fmt.Println("role表初始化成功!")

	if err := db.Db().AutoMigrate(&Kyc{}).Error; err != nil {
		panic("初始化Kyc表失败:" + err.Error())
	}
	fmt.Println("Kyc表初始化成功!")
}

func initRoles() {
	//如果表中有数据，则无需初始化
	var count int64
	if err := db.Db().Table(Role{}.TableName()).Where("code > -1").Count(&count).Error; err != nil {
		panic("Role数据初始化失败:" + err.Error())
	}

	if count > 0 {
		return
	}

	roles := make([]*Role, 0)

	admin := &Role{
		Code: ROLE_ADMIN,
		Name: "管理员",
	}
	roles = append(roles, admin)

	normalUser := &Role{
		Code: ROLE_NORMAL,
		Name: "普通用户",
	}
	roles = append(roles, normalUser)

	designer := &Role{
		Code: ROLE_DESIGNER,
		Name: "设计师",
	}
	roles = append(roles, designer)

	patternMaker := &Role{
		Code: ROLE_PATTERN,
		Name: "版型师",
	}
	roles = append(roles, patternMaker)

	db.Db().Table(Role{}.TableName()).Create(&admin)
	db.Db().Table(Role{}.TableName()).Create(&normalUser)
	db.Db().Table(Role{}.TableName()).Create(&designer)
	db.Db().Table(Role{}.TableName()).Create(&patternMaker)

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
	Phone string //手机号
	Token []byte //用户token
}

type SetTokenResp struct {
	Token string
}

type KycReq struct {
	UserId   string `json:"userId" gorm:"comment:'用户uuid'"`    //用户uuid
	Name     string `json:"name" gorm:"comment:'姓名'"`          //姓名
	Phone    string `json:"phone" gorm:"comment:'手机号'"`        //手机号
	RoleCode int    `json:"roleCode" gorm:"comment:'申请成为的角色'"` //申请成为的角色
}

type KycResp struct {
}
