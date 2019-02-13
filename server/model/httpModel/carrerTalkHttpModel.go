package httpModel

type HttpCareerTalkListModel struct {
	MeetingID       string  `json:"meeting_id"`
	CollegeIconURL  string  `json:"college_icon_url"`
	StartTime       tString `json:"start_time"`
	College         string  `json:"college"`
	CompanyName     string  `json:"company_name"`
	SimplifyAddress string  `json:"simplify_address"`
	//ContentType     string  `json:"content_type"`
}
