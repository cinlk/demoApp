package dbModel

import "github.com/jinzhu/gorm"

type ForumArticle struct {
	gorm.Model `json:"-"`
	Uuid       string `gorm:"unique;primary_key" json:"uuid"`
	Title      string `gorm:"not null" json:"title"`
	Content    string `gorm:"not null" json:"content"`
	UserId     string `gorm:"not null" json:"user_id"`
	// 从关联表获取统计数据
	//ThumbUpCount int    `gorm:"default:0" json:"thumb_up_count"`
	//ReplayCount  int    `gorm:"default:0" json:"replay_count"`
	ReadCount int `gorm:"default:0" json:"read_count"`
	// 帖子类型
	Type string `json:"type"`
	// 帖子内容审查结果 TODO
	Validate     bool   `gorm:"default:true" json:"validate"`
	AuditOptions string `json:"audit_options"`
	// 关联的回复
	ReplyItem []ReplyForumPost `gorm:"ForeignKey:PostUuid;AssociationForeignKey:PostUuid"`
}

// 热门帖子(每天离线计算统计 普通帖子中的热门帖子)
type ForumHotestArticle struct {
	gorm.Model `json:"-"`
	// 关联的帖子 id
	Uuid string `gorm:"unique" json:"uuid"`
}

// 一级回贴内容
type ReplyForumPost struct {
	gorm.Model `json:"-"`
	// 帖子id
	PostUuid string `gorm:"ForeignKey:PostUuid" json:"post_uuid"`
	// 回复id
	ReplyId string `gorm:"unique;primary_key" json:"reply_id"`
	// 回复的用户id
	UserId string `gorm:"not null" json:"user_id"`

	// 回复的内容
	Content      string            `json:"content"`
	SecondReplys []SecondReplyPost `gorm:"ForeignKey:ReplyId;AssociationForeignKey:ReplyId"`
}

// 二级回复 （对一级回复的回复）
type SecondReplyPost struct {
	gorm.Model `json:"-"`
	// 一级回复的id
	ReplyId       string `gorm:"ForeignKey:ReplyId;not null" json:"reply_id"`
	UserId        string `gorm:"not null" json:"user_id"`
	Content       string `json:"content"`
	SecondReplyId string `gorm:"unique;primary_key" json:"second_reply_id"`
}

// 三级回复 TODO

// 用户点赞的帖子
type UserLikePost struct {
	gorm.Model `json:"-"`
	// 用户id
	UserId string `gorm:"not null" json:"user_id"`
	// 帖子id
	PostUuid string `gorm:"not null" json:"post_uuid"`
}

// 用户点赞的一级回复
type UserLikeReply struct {
	gorm.Model `json:"-"`
	UserId     string `gorm:"not null" json:"user_id"`
	ReplyId    string `gorm:"not null" json:"reply_id"`
}

// 用户收藏的帖子
type UserCollectedPost struct {
	gorm.Model `json:"-"`
	UserId     string `gorm:"not null" json:"user_id"`
	PostUuid   string `gorm:"not null" json:"post_uuid"`
}
