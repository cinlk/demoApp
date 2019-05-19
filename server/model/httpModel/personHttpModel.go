package httpModel

import "github.com/lib/pq"

type HttpPersonAvatarModel struct {
	IconUrl string `json:"icon_url"`
}

type DeliveryJob struct {
	Type        string         `json:"type"`
	JobId       string         `json:"job_id"`
	JobName     string         `json:"job_name"`
	CreatedTime tString        `json:"created_time"`
	Status      int            `json:"status"`
	CompanyName string         `json:"company_name"`
	CompanyId   string         `json:"company_id"`
	CompanyIcon string         `json:"company_icon"`
	Address     pq.StringArray `json:"address"`
	FeedBack    string         `json:"feed_back"`
}

type DeliveryJobStatusHistory struct {
	Status   int     `json:"status"`
	Time     tString `json:"time"`
	Describe string  `json:"describe"`
}
