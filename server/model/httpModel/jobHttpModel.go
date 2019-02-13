package httpModel

import (
	"github.com/lib/pq"
)

type HttpJobModel struct {
	httpBaseMode

	Validate      bool                `json:"validate"`
	Applied       bool                `json:"applied"`
	CreatedTime   int64               `json:"created_time"`
	Collected     bool                `json:"collected"`
	Talked        bool                `json:"is_talk"`
	Type          string              `json:"type"`
	Benefits      string              `json:"benefits"`
	Comapany      *HttpCompanyModel   `json:"comapany,omitempty"`
	Recruiter     *HttpRecruiterModel `json:"hr,omitempty"`
	NeedSkills    string              `json:"requirement"`
	WorkContent   string              `json:"works"`
	BussinesField []string            `json:"industry"`
	Major         []string            `json:"major"`
	Tags          []string            `json:"job_tags"`
	ReviewCounts  int64               `json:"read_num"`
	LocationCity  []string            `json:"city"`
	Salary        string              `json:"salary"`
	Education     string              `json:"education"`
	ApplyEndTime  int64               `json:"apply_end_time"`
}

type HttpInternJobModel struct {
	HttpJobModel
	Days        int `json:"per_day"`
	Months      int `json:"month"`
	CanTransfer int `json:"is_staff"`
}

type HttpJobListModel struct {
	JobID       string         `json:"job_id"`
	Type        string         `json:"type"`
	IconURL     string         `json:"icon_url"`
	CompanyName string         `json:"company_name"`
	JobName     string         `json:"job_name"`
	Address     pq.StringArray `json:"address"`
	Degree      string         `json:"degree"`
	ReviewCount int64          `json:"review_count"`
	CreatedTime tString        `json:"created_time"`
}
