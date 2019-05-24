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

