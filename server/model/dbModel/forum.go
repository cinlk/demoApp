package dbModel

import "github.com/jinzhu/gorm"

type ForumArticle struct {
	gorm.Model   `json:"-"`
	Uuid         string `gorm:"unique;primary_key" json:"uuid"`
	Title        string `gorm:"not null" json:"title"`
	UserId       string `gorm:"not null" json:"user_id"`
	ThumbUpCount int    `json:"thumb_up_count"`
	ReplayCount  int    `json:"replay_count"`
	ReadCount    int    `json:"read_count"`
	// 帖子类型
	Type string `json:"type"`
	// 帖子内容审查结果 TODO
	Validate     bool   `gorm:"default:false" json:"validate"`
	AuditOptions string `json:"audit_options"`
}

// 热门帖子(每天离线计算统计 普通帖子中的热门帖子)
type ForumHotestArticle struct {
	gorm.Model `json:"-"`
	// 关联的帖子 id
	Uuid string `gorm:"unique" json:"uuid"`
}
