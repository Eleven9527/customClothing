package orderService

import (
	"customClothing/src/db"
	"fmt"
	"github.com/jinzhu/gorm"
)

const (
	STATUS_PENDING  = 1 //待接单
	STATUS_EXPIRED  = 2 //已过期
	STATUS_PROCESS  = 3 //进行中
	STATUS_COMPLETE = 4 //已完成
	STATUS_CANCEL   = 5 //已取消
)

type Order struct {
	OrderId              string  `json:"orderId" gorm:"comment:'订单uuid'"`              //订单uuid
	PartA                string  `json:"partA" gorm:"comment:'甲方uuid'"`                //甲方uuid
	PartB                string  `json:"partB" gorm:"comment:'乙方uuid'"`                //乙方uuid
	Status               int     `json:"status" gorm:"comment:'订单状态'"`                 //订单状态
	Cost                 float64 `json:"cost" gorm:"comment:'费用'"`                     //费用
	DesignArtwork        string  `json:"designArtwork" gorm:"comment:'设计图稿地址'"`        //设计图稿地址
	PatternArtwork       string  `json:"patternArtwork" gorm:"comment:'版型图稿地址'"`       //版型图稿地址
	PatternMakingProcess string  `json:"patternMakingProcess" gorm:"comment:'版型制作工艺'"` //版型制作工艺
	SampleImage          string  `json:"sampleImage" gorm:"comment:'样品成衣图地址'"`         //样品成衣图地址
	ShowVideo            string  `json:"showVideo" gorm:"comment:'模特展示视频地址'"`          //模特展示视频地址
	gorm.Model
}

func (Order) TableName() string {
	return "order"
}

type Reporter struct {
	Whistleblower  string `json:"whistleblower" gorm:"comment:'举报人uuid'"`   //举报人uuid
	ReportedPerson string `json:"reportedPerson" gorm:"comment:'被举报人uuid'"` //被举报人uuid
	Description    string `json:"description" gorm:"comment:'举报描述'"`        //举报描述
	OrderId        string `json:"orderId" gorm:"comment:'订单uuid'"`          //订单uuid
	gorm.Model
}

func (Reporter) TableName() string {
	return "reporter"
}

func InitOrderServiceDb() {
	if err := db.Db().AutoMigrate(&Order{}).Error; err != nil {
		panic("初始化Order表失败:" + err.Error())
	}
	fmt.Println("Order表初始化成功!")

	if err := db.Db().AutoMigrate(&Reporter{}).Error; err != nil {
		panic("初始化Reporter表失败:" + err.Error())
	}
	fmt.Println("Reporter表初始化成功!")
}

type ListOrdersReq struct {
	UserId   string `json:"userId"`   //用户uuid
	RoleCode int    `json:"roleCode"` //用户角色code
	Status   int    `json:"status"`   //订单状态
	PageNum  uint   `json:"pageNum"`  //页数
	PageSize uint   `json:"pageSize"` //每页条数
}

type ListOrdersResp struct {
	Total  int64    `json:"total"`  //总条数
	Orders []*Order `json:"orders"` //订单列表
}

type GetOrderReq struct {
	OrderId string `json:"orderId"` //订单uuid
}

type GetOrderResp struct {
	Order *Order `json:"order"`
}

type UpdateCostReq struct {
	OrderId string  `json:"orderId"` //订单uuid
	Cost    float64 `json:"cost"`    //费用
}

type UpdateCostResp struct {
}

type AddOrderReq struct {
	PartA string  `json:"partA"` //甲方uuid
	Cost  float64 `json:"cost"`  //费用
}

type AddOrderResp struct {
}

type CancelOrderReq struct {
	OrderId string `json:"orderId"` //订单uuid
}

type CancelOrderResp struct {
}

type ConfirmOrderReq struct {
	OrderId string `json:"orderId"` //订单uuid
}

type ConfirmOrderResp struct {
}

type ReportOrderReq struct {
	Whistleblower  string `json:"whistleblower" gorm:"comment:'举报人uuid'"`   //举报人uuid
	ReportedPerson string `json:"reportedPerson" gorm:"comment:'被举报人uuid'"` //被举报人uuid
	Description    string `json:"description" gorm:"comment:'举报描述'"`        //举报描述
	OrderId        string `json:"orderId" gorm:"comment:'订单uuid'"`          //订单uuid
}

type ReportOrderResp struct {
}

type UploadDesignArtworkReq struct {
	OrderId string `json:"orderId"` //订单uuid
}

type UploadDesignArtworkResp struct {
}

type UpdateDesignArtworkReq struct {
	OrderId string `json:"orderId"` //订单uuid
	Url     string `json:"url"`     //图稿url地址
}

type UpdateDesignArtworkResp struct {
}

type UploadPatternArtworkReq struct {
	OrderId string `json:"orderId"` //订单uuid
}

type UploadPatternArtworkResp struct {
}

type UpdatePatternArtworkReq struct {
	OrderId string `json:"orderId"` //订单uuid
	Url     string `json:"url"`     //图稿url地址
}

type UpdatePatternArtworkResp struct {
}

type UploadPatternMakingProcessReq struct {
	OrderId string `json:"orderId"` //订单uuid
	Content string `json:"content"` //工艺内容
}

type UploadPatternMakingProcessResp struct {
}

type UpdatePatternMakingProcessReq struct {
	OrderId string `json:"orderId"` //订单uuid
	Content string `json:"content"` //工艺内容
}

type UpdatePatternMakingProcessResp struct {
}

type UploadSampleImageReq struct {
	OrderId string `json:"orderId"` //订单uuid
	Url     string `json:"url"`     //成衣图url地址
}

type UploadSampleImageResp struct {
}

type UpdateSampleImageReq struct {
	OrderId string `json:"orderId"` //订单uuid
	Url     string `json:"url"`     //成衣图url地址
}

type UpdateSampleImageResp struct {
}

type UploadShowVideoReq struct {
	OrderId string `json:"orderId"` //订单uuid
	Url     string `json:"url"`     //视频url地址
}

type UploadShowVideoResp struct {
}

type UpdateShowVideoReq struct {
	OrderId string `json:"orderId"` //订单uuid
	Url     string `json:"url"`     //视频url地址
}

type UpdateShowVideoResp struct {
}
