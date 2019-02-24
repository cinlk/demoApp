package httpModel

import "github.com/lib/pq"

type HttpCompanyModel struct {
	httpBaseMode
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
