package dbModel

import (
	"github.com/jinzhu/gorm"
	"time"
)

type ResumeType string

func (r ResumeType) Validate() bool {
	switch r {
	case ResumeText, ResumeAttach:
		return true
	default:
		return false
	}
}

const (
	ResumeText   ResumeType = "text"
	ResumeAttach ResumeType = "attache"
)

// 投递历史记录(网申，校招 和 实习)
type UserDeliveryStatusHistory struct {
	gorm.Model `json:"-"`
	UserId     string     `gorm:"not null" json:"user_id"`
	JobId      string     `gorm:"not null" json:"job_id"`
	Type       JobType    `gorm:"type:mood" json:"type"`
	Status     int        `gorm:"default:0" json:"status"`
	Describe   string     `json:"describe"`
	Time       *time.Time `json:"time"`
}

// 个人简历 附件简历 和 文本简历

type MyResume struct {
	gorm.Model `json:"-"`

	UserId string `gorm:"not null" json:"user_id"`
	Uuid   string `gorm:"unique; primary_key"`
	// 简历名称
	Name string `gorm:"not null"`
	// 附件 和 文本简历
	Type string `gorm:"type:resume;not null"`
	// 默认投递的简历(只有一个)
	IsPrimary bool `gorm:"default:false" json:"is_primary"`

	AttachResume AttachFileResume `gorm:"ForeignKey:ResumeId;AssociationForeignKey:ResumeId" json:"attach_resume"`
	TextResume   TextResume       `gorm:"ForeignKey:ResumeId;AssociationForeignKey:ResumeId" json:"text_resume"`
}

type AttachFileResume struct {
	gorm.Model `json:"-"`
	// 关联的简历id
	ResumeId string `gorm:"unique;not null" json:"resume_id"`
	// 文件在七牛云的地址
	FileUrl string `gorm:"not null" json:"file_url"`
	// 关联的用户id
	//UserId string `gorm:"not null" json:"user_id"`
}

type TextResume struct {
	gorm.Model `json:"-"`
	// 关联的简历id
	ResumeId string `gorm:"unique;not null" json:"resume_id"`
	// 完善程度
	CompleteLevel int `gorm:"default:1"`
	// 关联的表
	BaseInfo TextResumeBaseInfo `gorm:"ForeignKey:ResumeId;AssociationForeignKey:ResumeId"`
	// 教育
	Education []TextResumeEducation  `gorm:"ForeignKey:ResumeId;AssociationForeignKey:ResumeId"`
	Works []TextResumeWorkExperience `gorm:"ForeignKey:ResumeId;AssociationForeignKey:ResumeId"`
	Projects []TextResumeProject `gorm:"ForeignKey:ResumeId;AssociationForeignKey:ResumeId"`
	CollegeActive []TextResumeCollegeActivity  `gorm:"ForeignKey:ResumeId;AssociationForeignKey:ResumeId"`
	SocialPractice []TextResumeSocialPractice `gorm:"ForeignKey:ResumeId;AssociationForeignKey:ResumeId"`
	Skills []TextResumeSkills `gorm:"ForeignKey:ResumeId;AssociationForeignKey:ResumeId"`
	Others []TextResumeOther `gorm:"ForeignKey:ResumeId;AssociationForeignKey:ResumeId"`
	Estimate TextResumeEstimate `gorm:"ForeignKey:ResumeId;AssociationForeignKey:ResumeId"`
}

type TextResumeBaseInfo struct {
	gorm.Model `json:"-"`
	ResumeId string `gorm:"not null;unique" json:"resume_id"`
	Icon string `json:"icon"`
	Name string `json:"name"`
	College string `json:"college"`
	Gender string `json:"gender"`
	City string `json:"city"`
	Degree string `json:"degree"`
	Birthday string `json:"birthday"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}


type TextResumeEducation struct {
	gorm.Model `json:"-"`
	ResumeId string `gorm:"not null" json:"resume_id"`
	College string `json:"college"`
	// 专业
	Major string `json:"major"`
	// 排名
	Rank string `json:"rank"`
	// 学历
	Degree string `json:"degree"`

	// 描述
	Describe string `json:"describe"`

	StartTime *time.Time `json:"start_time"`
	EndTime *time.Time  `json:"end_time"`


}

// 实习/兼职/工作经历
type TextResumeWorkExperience struct {
	gorm.Model `json:"-"`
	ResumeId string `gorm:"not null" json:"resume_id"`
	CompanyName string `json:"company_name"`
	WorkType string `json:"work_type"`
	City string `json:"city"`
	Position string `json:"position"`
	// 描述
	Describe string `json:"describe"`
	StartTime *time.Time `json:"start_time"`
	EndTime *time.Time  `json:"end_time"`

}

// 参与的项目或比赛
type TextResumeProject struct {
	gorm.Model `json:"-"`
	ResumeId string `gorm:"not null" json:"resume_id"`

	ProjectName string `json:"project_name"`
	// 人数规模 ??
	ProjectLevel  string `json:"project_level"`
	Position string `json:"position"`

	Describe string `json:"describe"`
	StartTime *time.Time `json:"start_time"`
	EndTime *time.Time  `json:"end_time"`

}

type TextResumeCollegeActivity struct {
	gorm.Model `json:"-"`
	ResumeId string `gorm:"not null" json:"resume_id"`
	College string `json:"college"`
	Orgnization string `json:"orgnization"`
	Position string `json:"position"`

	Describe string `json:"describe"`
	StartTime *time.Time `json:"start_time"`
	EndTime *time.Time  `json:"end_time"`

}

type TextResumeSocialPractice struct {
	gorm.Model `json:"-"`
	ResumeId string `gorm:"not null" json:"resume_id"`
	Name string
	Describe string `json:"describe"`
	StartTime *time.Time `json:"start_time"`
	EndTime *time.Time  `json:"end_time"`

}

type TextResumeSkills struct {
	gorm.Model `json:"-"`
	ResumeId string `gorm:"not null" json:"resume_id"`
	SkillName string `json:"skill_name"`
	Describe string `json:"describe"`

}

type TextResumeOther struct {
	gorm.Model `json:"-"`
	ResumeId string `gorm:"not null" json:"resume_id"`
	Title string
	Describe string `json:"describe"`

}

type TextResumeEstimate struct {
	gorm.Model `json:"-"`
	ResumeId string `gorm:"not null" json:"resume_id"`
	Content string `json:"content"`
}