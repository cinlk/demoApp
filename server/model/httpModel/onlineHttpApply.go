package httpModel

import "github.com/lib/pq"

type HttpOnlineApplyModel struct {
	httpBaseMode
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
