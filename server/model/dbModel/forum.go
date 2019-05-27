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
	ReplyId string `gorm:"ForeignKey:ReplyId;not null" json:"reply_id"`
	// 发起回复的用户
	UserId string `gorm:"not null" json:"user_id"`
	// 目标回复用户, 默认是一级回复的用户
	TalkedUserId  string `gorm:"not null" json:"talked_user_id"`
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

// 用户点赞的二级回复
type UserLikeSubReply struct {
	gorm.Model    `json:"-"`
	UserId        string `gorm:"not null" json:"user_id"`
	SecondReplyId string `gorm:"not null" json:"second_reply_id"`
}

// 用户收藏的帖子
type UserCollectedPost struct {
	gorm.Model `json:"-"`
	UserId     string `gorm:"not null" json:"user_id"`
	PostUuid   string `gorm:"not null" json:"post_uuid"`
	Groups      []UserCollectedGroup `gorm:"many2many:post_group" json:"group"`
	
}

// 用户举报的帖子
type UserAlertPost struct {
	gorm.Model `json:"-"`
	UserId     string `gorm:"unique; not null" json:"user_id"`
	PostId     string `gorm:"not null" json:"post_id"`
	Content    string `gorm:"not null" json:"content"`
}

// 用户举报一级回复
type UserAlertReply struct {
	gorm.Model `json:"-"`
	UserId     string `gorm:"unique;not null" json:"user_id"`
	Content    string
	ReplyId    string `gorm:"not null" json:"reply_id"`
}

// 用户举报二级回复
type UserAlertSubReply struct {
	gorm.Model    `json:"-"`
	UserId        string `gorm:"not null" json:"user_id"`
	Content       string `json:"content"`
	SecondReplyId string `gorm:"not null" json:"second_reply_id"`
}


// 用户收藏的帖子分组 (默认分组)
type UserCollectedGroup struct {
	gorm.Model `json:"-"`
	GroupName string `json:"group_name"`
	UserId string `json:"user_id"`
	Posts []UserCollectedPost `gorm:"many2many:post_group" json:"posts"`
}

 


