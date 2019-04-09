package httpModel

import "github.com/lib/pq"

// 发布者
type HttpRecruiterModel struct {
	// 系统自带的id
	UserID     string  `json:"user_id"`
	// leancloud 账号id
	LeanCloudAccount string `json:"lean_cloud_account"`
	Name       string  `json:"name"`
	UserIcon   string  `json:"user_icon"`
	Title      string  `json:"title"`
	OnlineTime tString `json:"online_time"`
	Company    string  `json:"company"`
	CompanyID  string  `json:"-"`
}

type HttpRecruiterJobsModel struct {
	Id          string         `json:"id"`
	Name        string         `json:"name"`
	Address     pq.StringArray `json:"address"`
	Education   string         `json:"education"`
	Type        string         `json:"type"`
	CreatedTime tString        `json:"created_time,omitempty"`
}

type HttpRecruiterMainModel struct {
	Recruiter HttpRecruiterModel       `json:"recruiter"`
	Company   HttpSimpleCompanyModel   `json:"company"`
	Jobs      []HttpRecruiterJobsModel `json:"jobs,omitempty"`
}
