package httpModel

import (
	"github.com/lib/pq"
)

type HttpGraduateModel struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	IsValidate  bool    `json:"is_validate"`
	IsApply     bool    `json:"is_apply"`
	CreatedTime tString `json:"created_time,omitempty"`
	IsCollected bool    `json:"is_collected"`
	IconURL     string  `json:"icon_url"`

	IsTalk bool   `json:"is_talk"`
	Type   string `json:"type"`
	// 福利描述
	Benefits      string                 `json:"benefits"`
	CompanyId     string                 `json:"company_id"`
	Company       HttpSimpleCompanyModel `json:"company,omitempty"`
	RecruiterUUID string                 `json:"recruiter_uuid"`
	Recruiter     HttpRecruiterModel     `json:"recruiter,omitempty"`
	Address       pq.StringArray         `json:"address"`
	City          pq.StringArray         `json:"city"`
	NeedSkills    string                 `json:"need_skills"`
	WorkContent   string                 `json:"work_content"`
	BusinesField  pq.StringArray         `json:"business_field"`
	Major         pq.StringArray         `json:"major"`
	JobTags       pq.StringArray         `json:"job_tags"`
	ReviewCounts  int64                  `json:"review_counts"`

	Salary       string  `json:"salary"`
	Education    string  `json:"education"`
	ApplyEndTime tString `json:"apply_end_time,omitempty"`

	// 会话号
	ConversationId string `json:"conversation_id,omitempty"`
}

type HttpInternJobModel struct {
	HttpGraduateModel
	PayDay      int  `json:"pay_day"`
	Days        int  `json:"days"`
	Months      int  `json:"months"`
	CanTransfer bool `json:"can_transfer"`
}

type HttpJobListModel struct {
	JobID       string         `json:"job_id"`
	Type        string         `json:"type"`
	IconURL     string         `json:"icon_url"`
	CompanyName string         `json:"company_name"`
	JobName     string         `json:"job_name"`
	Address     pq.StringArray `json:"address,omitempty"`
	Degree      string         `json:"degree"`
	ReviewCount int64          `json:"review_count"`
	CreatedTime tString        `json:"created_time,omitempty"`
	// search 需要的数据
	CompanyType   string         `json:"company_type"`
	BusinessField pq.StringArray `json:"business_field,omitempty"`
}

// 实习
type HttpInternListModel struct {
	JobID       string         `json:"job_id"`
	Type        string         `json:"type"`
	IconURL     string         `json:"icon_url"`
	CompanyName string         `json:"company_name"`
	JobName     string         `json:"job_name"`
	Address     pq.StringArray `json:"address,omitempty"`
	Degree      string         `json:"degree"`
	ReviewCount int64          `json:"review_count"`
	CreatedTime tString        `json:"created_time,omitempty"`
	// search 需要的数据
	Days   int `json:"days"`
	Months int `json:"months"`
	PayDay int `json:"pay_day"`
	// 转正
	IsTransfer    bool           `json:"is_transfer"`
	BusinessField pq.StringArray `json:"business_field"`
}

type HttpCompanyTagJobsModel struct {
	Id        string         `json:"id"`
	Type      string         `json:"type"`
	Name      string         `json:"name"`
	Address   pq.StringArray `json:"address,omitempty"`
	Education string         `json:"education"`
	//IconURL     string         `json:"icon_url"`
	CreatedTime tString        `json:"created_time,omitempty"`
	Tags        pq.StringArray `json:"tags,omitempty"`
}
