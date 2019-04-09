package dbModel

import "github.com/jinzhu/gorm"

// 单聊
type SingleConversation struct {
	gorm.Model `json:"-"`
	// 用户 的id
	UserID string `json:"user_id"`
	// 用户 的id
	RecruiterID string `json:"recruiter_id"`
	JobID       string `json:"job_id"`
	// 会话有效
	IsValidate bool `gorm:"type:boolean;default:true" json:"is_validate"`

	// 屏蔽 来自谁(seeker 或者 recruiter)
	// d的 我屏蔽 他 我能解除屏蔽
	// conversation remove member TODO
	ShieldFrom     string `json:"shield_from"`
	ConversationID string `gorm:"unique" json:"conversation_id"`
}

// 需要么 TODO
type Message struct {
}
