package httpModel

import "github.com/lib/pq"

type HttpOnlineApplyModel struct {
	Id   string `json:"id"`
	CompanyID   string `json:"company_id"`
	Link string `json:"link,omitempty"`
	Name string `json:"name,omitempty"`
	IconURL string `json:"icon_url,omitempty"`
	CreatedTime tString `json:"created_time"`
	Company HttpCompanyModel `json:"company,omitempty"`
	EndTime  tString `json:"end_time"`
	Citys  pq.StringArray `json:"citys,omitempty"`
	Positions pq.StringArray `json:"positions,omitempty"`
	Major pq.StringArray `json:"major,omitempty"`
	Content string `json:"content"`
	ContentType  string `json:"content_type"`
	OuterSide  bool `json:"outer_side"`

	// user relation
	IsCollected  bool `json:"is_collected"`


}

type HttpOnlineApplyListModel struct {
	OnlineApplyID  string         `json:"online_apply_id"`
	CompanyIconURL string         `json:"company_icon_url"`
	Citys          pq.StringArray `json:"citys"`
	BusinessField  pq.StringArray `json:"business_field"`
	EndTime        tString        `json:"end_time"`
	Name           string         `json:"name"`
	CompanyName    string         `json:"company_name"`
	OutSide        bool           `json:"out_side"`
	Link           string         `json:"link"`
	
}
