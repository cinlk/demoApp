package dbModel

import (
	"github.com/jinzhu/gorm"
	"time"
)

// 投递历史记录(网申，校招 和 实习)
type UserDeliveryStatusHistory struct {
	gorm.Model `json:"-"`
	UserId     string     `gorm:"not null" json:"user_id"`
	JobId      string     `gorm:"not null" json:"job_id"`
	Type       JobType    `gorm:"type:mood" json:"type"`
	Status     int        `gorm:"default:0" json:"status"`
	Describe   string     `json:"describe"`
	Time       *time.Time `json:"time"`
}
