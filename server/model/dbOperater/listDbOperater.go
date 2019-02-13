package dbOperater

import (
	"demoApp/server/model/dbModel"
	"demoApp/server/model/httpModel"
	"github.com/jinzhu/gorm"
	"goframework/orm"
)

type ListDboperater struct {
	orm *gorm.DB
}

func (l *ListDboperater) Banners(limit int) []dbModel.Banners {

	var banner = []dbModel.Banners{}
	_ = orm.DB.Model(&dbModel.Banners{}).Find(&banner).Limit(limit)

	return banner

}

func (l *ListDboperater) LatestNews() []dbModel.LatestNews {

	var news = []dbModel.LatestNews{}
	_ = orm.DB.Model(&dbModel.LatestNews{}).Order("created_at").Find(&news)
	return news
}

func (l *ListDboperater) LatestJobCategory(limit int) []dbModel.JobCategory {
	var jobs = []dbModel.JobCategory{}
	_ = orm.DB.Model(&dbModel.JobCategory{}).Find(&jobs).Limit(limit)

	return jobs
}

func (l *ListDboperater) TopJobs() []dbModel.TopJobs {

	var jobs = []dbModel.TopJobs{}
	_ = orm.DB.Model(&dbModel.TopJobs{}).Find(&jobs)
	return jobs
}

// 首页的 xxx特性的会议数据 TODO
func (l *ListDboperater) CarrerTalks(limit int) []httpModel.HttpCareerTalkListModel {

	var talks = []httpModel.HttpCareerTalkListModel{}

	_ = orm.DB.Model(&dbModel.CareerTalk{}).
		Joins("left join company on company.id = career_talk.company_id").
		Select("career_talk.id as meeting_id, career_talk.college, career_talk.start_time, " +
			"career_talk.simplify_address, career_talk.icon_url as college_icon_url, " +
			"company.name as company_name").
		Limit(limit).
		Scan(&talks)

	return talks
}

func (l *ListDboperater) OnlineApplyClass() []dbModel.ApplyClassify {

	var c = []dbModel.ApplyClassify{}
	_ = orm.DB.Model(&dbModel.ApplyClassify{}).Find(&c)
	return c
}

func (l *ListDboperater) JobList(offset, limit int) []httpModel.HttpJobListModel {

	// first compuse jobs
	// second intern jobs
	// TODO
	jobs := []httpModel.HttpJobListModel{}
	_ = orm.DB.Model(&dbModel.CompuseJobs{}).Joins("left join  company on company.id = compuse_jobs.company_id").
		Select("compuse_jobs.id as job_id, compuse_jobs.icon_url, compuse_jobs.type, compuse_jobs.name as job_name," +
			"compuse_jobs.location_city as address, compuse_jobs.education as degree, " +
			"compuse_jobs.review_counts as review_count, compuse_jobs.created_time, " +
			"company.name as company_name ").
		Offset(offset).Limit(limit).Scan(&jobs)

	return jobs

}

func (l *ListDboperater) Jobs(offset, limit int) []dbModel.CompuseJobs {

	var jobs = []dbModel.CompuseJobs{}

	_ = orm.DB.Model(&dbModel.CompuseJobs{}).Offset(offset).Limit(limit).Find(&jobs)

	// 获取关联的数据 公司 和发布职位者
	for i := 0; i < len(jobs); i++ {
		orm.DB.Model(&dbModel.Company{}).Where("id = ?", jobs[i].CompanyID).Find(&jobs[i].Company)
	}
	return jobs
}

func NewListDboperater() *ListDboperater {

	return &ListDboperater{
		orm: orm.DB,
	}

}
