package httpModel

import "github.com/lib/pq"

type HttpCompanyModel struct {
	Id            string         `json:"id"`
	Link          string         `json:"link,omitempty"`
	Name          string         `json:"name,omitempty"`
	Type 		  string         `json:"type,omitempty"`
	IconURL       string         `json:"icon_url,omitempty"`
	CreatedTime   tString        `json:"created_time"`
	Describe      string         `json:"describe"`
	SimpleDes     string         `json:"simple_des"`
	Citys         pq.StringArray `json:"citys,omitempty"`
	Staff         string         `json:"staff"`
	IsValidate    bool           `json:"is_validate"`
	IsCollected   bool           `json:"is_collected"`
	Tags          pq.StringArray `json:"tags,omitempty"`
	JobTags       pq.StringArray `json:"job_tags"`
	BusinessField pq.StringArray `json:"business_field,omitempty"`
	// 类型修改  TODO
	ReviewCounts int64 `json:"review_counts"`

	// related data
	Jobs        []HttpJobListModel        `json:"jobs,omitempty"`
	CareerTalks []HttpCareerTalkListModel `json:"career_talks,omitempty"`
}

type HttpNearByCompanyModel struct {
	CompanyID      string         `json:"company_id"`
	CompanyIconURL string         `json:"company_icon_url"`
	CompanyName    string         `json:"company_name"`
	BusinessField  pq.StringArray `json:"business_field"`
	ReviewCount    int64          `json:"review_count"`
	Distance       float64        `json:"distance"`
}

type HttpCompanyListModel struct {
	CompanyID      string         `json:"company_id"`
	CompanyIconURL string         `json:"company_icon_url"`
	CompanyName    string         `json:"company_name"`
	ReviewCounts   int64          `json:"review_counts"`
	Citys          pq.StringArray `json:"citys"`
	BusinessField  pq.StringArray `json:"business_field"`
}

type HttpSimpleCompanyModel struct {
	CompanyID     string         `json:"company_id"`
	IconURL       string         `json:"icon_url"`
	CompanyName   string         `json:"company_name"`
	Citys         pq.StringArray `json:"citys"`
	BusinessField pq.StringArray `json:"business_field"`
	Staff         string         `json:"staff"`
}
