package orderService

import (
	"customClothing/src/db"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

var (
	OrderStatus = map[int]string{
		1: "待接单",
		2: "已过期",
		3: "进行中",
		4: "已完成",
		5: "已取消",
	}
	STATUS_PENDING  = OrderStatus[1]
	STATUS_EXPIRED  = OrderStatus[2]
	STATUS_PROCESS  = OrderStatus[3]
	STATUS_COMPLETE = OrderStatus[4]
	STATUS_CANCEL   = OrderStatus[5]
)

var Part = map[int]string{
	1: "上衣",
	2: "裤子",
	3: "裙子",
	4: "套装",
	5: "其他",
	6: "帽子",
}

var WearingOccasion = map[int]string{
	1: "专业运动",
	2: "休闲运动",
	3: "商务正式",
	4: "商务休闲",
	5: "小礼服",
	6: "晚礼服",
	7: "演出服",
	8: "居家休闲",
}

var Style = map[int]string{
	1: "成熟大方",
	2: "活泼清纯",
	3: "低调内敛",
	4: "鲜艳奔放",
	5: "嘻哈自由",
	6: "", //自填
}

var Step = map[int]string{
	1: "需要设计图稿确认",
	2: "需要版型确认",
	3: "需要版型确认制作工艺",
	4: "需要样品成衣确认",
	5: "需要模特展示确认",
}

type Order struct {
	OrderId              string  `json:"orderId" gorm:"comment:'订单uuid'"`              //订单uuid
	PartA                string  `json:"partA" gorm:"comment:'甲方uuid'"`                //甲方uuid
	PartB                string  `json:"partB" gorm:"comment:'乙方uuid'"`                //乙方uuid
	Status               string  `json:"status" gorm:"comment:'订单状态'"`                 //订单状态
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

type OrderDetail struct {
	OrderId                string                 `json:"orderId" gorm:"comment:'订单uuid'"`               //订单uuid
	Cost                   float64                `json:"cost"  gorm:"comment:'费用'"`                     //费用
	DeadLine               time.Time              `json:"deadLine"  gorm:"comment:'截止时间'"`               //截止时间
	ProvideFabric          bool                   `json:"provideFabric"  gorm:"comment:'是否需要提供面料'"`      //是否需要提供面料
	Part                   string                 `json:"part"  gorm:"comment:'制作类别'"`                   //制作类别
	ReferenceSize          ReferenceSize          `json:"referenceSize"  gorm:"comment:'参考尺寸'"`          //参考尺寸
	WearingOccasion        string                 `json:"wearingOccasion"  gorm:"comment:'穿着场合'"`        //穿着场合
	PatternRequest         PatternRequest         `json:"patternRequest"  gorm:"comment:'版型工艺制作要求'"`     //版型工艺制作要求
	Style                  string                 `json:"style"  gorm:"comment:'风格描述'"`                  //风格描述
	ConsumptionPositioning ConsumptionPositioning `json:"consumptionPositioning"  gorm:"comment:'消费定位'"` //消费定位
	ProductionTime         int                    `json:"productionTime"  gorm:"comment:'制作天数要求'"`       //制作天数要求
	ConfirmSteps           string                 `json:"confirmSteps"  gorm:"comment:'步骤确认'"`           //步骤确认
	SampleCost             float64                `json:"sampleCost" gorm:"comment:'成衣样衣费用'"`            //成衣样衣费用
	PatternCost            float64                `json:"patternCost"  gorm:"comment:'版型费用'"`            //版型费用
	Address                Address                `json:"address"  gorm:"comment:'收货地址'"`                //收货地址
	Description            string                 `json:"description"  gorm:"comment:'需求大致介绍'"`          //需求大致介绍
}

func (OrderDetail) TableName() string {
	return "orderDetail"
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

	if err := db.Db().AutoMigrate(&OrderDetail{}).Error; err != nil {
		panic("初始化OrderDetail表失败:" + err.Error())
	}
	fmt.Println("OrderDetail表初始化成功!")
}

type ListOrdersReq struct {
	UserId   string `json:"userId"`   //用户uuid
	RoleCode int    `json:"roleCode"` //用户角色code: 0 = 管理员，1 = 普通用户，2 = 设计师，3 = 版型师
	Status   int    `json:"status"`   //订单状态: 1 = "待接单", 2 = "已过期", 3 = "进行中", 4 = "已完成", 5 = "已取消"
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

type GetOrderSumReq struct {
}

type GetOrderSumResp struct {
	Sum    int     `json:"sum"`    //总交易次数
	Amount float64 `json:"amount"` //总交易金额
}

type Coat struct {
	Bust             float64 `json:"bust"`             //净胸围[单位cm]
	ShoulderWidth    float64 `json:"shoulderWidth"`    //肩宽[单位cm]
	ClothLength      float64 `json:"clothLength"`      //衣长[单位cm]
	SleeveLength     float64 `json:"sleeveLength"`     //袖长[单位cm]
	ArmCircumference float64 `json:"armCircumference"` //手臂围[单位cm]
}

type Bottom struct {
	HipCircumference   float64 `json:"hipCircumference"`   //净臀围[单位cm]
	Waistline          float64 `json:"waistline"`          //净腰围[单位cm]
	Length             float64 `json:"length"`             //裤长[单位cm]
	ThighCircumference float64 `json:"thighCircumference"` //净大腿围[单位cm]
	HeadCircumference  float64 `json:"headCircumference"`  //净头围[单位cm]
}

type ReferenceSize struct {
	Height int    `json:"height"` //身高
	Weight int    `json:"weight"` //体重
	Coat   Coat   `json:"coat"`   //上衣
	Bottom Bottom `json:"bottom"` //下装
}

type PatternRequest struct {
	Looseness        int    `json:"looseness"`        //宽松度：紧身、合体、宽松；
	PatternStructure int    `json:"patternStructure"` //版型结构：立体或者平面；
	Contour          string `json:"contour"`          //轮廓设计：A H X Y T V O；
	LocalDesign      int    `json:"localDesign"`      //局部设计：绣花、印花、领口、袖口、门襟；
	Function         string `json:"function"`         //功能要求：自填；
}

type ConsumptionPositioning struct {
	Age        int    `json:"age"`        //年龄层：18-24，24-28，26-30，33-38等
	Earn       int    `json:"earn"`       //收入范围：5000-8000，7000-9000，10000-11000等
	Profession string `json:"profession"` //职业：自填
}

type Address struct {
	Province string `json:"province"` //省
	City     string `json:"city"`     //市
	District string `json:"district"` //区
	Detail   string `json:"detail"`   //详细地址
	Contact  string `json:"contact"`  //联系人
	Phone    string `json:"phone"`    //手机号
}

type PublishOrderReq struct {
	UserId                 string                 `json:"userId"`                 //甲方uuid
	Cost                   float64                `json:"cost"`                   //费用
	DeadLine               time.Time              `json:"deadLine"`               //截止时间
	ProvideFabric          bool                   `json:"provideFabric"`          //是否需要提供面料
	Part                   int                    `json:"part"`                   //制作类别
	ReferenceSize          ReferenceSize          `json:"referenceSize"`          //参考尺寸
	WearingOccasion        int                    `json:"wearingOccasion"`        //穿着场合
	PatternRequest         PatternRequest         `json:"patternRequest"`         //版型工艺制作要求
	Style                  int                    `json:"style"`                  //风格描述
	CustomStyle            string                 `json:"customStyle"`            //自定义风格描述
	ConsumptionPositioning ConsumptionPositioning `json:"consumptionPositioning"` //消费定位
	ProductionTime         int                    `json:"productionTime"`         //制作天数要求
	ConfirmSteps           []int                  `json:"confirmSteps"`           //步骤确认
	SampleCost             float64                `json:"sampleCost"`             //成衣样衣费用
	PatternCost            float64                `json:"patternCost"`            //版型费用
	Address                Address                `json:"address"`                //收货地址
	Description            string                 `json:"description"`            //需求大致介绍
}

type PublishOrderResp struct {
}

type GetOrderDetailReq struct {
	OrderId string `json:"orderId"` //订单uuid
}

type GetOrderDetailResp struct {
	Detail *OrderDetail `json:"detail"`
}

type PickupOrderReq struct {
	OrderId string `json:"orderId"` //订单uuid
	UserId  string `json:"userId"`  //乙方uuid
}

type PickupOrderResp struct {
}

type DeleteOrderReq struct {
	AdminId string `json:"adminId"` //管理员uuid
	OrderId string `json:"orderId"` //订单uuid
}

type DeleteOrderResp struct {
}
