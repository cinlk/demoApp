package dbModel

import (
	"github.com/go-pg/pg/orm"
	"time"
)

type AppGuidanceItem struct {
	orm.Model `json:"-"`
	ImageURL  string `json:"image_url"`
	Title     string `json:"title"`
	Detail    string `json:"detail"`
}

type BaseModel struct {
	Id          string     `gorm:"primary_key;unique" json:"id"`
	Link        string     `json:"link"`
	Name        string     `json:"name"`
	Icon        string     `json:"icon"`
	IsValidate  bool       `gorm:"default:true" json:"is_validate"`
	CreatedTime *time.Time `gorm:"default:now()" json:"created_time"`
}
