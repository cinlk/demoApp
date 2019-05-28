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

type ResumeListModel struct {
	Name        string  `json:"name"`
	IsPrimary   bool    `json:"is_primary"`
	ResumeId    string  `json:"resume_id"`
	CreatedTime tString `json:"created_time"`
	Type        string  `json:"type"`
}


type textResumeBaseInfo struct {
	Id string `json:"id"`
	Icon string `json:"icon"`
	Name string `json:"name"`
	College string `json:"college"`
	Gender string `json:"gender"`
	City string `json:"city"`
	Degree string `json:"degree"`
	Birthday string `json:"birthday"`
	Phone string `json:"phone"`
	Email string `json:"email"`

}

type textResumeEducation struct {
	Id string `json:"id"`
	College string `json:"college"`
	Degree string `json:"degree"`
	Major string `json:"major"`
	Rank string `json:"rank"`
	Describe string `json:"describe"`
	StartTime resumeTimeString `json:"start_time"`
	EndTime resumeTimeString `json:"end_time"`

}

type textResumeWork struct {
	Id string `json:"id"`
	CompanyName string `json:"company_name"`
	WorkType string `json:"work_type"`
	City string `json:"city"`
	Position string `json:"position"`
	Describe string `json:"describe"`
	StartTime resumeTimeString `json:"start_time"`
	EndTime resumeTimeString `json:"end_time"`
	
}

type textResumeProject struct {
	Id string `json:"id"`
	ProjectName string `json:"project_name"`
	ProjectLevel string `json:"project_level"`
	Position string `json:"position"`
	Describe string `json:"describe"`
	StartTime resumeTimeString `json:"start_time"`
	EndTime resumeTimeString `json:"end_time"`

}

type textResumeCollegeActive struct {
	Id string `json:"id"`
	College string `json:"college"`
	Orgnization string `json:"orgnization"`
	Position string `json:"position"`
	Describe string `json:"describe"`
	StartTime resumeTimeString `json:"start_time"`
	EndTime resumeTimeString `json:"end_time"`

}

type textResumeSocialPractice struct {
	Id string `json:"id"`
	PracticeName string `json:"practice_name"`
	Describe string `json:"describe"`
	StartTime resumeTimeString `json:"start_time"`
	EndTime resumeTimeString `json:"end_time"`
	
}

type textResumeSkill struct {
	Id string `json:"id"`
	SkillName string `json:"skill_name"`
	Describe string `json:"describe"`
}
type textResumeOther struct {
	Id string `json:"id"`
	Title string `json:"title"`
	Describe string `json:"describe"`
}

type textResumeEstimate struct {
	Id string `json:"id"`
	Content string `json:"content"`
}

type TextResumeContentModel struct {
	ResumeId string `json:"resume_id"`
	Level int `json:"level"`
	BaseInfo textResumeBaseInfo `json:"base_person_info"`
	Educations []textResumeEducation `json:"education_infos"`
	Works []textResumeWork `json:"work_infos"`
	Projects []textResumeProject `json:"project_infos"`
	Activities []textResumeCollegeActive `json:"colleage_activities"`
	Practices []textResumeSocialPractice `json:"practice"`
	Skills []textResumeSkill `json:"skills"`
	Others []textResumeOther `json:"other"`
	SelfEstimate textResumeEstimate `json:"self_estimate"`
}


type CollectedJobModel struct {

	JobId string `json:"job_id"`
	CompanyId string `json:"-"`
	IconURL string `json:"icon_url"`
	CompanyName string `json:"company_name"`
	Name string `json:"name"`
	CreatedTime tString `json:"created_time"`
}


type CollectedCompanyModel struct {
	CompanyId string `json:"company_id"`
	IconURL   string `json:"icon_url"`
	Name  	  string `json:"name"`
	Type	  string `json:"type"`
	Citys	  pq.StringArray `json:"citys"`
	CreatedTime tString `json:"created_time"`
}

type CollectedCareerTalkModel struct {
	MeetingId string `json:"meeting_id"`
	CollegeIconUrl string `json:"college_icon_url"`
	CompanyId string `json:"-"`
	CompanyName string `json:"company_name"`
	Name string `json:"name"`
	College string `json:"college"`
	SimplifyAddress string `json:"simplify_address"`
	//CreatedTime tString `json:"created_time"`
}

type CollectedOnlineApplyModel struct {
	OnlineApplyId string `json:"online_apply_id"`
	IconURL string `json:"icon_url"`
	CompanyId string `json:"-"`
	CompanyName string `json:"company_name"`
	Name string `json:"name"`
	Positions []string `json:"positions"`
	CreatedTime tString `json:"created_time"`
	
}

type CollectedPostModel struct {
	PostId string `json:"post_id"`
	Name string `json:"name"`
	GroupName []string `json:"group_name"`
}


type JobSubscribeCondition struct {
	Id string `json:"id"`
	Type string `json:"type"`
	// 行业领域
	Fields  string `json:"fields"`
	Citys   pq.StringArray `json:"citys"`
	Degree  string `json:"degree"`
	InternDay     string `json:"intern_day"`
	InternMonth   string `json:"intern_month"`
	InternSalary  string `json:"intern_salary"`
	Salary  string `json:"salary"`
	CreatedTime tString `json:"created_time"`

}

type NotiyMessageSettings struct {

	Type string `json:"type"`
	Open bool `json:"open"`
}

type UserDefaultTalkMessage struct {
	Messages pq.StringArray `json:"messages"`
	Number int `json:"number"`
	Open bool `json:"open"`
}