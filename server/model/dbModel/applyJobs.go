package dbModel

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"time"
)

type jobType string

const (
	intern      jobType = "intern"
	graduate    jobType = "graduate"
	onlineApply jobType = "onlineApply"
	All         jobType = "all"
)

type CompuseJobs struct {
	BaseModel
	// 职位类型 TODO
	Type jobType `gorm:"type:mood" json:"type"`
	// 职位附加描述
	Benefits    string `json:"benefits"`
	NeedSkills  string `json:"need_skills"`
	WorkContent string `json:"work_content"`
	// 行业领域
	BusinessField pq.StringArray `gorm:"type:text[]" json:"business_field"`
	// 专业领域
	Major pq.StringArray `gorm:"type:text[]" json:"major"`
	// 标签
	Tags         pq.StringArray `gorm:"type:text[]" json:"tags"`
	LocationCity pq.StringArray `gorm:"type:text[]" json:"location_city"`
	Salary       string         `json:"salary"`
	Education    string         `json:"education"`
	// 申请截止时间
	ApplyEndTime *time.Time `json:"apply_end_time"`
	// 浏览次数
	ReviewCounts int64 `json:"review_counts"`
	// 户口
	HasResidence bool `json:"has_residence"`

	CompanyID    string        `gorm:"ForeignKey:CompanyID" json:"company_id"`
	Company      Company       `gorm:"ForeignKey:CompanyID;AssociationForeignKey:CompanyID" json:"company"`
	RecruiterID  string        `gorm:"ForeignKey:RecruiterID" json:"recruiter_id"`
	Recruiter    Recruiter     `gorm:"ForeignKey:Uuid;AssociationForeignKey:Uuid" json:"-"`
	UserApplyJob UserApplyJobs `gorm:"ForeignKey:JobId;AssociationForeignKey:JobId" json:"-"`
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
	gorm.Model `json:"-"`
	JobId      string `gorm:"ForeignKey:JobId" json:"job_id"`
	UserId     string `gorm:"ForeignKey:UserId" json:"user_id"`
	Collected  bool   `gorm:"default:false" json:"collected"`
	Applied    bool   `gorm:"default:false" json:"applied"`
	Talked     bool   `gorm:"default:false" json:"talked"`
}
