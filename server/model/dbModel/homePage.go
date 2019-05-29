package dbModel

import (
	"github.com/jinzhu/gorm"
)

type Banners struct {
	gorm.Model `json:"-"`
	ImageURL   string `json:"image_url"`
	Link       string `gorm:"unique" json:"link"`
}

type LatestNews struct {
	gorm.Model `json:"-"`
	Title      string `gorm:"unique" json:"title"`
}

type JobCategory struct {
	gorm.Model `json:"-"`
	ImageURL   string `json:"image_url"`
	Title      string `gorm:"unique" json:"title"`
}

// 专栏推荐职位? TODO
type TopJobs struct {
	gorm.Model `json:"-"`
	ImageURL   string `json:"image_url"`
	Title      string `gorm:"unique" json:"title"`
	Link       string `json:"link"`
}


// app 相关信息
type AppInfo struct {
	//gorm.Model `json:"-"`
	Wechat string `json:"wechat"`
	ServicePhone string `json:"service_phone"`
	AppId string `gorm:"primary_key" json:"app_id"`
	AppIcon string `json:"app_icon"`
	AppName string `json:"app_name"`
	AppDescribe string `json:"app_describe"`
	Company string `json:"company"`
	Version string `json:"version"`
	CopyRight string `json:"copy_right"`
	AgreeMent string `json:"agree_ment"`

}