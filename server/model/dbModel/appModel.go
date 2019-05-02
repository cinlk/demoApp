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

// 系统通知消息
type SystemMessage struct {
	gorm.Model `json:"-"`
	Content    string `json:"content"`
	Title      string `json:"title"`
}

// 用户查看系统消息
type UserCheckSystemMessage struct {
	gorm.Model `json:"-"`
	CheckTime  *time.Time
	// 关联的用户id
	UserId string `gorm:"unique" json:"user_id"`
}

// 最新点赞时间记录
type ForumThumbUpTime struct {
	gorm.Model `json:"-"`
	// 最新点赞时间
	LatestThumbTime *time.Time
	// 自己坚持时间
	CheckTime *time.Time
	// 关联的用户id
	UserId string `gorm:"unique" json:"user_id"`
}

// 最新回复我的帖子时间记录

type ForumReplyMyTime struct {
	gorm.Model      `json:"-"`
	LatestReplyTime *time.Time
	CheckTime       *time.Time
	// 关联的用户id
	UserId string `gorm:"unique" json:"user_id"`
}
