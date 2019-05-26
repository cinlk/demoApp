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

type HttpCareerTalkModel struct {
	Id              string                 `json:"id"`
	Name            string                 `json:"name"`
	Link            string                 `json:"link"`
	IsValidate      bool                   `json:"is_validate"`
	CreatedTime     tString                `json:"created_time,omitempty"`
	IconURL         string                 `json:"icon_url"`
	Company         HttpSimpleCompanyModel `json:"company"`
	College         string                 `json:"college"`
	Address         string                 `json:"address"`
	SimplifyAddress string                 `json:"simplify_address"`
	StartTime       tString                `json:"start_time,omitempty"`
	EndTime         tString                `json:"end_time,omitempty"`
	Content         string                 `json:"content"`
	City 			string 				   `json:"city"`
	ContentType     string                 `json:"content_type"`
	Majors          pq.StringArray         `json:"majors,omitempty"`
	BusinessField   pq.StringArray         `json:"business_field,omitempty"`
	Reference       string                 `json:"reference"`
	CompanyID       string                 `json:"-"`

	//
	IsCollected bool `json:"is_collected"`
}

//type HttpCompanyCareerTalksModel struct {
//	MeetingID string `json:"meeting_id"`
//	//CollegeIconURL string  `json:"college_icon_url"`
//	StartTime tString `json:"start_time,omitempty"`
//	College   string  `json:"college"`
//	//CompanyName     string  `json:"company_name"`
//}
