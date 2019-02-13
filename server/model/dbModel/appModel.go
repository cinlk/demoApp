package dbModel

import (
	"github.com/jinzhu/gorm"
	"time"
)

type AppGuidanceItem struct {
	gorm.Model `json:"-"`
	ImageURL   string `json:"image_url"`
	Title      string `json:"title"`
	Detail     string `json:"detail"`
}

type BaseModel struct {
	Id          string     `gorm:"primary_key;unique" json:"id"`
	Link        string     `json:"link"`
	Name        string     `json:"name"`
	IconURL     string     `json:"icon_url"`
	IsValidate  bool       `gorm:"default:true" json:"is_validate"`
	CreatedTime *time.Time `gorm:"default:now()" json:"created_time"`
}

// 新闻栏目
type NewsModel struct {
	gorm.Model  `json:"-"`
	Title       string     `json:"title"`
	Author      string     `json:"author"`
	Icon        string     `json:"icon"`
	Link        string     `json:"link"`
	Time        *time.Time `gorm:"default:now()" json:"-"`
	CreatedTime int64      `gorm:"-" json:"time"`
}
