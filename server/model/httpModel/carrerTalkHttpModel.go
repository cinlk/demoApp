package httpModel

import "github.com/lib/pq"

type HttpCareerTalkListModel struct {
	MeetingID       string  `json:"meeting_id"`
	CollegeIconURL  string  `json:"college_icon_url"`
	StartTime       tString `json:"start_time"`
	College         string  `json:"college"`
	CompanyName     string  `json:"company_name"`
	SimplifyAddress string  `json:"simplify_address"`
	City            string  `json:"city"`
	//ContentType     string  `json:"content_type"`
	BusinessField pq.StringArray `json:"business_field"`
}

type HttpNearByCareerTalkModel struct {
	MeetingID      string `json:"meeting_id"`
	CollegeIconURL string `json:"college_icon_url"`
	// 公里
	Distance  float64 `json:"distance"`
	StartTime tString `json:"start_time"`
	College   string  `json:"college"`
	Address   string  `json:"address"`
}
