package dbModel

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"time"
)

type contentType string

const (
	Text contentType = "text"
	HTML contentType = "html"
)

type CareerTalk struct {
	BaseModel `json:"-"`
	College   string `json:"college"`
	// 具体地址
	Address string `json:"address"`
	// 简写地址
	SimplifyAddress string     `json:"simplify_address"`
	StartTime       *time.Time `gorm:"default:now()" json:"start_time"`
	EndTime         *time.Time `gorm:"default:now()" json:"end_time"`
	// 引用来源 url 地址 TODO
	Reference string `json:"reference"`
	// 内容
	Content string `gorm:"not null" json:"content"`
	// 内容格式 text 或 html
	ContentType   contentType    `gorm:"type:contentType"  json:"content_type"`
	BusinessField pq.StringArray `gorm:"type:text[]" json:"business_field"`
	Majors        pq.StringArray `gorm:"type:text[]" json:"majors"`

	// 关联的公司
	CompanyID string  `gorm:"ForeignKey:CompanyID;not null" json:"company_id"`
	Company   Company `gorm:"ForeignKey:CompanyID;AssociationForeignKey:CompanyID" json:"company"`
	// 用户
	//UserCarrerTalk []UserApplyCarrerTalk `gorm:"ForeignKey:CarrerTalkID;AssociationForeignKey:CarrerTalkID" json:"users_talks"`
}

//
type UserApplyCarrerTalk struct {
	gorm.Model   `json:"-"`
	UserId       string `gorm:"ForeignKey:UserId;not null" json:"user_id"`
	CarrerTalkID string `gorm:"ForeignKey:CarrerTalkID;not null" json:"carrer_talk_id"`
	IsCollected  bool   `gorm:"default:false" json:"is_collected"`
	//IsApply      bool   `gorm:"default:false" json:"is_apply"`
}