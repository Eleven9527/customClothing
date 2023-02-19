package orderService

import (
	"customClothing/src/db"
	"fmt"
	"github.com/jinzhu/gorm"
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

func InitOrderServiceDb() {
	if err := db.Db().AutoMigrate(&Order{}).Error; err != nil {
		panic("初始化Kyc表失败:" + err.Error())
	}
	fmt.Println("Order表初始化成功!")
}

type ListOrdersReq struct {
	UserId   string `json:"userId"`   //用户uuid
	RoleCode int    `json:"roleCode"` //用户角色code
	Status   string `json:"status"`   //订单状态
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
