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
	Positions     pq.StringArray `gorm:"type:text[]" json:"positions"`
	BusinessField pq.StringArray `gorm:"type:text[]" json:"business_field"`

	// 是否是外部链接
	Outside bool `gorm:"default:true" json:"outside"`
	// 内容格式 text 或 html
	ContentType contentType `gorm:"type:contentType"  json:"content_type"`
	Content     string      `gorm:"not null" json:"content"`

	CompanyID string `gorm:"ForeignKey:CompanyID" json:"company_id"`
	//Company   Company `gorm:"ForeignKey:CompanyID;AssociationForeignKey:CompanyID" json:"company"`
}

type UserOnlineApply struct {
	gorm.Model    `json:"-"`
	UserId        string `gorm:"ForeignKey:UserId;not null" json:"user_id"`
	OnlineApplyID string `gorm:"ForeignKey:OnlineApplyID;not null" json:"online_apply_id"`
	// TODO 按照xx申请 控制次数
	//IsApply       bool   `gorm:"default:false" json:"is_apply"`
	//Position      string `json:"position"`
	IsCollected bool `json:"is_collected"`
}
