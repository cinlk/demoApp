package dbModel

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
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


// 订阅职位条件
type JobSubScribeCondition struct {
	gorm.Model `json:"-"`
	UserId string `gorm:"not null" json:"user_id"`
	Type   JobType `gorm:"type:mood" json:"type"`
	// 行业领域
	Fields  string `json:"fields"`
	Citys   pq.StringArray `gorm:"type:text[]" json:"citys"`
	Degree  string `json:"degree"`
	InternDay     string `json:"intern_day"`
	InternMonth   string `json:"intern_month"`
	InternSalary  string `json:"intern_salary"`
	Salary  string `json:"salary"`

}

//type GraduateJobSubscribe struct {
//	gorm.Model `json:"-"`
//	UserId string `gorm:"not null" json:"user_id"`
//	// 行业领域
//	Fields  pq.StringArray `json:"fields"`
//	Citys   pq.StringArray `json:"citys"`
//	Degree  string `json:"degree"`
//	Salary  string `json:"salary"`
//
//
//}

// 控制消息推送开关 TODO
type NotifyMessageSwitch struct {

	gorm.Model `json:"-"`
	// 消息类型 TODO
	Type string `json:"type"`
	UserId string `json:"user_id"`
	Open bool `gorm:"default:true" json:"open"`
}

// 夜间免打扰 开关
type NotifyMessageNightSwitch struct {
	gorm.Model `json:"-"`
	UserId string `gorm:"unique; not null"`
	Open bool `gorm:"default:false" json:"open"`
}


// 求职者 发送的打招呼默认语句 TODO
// 用户创建时 加入数据
type DefaultFirstMessage struct {
	gorm.Model `json:"-"`
	Messages pq.StringArray `gorm:"type:text[]; not null"`
	UserId string `gorm:"unique" json:"user_id"`
	DefaultNum   int  `json:"default_num"`
	// 是否开启打招呼用语
	Open bool  `gorm:"default:true" json:"open"`


}



// 用户提交的反馈信息 TODO
type UserFeedBackMessage struct {
	gorm.Model `json:"-"`
	// 可能有
	UserId  string `json:"user_id"`
	Problem string `json:"problem"`
	Describe string `json:"describe"`
	ImageOneUrl string  `json:"image_one_url"`
	ImageTwoUrl string `json:"image_two_url"`


}




// 用户简历对外可见  完善  TODO
type UserOpenResume struct {
	gorm.Model `json:"-"`
	UserId string `json:"user_id"`
	Open bool `gorm:"default:true" json:"open"`
}