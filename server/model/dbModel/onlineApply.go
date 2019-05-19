package dbModel

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"time"
)

type ApplyClassify struct {
	ImageURL string `json:"image_url"`
	Field    string `json:"field"`
}

type OnlineApply struct {
	BaseModel
	EndTime *time.Time `gorm:"default:now()" json:"end_time"`
	// 工作地址
	LocationCity pq.StringArray `gorm:"type:text[]" json:"location_city"`
	Major        pq.StringArray `gorm:"type:text[]" json:"major"`
	// 可投递职位
	//Positions     pq.StringArray `gorm:"type:text[]" json:"positions"`
	Positions     []OnlineApplyPosition `gorm:"ForeignKey:OnlineApplyId;AssociationForeignKey:OnlineApplyId"`
	BusinessField pq.StringArray        `gorm:"type:text[]" json:"business_field"`

	// 是否是外部链接
	Outside bool `gorm:"default:true" json:"outside"`
	// 内容格式 text 或 html
	ContentType contentType `gorm:"type:contentType"  json:"content_type"`
	Content     string      `gorm:"not null" json:"content"`

	CompanyID string `gorm:"ForeignKey:CompanyID" json:"company_id"`
	//Company   Company `gorm:"ForeignKey:CompanyID;AssociationForeignKey:CompanyID" json:"company"`
}

// 需要细化 网申的职位
// 更加职位投递简历
type OnlineApplyPosition struct {
	gorm.Model    `json:"-"`
	OnlineApplyId string `gorm:"ForeignKey:OnlineApplyId"`
	//Uuid string `gorm:""`
	Name string `gorm:"not null" json:"name"`
}

type UserOnlineApplyPosition struct {
	gorm.Model    `json:"-"`
	UserId        string `gorm:"ForeignKey:UserId;not null" json:"user_id"`
	OnlineApplyID string `gorm:"ForeignKey:OnlineApplyID;not null" json:"online_apply_id"`
	// TODO  职位id 已经投递状态记录
	PositionId uint `gorm:"ForeignKey:OnlineApplyPositionId" json:"position_id"`
	//IsApply       bool    `gorm:"default:false" json:"is_apply"`
	//Position      string `json:"position"`
	//IsCollected bool `json:"is_collected"`
	// 0 投递成功 1  2 3 4   当前最新状态
	Status uint `gorm:"default:0" json:"status"`
	// Hr 当前反馈
	FeedBack string `json:"feed_back"`
}

type UserCollectedOnlineApply struct {
	gorm.Model    `json:"-"`
	UserId        string `gorm:"ForeignKey:UserId;not null" json:"user_id"`
	OnlineApplyID string `gorm:"ForeignKey:OnlineApplyID;not null" json:"online_apply_id"`
	IsCollected   bool   `gorm:"default:false" json:"is_collected"`
}
