package dbModel

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type Company struct {
	BaseModel
	Type       string         `json:"type"`
	Latitude   float64        `gorm:"type:numeric" json:"latitude"`
	Longitude  float64        `gorm:"type:numeric" json:"longitude"`
	Citys      pq.StringArray `gorm:"type:text[]" json:"citys"`
	TotalStaff string         `json:"total_staff"`
	// 详细描述
	Describe string `json:"describe"`
	// 简单描述
	SimpleDescribe string         `json:"simple_describe"`
	BusinessField  pq.StringArray `gorm:"type:text[]" json:"business_field"`
	// 公司自身标签
	FeatureTags pq.StringArray `gorm:"type:text[]" json:"feature_tags"`
	// 职位类型标签
	JobTags pq.StringArray `gorm:"type:text[]" json:"job_tags"`
	// 被关注次数(被收藏次数)
	ReviewCounts int64 `json:"review_counts"`
	// 多个talks
	CarrerTalks  []CareerTalk  `gorm:"ForeignKey:companyID" json:"carrer_talks"`
	CompuseJobs  []CompuseJobs `gorm:"ForeignKey:companyID" json:"compuse_jobs"`
	InternJobs   []InternJobs  `gorm:"ForeignKey:companyID" json:"intern_jobs"`
	OnlineApplys []OnlineApply `gorm:"ForeignKey:companyID" json:"online_apply"`
	Recruiters   []Recruiter   `gorm:"ForeignKey:companyID" json:"recruiter"`
}

type UserCompanyRelate struct {
	gorm.Model
	CompanyID   string `gorm:"ForeignKey:CompanyID;not null" json:"company_id"`
	UserID      string `gorm:"ForeignKey:UserId;not null" json:"user_id"`
	IsCollected bool   `gorm:"default:false"`
}
