package dbModel

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"time"
)

type JobType string

func (t JobType) Validate() bool {
	switch t {
	case intern, graduate, onlineApply:
		return true
	default:
		return false
	}
}

const (
	intern      JobType = "intern"
	graduate    JobType = "graduate"
	onlineApply JobType = "onlineApply"
	All         JobType = "all"
)

type CompuseJobs struct {
	BaseModel
	// 职位类型 TODO
	Type JobType `gorm:"type:mood" json:"type"`
	// 职位附加描述
	Benefits string `json:"benefits"`
	// 职位地址
	Address     pq.StringArray `gorm:"type:text[]" json:"address"`
	NeedSkills  string         `json:"need_skills"`
	WorkContent string         `json:"work_content"`
	// 行业领域
	BusinessField pq.StringArray `gorm:"type:text[]" json:"business_field"`
	// 专业领域
	Major pq.StringArray `gorm:"type:text[]" json:"major"`
	// 职位类型标签（公司检索职位使用）
	Tags         pq.StringArray `gorm:"type:text[]" json:"tags"`
	LocationCity pq.StringArray `gorm:"type:text[]" json:"LocationCity"`
	Salary       string         `json:"salary"`
	Education    string         `json:"education"`
	// 申请截止时间
	ApplyEndTime *time.Time `json:"apply_end_time"`
	// 浏览次数
	ReviewCounts int64 `json:"review_counts"`
	// 户口
	HasResidence bool `json:"has_residence"`

	CompanyID string `gorm:"ForeignKey:CompanyID" json:"company_id"`
	//Company       Company       `gorm:"ForeignKey:CompanyID;AssociationForeignKey:CompanyID" json:"company"`
	RecruiterUUID string `gorm:"ForeignKey:RecruiterUUID" json:"recruiter_id"`
	//Recruiter     Recruiter     `gorm:"ForeignKey:Uuid;AssociationForeignKey:Uuid" json:"-"`
	//UserApplyJob  UserApplyJobs `gorm:"ForeignKey:JobId;AssociationForeignKey:JobId" json:"-"`

	Conversation SingleConversation `gorm:"AssociationForeignKey:JobID"`
}

type InternJobs struct {
	CompuseJobs
	// 每周实习天数
	Days   int `json:"days"`
	Months int `json:"months"`
	// 实习日薪
	PayDay int `json:"pay_day"`
	// 可以转正
	CanTransfer bool `json:"can_transfer"`
}

type UserApplyJobs struct {
	gorm.Model  `json:"-"`
	JobId       string `gorm:"ForeignKey:JobId" json:"job_id"`
	UserId      string `gorm:"ForeignKey:UserId" json:"user_id"`
	IsCollected bool   `gorm:"default:false" json:"is_collected"`
	IsApply     bool   `gorm:"default:false" json:"is_apply"`
	IsTalk      bool   `gorm:"default:false" json:"is_talk"`
	// 投递的状态 0 投递成功  1 2 3 4
	Status int `gorm:"default:0"`
	// Hr 当前反馈
	FeedBack string  `json:"feed_back"`
	JobType  JobType `gorm:"type:mood" json:"type"`
}
