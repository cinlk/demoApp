package dbModel

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"time"
)

type CarrerTalk struct {
	BaseModel    `json:"-"`
	Address      string         `json:"address"`
	ShortAddress string         `json:"short_address"`
	StartTime    *time.Time     `gorm:"default:now()" json:"start_time"`
	EndTime      *time.Time     `gorm:"default:now()" json:"end_time"`
	Reference    string         `json:"reference"`
	Content      string         `json:"content"`
	ContentType  string         `gorm:"default:'text'" json:"content_type"`
	Field        pq.StringArray `gorm:"type:text[]" json:"field"`
	Majors       pq.StringArray `gorm:"type:text[]" json:"majors"`

	// 关联的公司
	CompanyID string  `gorm:"ForeignKey:CompanyID;unique;not null" json:"company_id"`
	Company   Company `json:"company"`
	// 用户
	UserCarrerTalk []UserCarrerTalk `gorm:"ForeignKey:CarrerTalkID;AssociationForeignKey:CarrerTalkID" json:"users_talks"`
}

//
type UserCarrerTalk struct {
	gorm.Model   `json:"-"`
	UserId       string `json:"user_id"`
	CarrerTalkID string `json:"carrer_talk_id"`
	IsCollected  bool   `gorm:"default:false" json:"is_collected"`
	IsApply      bool   `gorm:"default:false" json:"is_apply"`
}
