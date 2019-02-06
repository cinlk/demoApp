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

type TopJobs struct {
	gorm.Model `json:"-"`
	ImageURL   string `json:"image_url"`
	Title      string `gorm:"unique" json:"title"`
	Link       string `json:"link"`
}
