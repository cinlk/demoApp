package dbModel

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type Notify string

const (
	System  Notify = "system"
	Channel Notify = "channel"
	Special Notify = "special"
)

// 通知类型表
type NotifyMessage struct {
	gorm.Model `json:"-"`
	// 通知类型
	Type    Notify `gorm:"type:notify" json:"type"`
	Content string `json:"content"`
	Title   string `json:"title"`
	// 需要通知的用户
	UserId pq.StringArray `gorm:"type:text[]"`
	// 推送的channels
	Channels pq.StringArray `gorm:"type:text[]"`
	// 通知id (通知结果 TODO)
	// 发送成功后，返回的用objectId 是该id，用于查询发送记录
	NotificationId string `gorm:"unique;not null"`
}
